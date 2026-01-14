// Package gallery provides the gallery service for browsing and rating generations.
package gallery

import (
	"context"
	"errors"
	"log/slog"
	"math"
	"time"

	"better-kiro-prompts/internal/config"
	"better-kiro-prompts/internal/logger"
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
	log         *slog.Logger
	pageSize    int
	defaultSort string
}

// NewService creates a new gallery service with default configuration.
func NewService(repo storage.Repository, rateLimiter *ratelimit.Limiter, log *logger.Logger) *Service {
	// Use defaults from config package
	defaultCfg := config.DefaultConfig()
	return NewServiceWithConfig(repo, rateLimiter, log, defaultCfg.Gallery)
}

// NewServiceWithConfig creates a new gallery service with the provided configuration.
func NewServiceWithConfig(repo storage.Repository, rateLimiter *ratelimit.Limiter, log *logger.Logger, cfg config.GalleryConfig) *Service {
	var slogger *slog.Logger
	if log != nil {
		slogger = log.App()
	}
	return &Service{
		repo:        repo,
		rateLimiter: rateLimiter,
		log:         slogger,
		pageSize:    cfg.PageSize,
		defaultSort: cfg.DefaultSort,
	}
}

// ListGenerations retrieves a paginated list of generations with optional filtering.
func (s *Service) ListGenerations(ctx context.Context, req ListRequest) (*ListResponse, error) {
	requestID := logger.GetRequestID(ctx)
	start := time.Now()

	// Log start
	if s.log != nil {
		s.log.Info("gallery_list_start",
			slog.String("request_id", requestID),
			slog.String("sort_by", req.SortBy),
			slog.Int("page", req.Page),
			slog.Int("page_size", req.PageSize),
			slog.Any("category_id", req.CategoryID),
		)
	}

	// Validate and normalize inputs
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = s.pageSize
	}
	if req.PageSize > MaxPageSize {
		req.PageSize = MaxPageSize
	}

	// Validate sort option
	if req.SortBy == "" {
		req.SortBy = s.defaultSort
	}
	if !ValidSortOptions[req.SortBy] {
		if s.log != nil {
			s.log.Warn("gallery_list_invalid_sort",
				slog.String("request_id", requestID),
				slog.String("sort_by", req.SortBy),
			)
		}
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
		if s.log != nil {
			s.log.Error("gallery_list_failed",
				slog.String("request_id", requestID),
				slog.String("error", err.Error()),
				slog.Duration("duration", time.Since(start)),
			)
		}
		return nil, err
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(req.PageSize)))
	if totalPages < 1 {
		totalPages = 1
	}

	// Log completion
	if s.log != nil {
		s.log.Info("gallery_list_complete",
			slog.String("request_id", requestID),
			slog.Int("item_count", len(items)),
			slog.Int("total", total),
			slog.Duration("duration", time.Since(start)),
		)
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
// Deprecated: Use GetGenerationWithView for IP-deduplicated view tracking.
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

// GetGenerationWithView retrieves a single generation by ID and records a view
// deduplicated by IP hash. Only increments view count for new unique views.
func (s *Service) GetGenerationWithView(ctx context.Context, id string, ipHash string) (*storage.Generation, error) {
	requestID := logger.GetRequestID(ctx)
	start := time.Now()

	// Log start
	if s.log != nil {
		s.log.Info("gallery_get_start",
			slog.String("request_id", requestID),
			slog.String("generation_id", id),
		)
	}

	if id == "" {
		if s.log != nil {
			s.log.Warn("gallery_get_invalid_input",
				slog.String("request_id", requestID),
				slog.String("error", "empty generation id"),
			)
		}
		return nil, ErrInvalidInput
	}

	// Get the generation
	gen, err := s.repo.GetGeneration(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			if s.log != nil {
				s.log.Warn("gallery_get_not_found",
					slog.String("request_id", requestID),
					slog.String("generation_id", id),
				)
			}
			return nil, ErrNotFound
		}
		if s.log != nil {
			s.log.Error("gallery_get_failed",
				slog.String("request_id", requestID),
				slog.String("generation_id", id),
				slog.String("error", err.Error()),
			)
		}
		return nil, err
	}

	// Record view with IP deduplication (fire and forget - don't fail if this fails)
	newView := false
	if ipHash != "" {
		newView, _ = s.repo.RecordView(ctx, id, ipHash)
		if s.log != nil {
			s.log.Debug("gallery_view_recorded",
				slog.String("request_id", requestID),
				slog.String("generation_id", id),
				slog.Bool("new_view", newView),
			)
		}
	}

	// Log completion
	if s.log != nil {
		s.log.Info("gallery_get_complete",
			slog.String("request_id", requestID),
			slog.String("generation_id", id),
			slog.Bool("new_view", newView),
			slog.Duration("duration", time.Since(start)),
		)
	}

	return gen, nil
}

// RateGeneration submits or updates a rating for a generation.
// Returns the retry-after duration if rate limited.
func (s *Service) RateGeneration(ctx context.Context, genID string, score int, voterHash string, clientIP string) (retryAfter int, err error) {
	requestID := logger.GetRequestID(ctx)
	start := time.Now()

	// Log start
	if s.log != nil {
		s.log.Info("gallery_rate_start",
			slog.String("request_id", requestID),
			slog.String("generation_id", genID),
			slog.Int("score", score),
		)
	}

	// Validate inputs
	if genID == "" || voterHash == "" {
		if s.log != nil {
			s.log.Warn("gallery_rate_invalid_input",
				slog.String("request_id", requestID),
				slog.String("error", "empty generation id or voter hash"),
			)
		}
		return 0, ErrInvalidInput
	}
	if score < 1 || score > 5 {
		if s.log != nil {
			s.log.Warn("gallery_rate_invalid_score",
				slog.String("request_id", requestID),
				slog.Int("score", score),
			)
		}
		return 0, ErrInvalidRating
	}

	// Check rate limit if limiter is configured
	if s.rateLimiter != nil {
		allowed, duration := s.rateLimiter.Allow(clientIP)
		if !allowed {
			if s.log != nil {
				s.log.Warn("gallery_rate_limited",
					slog.String("request_id", requestID),
					slog.String("generation_id", genID),
					slog.Duration("retry_after", duration),
				)
			}
			return int(duration.Seconds()), ErrRateLimited
		}
		if s.log != nil {
			s.log.Debug("gallery_rate_limit_allowed",
				slog.String("request_id", requestID),
			)
		}
	}

	// Verify generation exists
	_, err = s.repo.GetGeneration(ctx, genID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			if s.log != nil {
				s.log.Warn("gallery_rate_not_found",
					slog.String("request_id", requestID),
					slog.String("generation_id", genID),
				)
			}
			return 0, ErrNotFound
		}
		if s.log != nil {
			s.log.Error("gallery_rate_get_failed",
				slog.String("request_id", requestID),
				slog.String("generation_id", genID),
				slog.String("error", err.Error()),
			)
		}
		return 0, err
	}

	// Create or update rating
	err = s.repo.CreateOrUpdateRating(ctx, genID, score, voterHash)
	if err != nil {
		if s.log != nil {
			s.log.Error("gallery_rate_failed",
				slog.String("request_id", requestID),
				slog.String("generation_id", genID),
				slog.String("error", err.Error()),
			)
		}
		return 0, err
	}

	// Log completion
	if s.log != nil {
		s.log.Info("gallery_rate_complete",
			slog.String("request_id", requestID),
			slog.String("generation_id", genID),
			slog.Int("score", score),
			slog.Duration("duration", time.Since(start)),
		)
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
	defaultCfg := config.DefaultConfig()
	if pageSize <= 0 {
		pageSize = defaultCfg.Gallery.PageSize
	}
	if total <= 0 {
		return 1
	}
	return int(math.Ceil(float64(total) / float64(pageSize)))
}

// NormalizePageSize ensures page size is within valid bounds.
// Exported for use in property tests.
func NormalizePageSize(pageSize int) int {
	defaultCfg := config.DefaultConfig()
	if pageSize < 1 {
		return defaultCfg.Gallery.PageSize
	}
	if pageSize > MaxPageSize {
		return MaxPageSize
	}
	return pageSize
}
