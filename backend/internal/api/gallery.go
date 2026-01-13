// Package api provides HTTP handlers for the gallery endpoints.
package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"better-kiro-prompts/internal/gallery"
	"better-kiro-prompts/internal/ratelimit"
)

// GalleryHandler holds dependencies for gallery endpoints.
type GalleryHandler struct {
	service       *gallery.Service
	ratingLimiter *ratelimit.Limiter
}

// NewGalleryHandler creates a new handler with the given dependencies.
func NewGalleryHandler(service *gallery.Service, ratingLimiter *ratelimit.Limiter) *GalleryHandler {
	return &GalleryHandler{
		service:       service,
		ratingLimiter: ratingLimiter,
	}
}

// GalleryListResponse is the response for listing gallery items.
type GalleryListResponse struct {
	Items      []GalleryItem `json:"items"`
	Total      int           `json:"total"`
	Page       int           `json:"page"`
	PageSize   int           `json:"pageSize"`
	TotalPages int           `json:"totalPages"`
}

// GalleryItem represents a gallery item in list responses.
type GalleryItem struct {
	ID          string  `json:"id"`
	ProjectIdea string  `json:"projectIdea"`
	Category    string  `json:"category"`
	AvgRating   float64 `json:"avgRating"`
	RatingCount int     `json:"ratingCount"`
	ViewCount   int     `json:"viewCount"`
	CreatedAt   string  `json:"createdAt"`
	Preview     string  `json:"preview"`
}

// GalleryDetailResponse is the response for a single gallery item.
type GalleryDetailResponse struct {
	Generation GalleryDetail `json:"generation"`
	UserRating int           `json:"userRating"`
}

// GalleryDetail represents full generation details.
type GalleryDetail struct {
	ID              string          `json:"id"`
	ProjectIdea     string          `json:"projectIdea"`
	ExperienceLevel string          `json:"experienceLevel"`
	HookPreset      string          `json:"hookPreset"`
	Files           json.RawMessage `json:"files"`
	Category        string          `json:"category"`
	AvgRating       float64         `json:"avgRating"`
	RatingCount     int             `json:"ratingCount"`
	ViewCount       int             `json:"viewCount"`
	CreatedAt       string          `json:"createdAt"`
}

// RateRequest is the request body for rating a generation.
type RateRequest struct {
	Score     int    `json:"score"`
	VoterHash string `json:"voterHash"`
}

// RateResponse is the response for rating a generation.
type RateResponse struct {
	Success bool `json:"success"`
}

// HandleListGallery handles GET /api/gallery.
func (h *GalleryHandler) HandleListGallery(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	query := r.URL.Query()

	// Parse category filter
	var categoryID *int
	if catStr := query.Get("category"); catStr != "" {
		cat, err := strconv.Atoi(catStr)
		if err != nil {
			WriteValidationError(w, r, "Invalid category ID")
			return
		}
		categoryID = &cat
	}

	// Parse sort option
	sortBy := query.Get("sort")
	if sortBy == "" {
		sortBy = "newest"
	}

	// Parse pagination
	page := 1
	if pageStr := query.Get("page"); pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err != nil || p < 1 {
			WriteValidationError(w, r, "Invalid page number")
			return
		}
		page = p
	}

	pageSize := gallery.DefaultPageSize
	if sizeStr := query.Get("pageSize"); sizeStr != "" {
		s, err := strconv.Atoi(sizeStr)
		if err != nil || s < 1 {
			WriteValidationError(w, r, "Invalid page size")
			return
		}
		pageSize = s
	}

	// Call service
	resp, err := h.service.ListGenerations(r.Context(), gallery.ListRequest{
		CategoryID: categoryID,
		SortBy:     sortBy,
		Page:       page,
		PageSize:   pageSize,
	})
	if err != nil {
		if errors.Is(err, gallery.ErrInvalidSort) {
			WriteValidationError(w, r, "Invalid sort option")
			return
		}
		WriteInternalError(w, r, "")
		return
	}

	// Convert to response format
	items := make([]GalleryItem, len(resp.Items))
	for i, gen := range resp.Items {
		items[i] = GalleryItem{
			ID:          gen.ID,
			ProjectIdea: gen.ProjectIdea,
			Category:    gen.CategoryName,
			AvgRating:   gen.AvgRating,
			RatingCount: gen.RatingCount,
			ViewCount:   gen.ViewCount,
			CreatedAt:   gen.CreatedAt.Format("2006-01-02T15:04:05Z"),
			Preview:     truncateString(gen.ProjectIdea, 200),
		}
	}

	writeJSON(w, http.StatusOK, GalleryListResponse{
		Items:      items,
		Total:      resp.Total,
		Page:       resp.Page,
		PageSize:   resp.PageSize,
		TotalPages: resp.TotalPages,
	})
}

// HandleGetGalleryItem handles GET /api/gallery/{id}.
func (h *GalleryHandler) HandleGetGalleryItem(w http.ResponseWriter, r *http.Request) {
	// Extract ID from path using Go 1.22+ PathValue
	id := r.PathValue("id")
	if id == "" {
		WriteValidationError(w, r, "Invalid generation ID")
		return
	}

	// Hash the client IP for view tracking and rating lookup
	clientIP := getClientIP(r)
	ipHash := hashIP(clientIP)

	// Get generation with IP-deduplicated view tracking
	gen, err := h.service.GetGenerationWithView(r.Context(), id, ipHash)
	if err != nil {
		if errors.Is(err, gallery.ErrNotFound) {
			WriteNotFound(w, r, "Generation not found")
			return
		}
		if errors.Is(err, gallery.ErrInvalidInput) {
			WriteValidationError(w, r, "Invalid generation ID")
			return
		}
		WriteInternalError(w, r, "")
		return
	}

	// Get user rating using IP hash (Requirements 5.2, 5.4)
	userRating, _ := h.service.GetUserRating(r.Context(), id, ipHash)

	writeJSON(w, http.StatusOK, GalleryDetailResponse{
		Generation: GalleryDetail{
			ID:              gen.ID,
			ProjectIdea:     gen.ProjectIdea,
			ExperienceLevel: gen.ExperienceLevel,
			HookPreset:      gen.HookPreset,
			Files:           gen.Files,
			Category:        gen.CategoryName,
			AvgRating:       gen.AvgRating,
			RatingCount:     gen.RatingCount,
			ViewCount:       gen.ViewCount,
			CreatedAt:       gen.CreatedAt.Format("2006-01-02T15:04:05Z"),
		},
		UserRating: userRating,
	})
}

// HandleRateGalleryItem handles POST /api/gallery/{id}/rate.
// Uses IP hash for vote deduplication per Requirements 5.2, 5.4, 5.5.
func (h *GalleryHandler) HandleRateGalleryItem(w http.ResponseWriter, r *http.Request) {
	// Extract ID from path using Go 1.22+ PathValue
	id := r.PathValue("id")
	if id == "" {
		WriteValidationError(w, r, "Invalid generation ID")
		return
	}

	// Parse request body
	var req RateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteBadRequest(w, r, "Invalid request body")
		return
	}

	// Validate score
	if req.Score < 1 || req.Score > 5 {
		WriteValidationError(w, r, "Score must be between 1 and 5")
		return
	}

	// Check rating rate limit
	ip := getClientIP(r)
	if h.ratingLimiter != nil {
		allowed, retryAfter := h.ratingLimiter.Allow(ip)
		if !allowed {
			WriteRateLimited(w, r, int(retryAfter.Seconds()))
			return
		}
	}

	// Use IP hash for voter identification (Requirements 5.2, 5.4, 5.5)
	// This ensures one vote per IP address per generation
	ipHash := hashIP(ip)

	// Submit rating using IP hash for deduplication
	retryAfter, err := h.service.RateGeneration(r.Context(), id, req.Score, ipHash, ip)
	if err != nil {
		if errors.Is(err, gallery.ErrNotFound) {
			WriteNotFound(w, r, "Generation not found")
			return
		}
		if errors.Is(err, gallery.ErrInvalidRating) {
			WriteValidationError(w, r, "Score must be between 1 and 5")
			return
		}
		if errors.Is(err, gallery.ErrInvalidInput) {
			WriteValidationError(w, r, "Invalid input")
			return
		}
		if errors.Is(err, gallery.ErrRateLimited) {
			WriteRateLimited(w, r, retryAfter)
			return
		}
		WriteInternalError(w, r, "")
		return
	}

	writeJSON(w, http.StatusOK, RateResponse{Success: true})
}

// truncateString truncates a string to the given length, adding "..." if truncated.
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// hashIP creates a SHA-256 hash of an IP address for privacy-preserving storage.
// The hash is returned as a lowercase hex string.
func hashIP(ip string) string {
	hash := sha256.Sum256([]byte(ip))
	return hex.EncodeToString(hash[:])
}
