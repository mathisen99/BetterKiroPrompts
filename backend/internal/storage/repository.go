// Package storage provides database storage for generations and ratings.
package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// Common errors
var (
	ErrNotFound      = errors.New("record not found")
	ErrDuplicateKey  = errors.New("duplicate key violation")
	ErrInvalidInput  = errors.New("invalid input")
	ErrDatabaseError = errors.New("database error")
)

// Generation represents a stored generation record.
type Generation struct {
	ID              string          `json:"id"`
	ProjectIdea     string          `json:"projectIdea"`
	ExperienceLevel string          `json:"experienceLevel"`
	HookPreset      string          `json:"hookPreset"`
	Files           json.RawMessage `json:"files"`
	CategoryID      int             `json:"categoryId"`
	CategoryName    string          `json:"categoryName,omitempty"`
	AvgRating       float64         `json:"avgRating"`
	RatingCount     int             `json:"ratingCount"`
	ViewCount       int             `json:"viewCount"`
	CreatedAt       time.Time       `json:"createdAt"`
}

// ListFilter defines filtering and pagination options for listing generations.
type ListFilter struct {
	CategoryID *int
	SortBy     string // "newest", "highest_rated", "most_viewed"
	Page       int
	PageSize   int
}

// Repository defines the interface for storage operations.
type Repository interface {
	// Generations
	CreateGeneration(ctx context.Context, gen *Generation) error
	GetGeneration(ctx context.Context, id string) (*Generation, error)
	ListGenerations(ctx context.Context, filter ListFilter) ([]Generation, int, error)
	IncrementViewCount(ctx context.Context, id string) error

	// Ratings
	CreateOrUpdateRating(ctx context.Context, genID string, score int, voterHash string) error
	GetUserRating(ctx context.Context, genID string, voterHash string) (int, error)

	// Categories
	GetCategoryByKeywords(ctx context.Context, text string) (int, error)
	GetCategories(ctx context.Context) ([]Category, error)
}

// Category represents a generation category.
type Category struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Keywords []string `json:"keywords"`
}

// PostgresRepository implements Repository using PostgreSQL.
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository creates a new PostgreSQL repository.
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// CreateGeneration stores a new generation in the database.
func (r *PostgresRepository) CreateGeneration(ctx context.Context, gen *Generation) error {
	if gen == nil {
		return ErrInvalidInput
	}

	query := `
		INSERT INTO generations (project_idea, experience_level, hook_preset, files, category_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`

	err := r.db.QueryRowContext(ctx, query,
		gen.ProjectIdea,
		gen.ExperienceLevel,
		gen.HookPreset,
		gen.Files,
		gen.CategoryID,
	).Scan(&gen.ID, &gen.CreatedAt)

	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return nil
}

// GetGeneration retrieves a generation by ID.
func (r *PostgresRepository) GetGeneration(ctx context.Context, id string) (*Generation, error) {
	query := `
		SELECT g.id, g.project_idea, g.experience_level, g.hook_preset, g.files,
		       g.category_id, c.name, g.avg_rating, g.rating_count, g.view_count, g.created_at
		FROM generations g
		LEFT JOIN categories c ON g.category_id = c.id
		WHERE g.id = $1`

	gen := &Generation{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&gen.ID,
		&gen.ProjectIdea,
		&gen.ExperienceLevel,
		&gen.HookPreset,
		&gen.Files,
		&gen.CategoryID,
		&gen.CategoryName,
		&gen.AvgRating,
		&gen.RatingCount,
		&gen.ViewCount,
		&gen.CreatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return gen, nil
}

// ListGenerations retrieves a paginated list of generations with optional filtering.
func (r *PostgresRepository) ListGenerations(ctx context.Context, filter ListFilter) ([]Generation, int, error) {
	// Set defaults
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 || filter.PageSize > 100 {
		filter.PageSize = 20
	}

	// Build query with optional category filter
	baseQuery := `
		FROM generations g
		LEFT JOIN categories c ON g.category_id = c.id`

	whereClause := ""
	args := []interface{}{}
	argIndex := 1

	if filter.CategoryID != nil {
		whereClause = fmt.Sprintf(" WHERE g.category_id = $%d", argIndex)
		args = append(args, *filter.CategoryID)
		argIndex++
	}

	// Count total
	countQuery := "SELECT COUNT(*)" + baseQuery + whereClause
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	// Determine sort order
	orderBy := " ORDER BY g.created_at DESC" // default: newest
	switch filter.SortBy {
	case "highest_rated":
		orderBy = " ORDER BY g.avg_rating DESC, g.rating_count DESC"
	case "most_viewed":
		orderBy = " ORDER BY g.view_count DESC"
	}

	// Build select query with pagination
	offset := (filter.Page - 1) * filter.PageSize
	selectQuery := fmt.Sprintf(`
		SELECT g.id, g.project_idea, g.experience_level, g.hook_preset, g.files,
		       g.category_id, c.name, g.avg_rating, g.rating_count, g.view_count, g.created_at
		%s%s%s
		LIMIT $%d OFFSET $%d`,
		baseQuery, whereClause, orderBy, argIndex, argIndex+1)

	args = append(args, filter.PageSize, offset)

	rows, err := r.db.QueryContext(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer func() { _ = rows.Close() }()

	generations := []Generation{}
	for rows.Next() {
		var gen Generation
		if err := rows.Scan(
			&gen.ID,
			&gen.ProjectIdea,
			&gen.ExperienceLevel,
			&gen.HookPreset,
			&gen.Files,
			&gen.CategoryID,
			&gen.CategoryName,
			&gen.AvgRating,
			&gen.RatingCount,
			&gen.ViewCount,
			&gen.CreatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}
		generations = append(generations, gen)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return generations, total, nil
}

// IncrementViewCount increments the view count for a generation.
func (r *PostgresRepository) IncrementViewCount(ctx context.Context, id string) error {
	query := `UPDATE generations SET view_count = view_count + 1 WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

// CreateOrUpdateRating creates or updates a rating for a generation.
func (r *PostgresRepository) CreateOrUpdateRating(ctx context.Context, genID string, score int, voterHash string) error {
	if score < 1 || score > 5 {
		return fmt.Errorf("%w: score must be between 1 and 5", ErrInvalidInput)
	}

	// Use upsert to handle both create and update
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer func() { _ = tx.Rollback() }()

	// Upsert the rating
	upsertQuery := `
		INSERT INTO ratings (generation_id, score, voter_hash)
		VALUES ($1, $2, $3)
		ON CONFLICT (generation_id, voter_hash)
		DO UPDATE SET score = $2, created_at = NOW()`

	_, err = tx.ExecContext(ctx, upsertQuery, genID, score, voterHash)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	// Recalculate average rating
	updateAvgQuery := `
		UPDATE generations
		SET avg_rating = (SELECT COALESCE(AVG(score), 0) FROM ratings WHERE generation_id = $1),
		    rating_count = (SELECT COUNT(*) FROM ratings WHERE generation_id = $1)
		WHERE id = $1`

	_, err = tx.ExecContext(ctx, updateAvgQuery, genID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return nil
}

// GetUserRating retrieves the user's rating for a generation.
func (r *PostgresRepository) GetUserRating(ctx context.Context, genID string, voterHash string) (int, error) {
	query := `SELECT score FROM ratings WHERE generation_id = $1 AND voter_hash = $2`

	var score int
	err := r.db.QueryRowContext(ctx, query, genID, voterHash).Scan(&score)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil // No rating yet
	}
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return score, nil
}

// GetCategoryByKeywords is implemented in category.go

// GetCategories retrieves all categories.
func (r *PostgresRepository) GetCategories(ctx context.Context) ([]Category, error) {
	query := `SELECT id, name, keywords FROM categories ORDER BY id`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer func() { _ = rows.Close() }()

	categories := []Category{}
	for rows.Next() {
		var cat Category
		var keywords []byte
		if err := rows.Scan(&cat.ID, &cat.Name, &keywords); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}
		// Parse PostgreSQL array format
		if err := parsePostgresArray(keywords, &cat.Keywords); err != nil {
			return nil, fmt.Errorf("%w: failed to parse keywords: %v", ErrDatabaseError, err)
		}
		categories = append(categories, cat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return categories, nil
}

// parsePostgresArray parses a PostgreSQL text array into a Go string slice.
func parsePostgresArray(data []byte, dest *[]string) error {
	str := string(data)
	if str == "{}" || str == "" {
		*dest = []string{}
		return nil
	}

	// Remove braces
	str = str[1 : len(str)-1]
	if str == "" {
		*dest = []string{}
		return nil
	}

	// Split by comma (simple case without quoted strings)
	parts := []string{}
	current := ""
	inQuotes := false

	for _, ch := range str {
		switch ch {
		case '"':
			inQuotes = !inQuotes
		case ',':
			if !inQuotes {
				parts = append(parts, current)
				current = ""
				continue
			}
			current += string(ch)
		default:
			current += string(ch)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}

	*dest = parts
	return nil
}
