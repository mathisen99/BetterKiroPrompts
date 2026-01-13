// Package gallery provides the gallery service for browsing and rating generations.
package gallery

import (
	"context"
	"errors"
	"math"

	"better-kiro-prompts/internal/ratelimit"
	"better-kiro-prompts/internal/storage"
)

// Service errors
var (
	ErrNotFound      = errors.New("generation not found")
	ErrInvalidInput  = errors.New("invalid input")
	ErrRateLimited   = errors.New("rate limited")
	ErrInvalidRating = errors.New("rating must be between 1 and 5")
	ErrInvalidPage   = errors.New("page must be positive")
	ErrInvalidSort   = errors.New("invalid sort option")
)

// DefaultPageSize is the default number of items per page.
const DefaultPageSize = 20

// MaxPageSize is the maximum allowed page size.
const MaxPageSize = 100

// ValidSortOptions defines the allowed sort options.
var ValidSortOptions = map[string]bool{
	"newest":        true,
	"highest_rated": true,
	"most_viewed":   true,
}

// ListRequest contains parameters for listing generations.
type ListRequest struct {
	CategoryID *int
	SortBy     string
	Page       int
	PageSize   int
}

// ListResponse contains the paginated list of generations.
type ListResponse struct {
	Items      []storage.Generation `json:"items"`
	Total      int                  `json:"total"`
	Page       int                  `json:"page"`
	PageSize   int                  `json:"pageSize"`
	TotalPages int                  `json:"totalPages"`
}

// Service provides gallery operations.
type Service struct {
	repo        storage.Repository
	rateLimiter *ratelimit.Limiter
}

// NewService creates a new gallery service.
func NewService(repo storage.Repository, rateLimiter *ratelimit.Limiter) *Service {
	return &Service{
		repo:        repo,
		rateLimiter: rateLimiter,
	}
}

// ListGenerations retrieves a paginated list of generations with optional filtering.
func (s *Service) ListGenerations(ctx context.Context, req ListRequest) (*ListResponse, error) {
	// Validate and normalize inputs
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = DefaultPageSize
	}
	if req.PageSize > MaxPageSize {
		req.PageSize = MaxPageSize
	}

	// Validate sort option
	if req.SortBy == "" {
		req.SortBy = "newest"
	}
	if !ValidSortOptions[req.SortBy] {
		return nil, ErrInvalidSort
	}

	// Build filter for repository
	filter := storage.ListFilter{
		CategoryID: req.CategoryID,
		SortBy:     req.SortBy,
		Page:       req.Page,
		PageSize:   req.PageSize,
	}

	// Fetch from repository
	items, total, err := s.repo.ListGenerations(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(req.PageSize)))
	if totalPages < 1 {
		totalPages = 1
	}

	return &ListResponse{
		Items:      items,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// GetGeneration retrieves a single generation by ID and increments view count.
func (s *Service) GetGeneration(ctx context.Context, id string) (*storage.Generation, error) {
	if id == "" {
		return nil, ErrInvalidInput
	}

	// Get the generation
	gen, err := s.repo.GetGeneration(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	// Increment view count (fire and forget - don't fail if this fails)
	_ = s.repo.IncrementViewCount(ctx, id)

	return gen, nil
}

// RateGeneration submits or updates a rating for a generation.
// Returns the retry-after duration if rate limited.
func (s *Service) RateGeneration(ctx context.Context, genID string, score int, voterHash string, clientIP string) (retryAfter int, err error) {
	// Validate inputs
	if genID == "" || voterHash == "" {
		return 0, ErrInvalidInput
	}
	if score < 1 || score > 5 {
		return 0, ErrInvalidRating
	}

	// Check rate limit if limiter is configured
	if s.rateLimiter != nil {
		allowed, duration := s.rateLimiter.Allow(clientIP)
		if !allowed {
			return int(duration.Seconds()), ErrRateLimited
		}
	}

	// Verify generation exists
	_, err = s.repo.GetGeneration(ctx, genID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return 0, ErrNotFound
		}
		return 0, err
	}

	// Create or update rating
	err = s.repo.CreateOrUpdateRating(ctx, genID, score, voterHash)
	if err != nil {
		return 0, err
	}

	return 0, nil
}

// GetUserRating retrieves the user's rating for a generation.
// Returns 0 if the user hasn't rated the generation.
func (s *Service) GetUserRating(ctx context.Context, genID string, voterHash string) (int, error) {
	if genID == "" || voterHash == "" {
		return 0, ErrInvalidInput
	}

	return s.repo.GetUserRating(ctx, genID, voterHash)
}

// GetCategories retrieves all available categories.
func (s *Service) GetCategories(ctx context.Context) ([]storage.Category, error) {
	return s.repo.GetCategories(ctx)
}

// CalculateTotalPages is a helper function to calculate total pages.
// Exported for use in property tests.
func CalculateTotalPages(total, pageSize int) int {
	if pageSize <= 0 {
		pageSize = DefaultPageSize
	}
	if total <= 0 {
		return 1
	}
	return int(math.Ceil(float64(total) / float64(pageSize)))
}

// NormalizePageSize ensures page size is within valid bounds.
// Exported for use in property tests.
func NormalizePageSize(pageSize int) int {
	if pageSize < 1 {
		return DefaultPageSize
	}
	if pageSize > MaxPageSize {
		return MaxPageSize
	}
	return pageSize
}
