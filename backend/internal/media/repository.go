package media

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository handles database operations for media items.
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new media Repository.
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// scanItem scans a database row into an Item.
func scanItem(row pgx.Row) (*Item, error) {
	var item Item
	var metaJSON []byte
	var genre []string

	err := row.Scan(
		&item.ID, &item.UserID, &item.Title, &item.MediaType,
		&item.Status, &item.Creator, &genre, &item.ReleaseYear,
		&item.CoverURL, &item.Notes, &item.Rating,
		&item.TMDBId, &item.MusicbrainzID, &item.IGDBId,
		&metaJSON, &item.CreatedAt, &item.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	item.Genre = genre
	if metaJSON != nil {
		if err := json.Unmarshal(metaJSON, &item.Metadata); err != nil {
			return nil, fmt.Errorf("unmarshal metadata: %w", err)
		}
	}
	if item.Metadata == nil {
		item.Metadata = map[string]any{}
	}
	return &item, nil
}

const itemColumns = `id, user_id, title, media_type, status, creator, genre,
	release_year, cover_url, notes, rating, tmdb_id, musicbrainz_id, igdb_id,
	metadata, created_at, updated_at`

// Create inserts a new media item.
func (r *Repository) Create(ctx context.Context, userID uuid.UUID, req CreateRequest, metaOverride map[string]any) (*Item, error) {
	meta := metaOverride
	if meta == nil {
		meta = map[string]any{}
	}
	metaJSON, err := json.Marshal(meta)
	if err != nil {
		return nil, fmt.Errorf("marshal metadata: %w", err)
	}

	genre := req.Genre
	if genre == nil {
		genre = []string{}
	}

	row := r.db.QueryRow(ctx, `
		INSERT INTO media_items (user_id, title, media_type, status, creator, genre,
			release_year, cover_url, notes, rating, metadata)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		RETURNING `+itemColumns,
		userID, req.Title, req.MediaType, req.Status, req.Creator, genre,
		req.ReleaseYear, req.CoverURL, req.Notes, req.Rating, metaJSON,
	)
	return scanItem(row)
}

// GetByID fetches a media item by ID.
func (r *Repository) GetByID(ctx context.Context, id, userID uuid.UUID) (*Item, error) {
	row := r.db.QueryRow(ctx,
		`SELECT `+itemColumns+` FROM media_items WHERE id=$1 AND user_id=$2`,
		id, userID,
	)
	item, err := scanItem(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("item not found")
		}
		return nil, fmt.Errorf("query item: %w", err)
	}
	return item, nil
}

// List returns paginated media items matching the filter.
func (r *Repository) List(ctx context.Context, f ListFilter) ([]*Item, int, error) {
	if f.PageSize <= 0 {
		f.PageSize = 20
	}
	if f.Page <= 0 {
		f.Page = 1
	}

	conditions := []string{"user_id = $1"}
	args := []any{f.UserID}
	argIdx := 2

	if f.MediaType != nil {
		conditions = append(conditions, fmt.Sprintf("media_type = $%d", argIdx))
		args = append(args, *f.MediaType)
		argIdx++
	}
	if f.Status != nil {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, *f.Status)
		argIdx++
	}
	if f.Genre != nil {
		conditions = append(conditions, fmt.Sprintf("$%d = ANY(genre)", argIdx))
		args = append(args, *f.Genre)
		argIdx++
	}

	where := "WHERE " + strings.Join(conditions, " AND ")
	countQuery := "SELECT COUNT(*) FROM media_items " + where
	query := fmt.Sprintf(
		`SELECT `+itemColumns+` FROM media_items %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d`,
		where, argIdx, argIdx+1,
	)
	countArgs := make([]any, argIdx-1)
	copy(countArgs, args[:argIdx-1])
	args = append(args, f.PageSize, (f.Page-1)*f.PageSize)

	var total int
	if err := r.db.QueryRow(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count media: %w", err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list media: %w", err)
	}
	defer rows.Close()

	items := make([]*Item, 0)
	for rows.Next() {
		item, err := scanItem(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("scan item: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	return items, total, nil
}

// Update modifies a media item's fields.
func (r *Repository) Update(ctx context.Context, id, userID uuid.UUID, req UpdateRequest) (*Item, error) {
	sets := []string{}
	args := []any{}
	argIdx := 1

	if req.Title != nil {
		sets = append(sets, fmt.Sprintf("title=$%d", argIdx))
		args = append(args, *req.Title)
		argIdx++
	}
	if req.Status != nil {
		sets = append(sets, fmt.Sprintf("status=$%d", argIdx))
		args = append(args, *req.Status)
		argIdx++
	}
	if req.Creator != nil {
		sets = append(sets, fmt.Sprintf("creator=$%d", argIdx))
		args = append(args, *req.Creator)
		argIdx++
	}
	if req.Genre != nil {
		sets = append(sets, fmt.Sprintf("genre=$%d", argIdx))
		args = append(args, req.Genre)
		argIdx++
	}
	if req.ReleaseYear != nil {
		sets = append(sets, fmt.Sprintf("release_year=$%d", argIdx))
		args = append(args, *req.ReleaseYear)
		argIdx++
	}
	if req.CoverURL != nil {
		sets = append(sets, fmt.Sprintf("cover_url=$%d", argIdx))
		args = append(args, *req.CoverURL)
		argIdx++
	}
	if req.Notes != nil {
		sets = append(sets, fmt.Sprintf("notes=$%d", argIdx))
		args = append(args, *req.Notes)
		argIdx++
	}
	if req.Rating != nil {
		sets = append(sets, fmt.Sprintf("rating=$%d", argIdx))
		args = append(args, *req.Rating)
		argIdx++
	}

	if len(sets) == 0 {
		return r.GetByID(ctx, id, userID)
	}

	args = append(args, id, userID)
	query := fmt.Sprintf(
		`UPDATE media_items SET %s WHERE id=$%d AND user_id=$%d RETURNING `+itemColumns,
		strings.Join(sets, ","), argIdx, argIdx+1,
	)

	row := r.db.QueryRow(ctx, query, args...)
	item, err := scanItem(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("item not found")
		}
		return nil, fmt.Errorf("update item: %w", err)
	}
	return item, nil
}

// Delete removes a media item.
func (r *Repository) Delete(ctx context.Context, id, userID uuid.UUID) error {
	result, err := r.db.Exec(ctx,
		"DELETE FROM media_items WHERE id=$1 AND user_id=$2", id, userID,
	)
	if err != nil {
		return fmt.Errorf("delete item: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("item not found")
	}
	return nil
}

// UpdateStatus patches only the status field.
func (r *Repository) UpdateStatus(ctx context.Context, id, userID uuid.UUID, status Status) (*Item, error) {
	row := r.db.QueryRow(ctx, `
		UPDATE media_items SET status=$1 WHERE id=$2 AND user_id=$3 RETURNING `+itemColumns,
		status, id, userID,
	)
	item, err := scanItem(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("item not found")
		}
		return nil, fmt.Errorf("update status: %w", err)
	}
	return item, nil
}

// Search performs a full-text search using tsvector + trigram fallback.
func (r *Repository) Search(ctx context.Context, userID uuid.UUID, query string, mediaType *MediaType, page, pageSize int) ([]*Item, int, error) {
	if pageSize <= 0 {
		pageSize = 20
	}
	if page <= 0 {
		page = 1
	}

	baseArgs := []any{userID, query}
	typeFilter := ""
	if mediaType != nil {
		typeFilter = " AND media_type=$3"
		baseArgs = append(baseArgs, *mediaType)
	}

	baseQuery := fmt.Sprintf(`
		FROM media_items
		WHERE user_id=$1%s
		AND (
			search_vector @@ plainto_tsquery('english', $2)
			OR title ILIKE '%%' || $2 || '%%'
			OR creator ILIKE '%%' || $2 || '%%'
		)`, typeFilter)

	var total int
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) "+baseQuery, baseArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count search: %w", err)
	}

	limitIdx := len(baseArgs) + 1
	offsetIdx := limitIdx + 1
	selectQuery := fmt.Sprintf(
		`SELECT `+itemColumns+` %s ORDER BY ts_rank(search_vector, plainto_tsquery('english', $2)) DESC LIMIT $%d OFFSET $%d`,
		baseQuery, limitIdx, offsetIdx,
	)
	allArgs := append(baseArgs, pageSize, (page-1)*pageSize)

	rows, err := r.db.Query(ctx, selectQuery, allArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("search media: %w", err)
	}
	defer rows.Close()

	items := make([]*Item, 0)
	for rows.Next() {
		item, err := scanItem(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("scan item: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	return items, total, nil
}

// GetAllForUser returns all items for a user (used for AI features).
func (r *Repository) GetAllForUser(ctx context.Context, userID uuid.UUID) ([]*Item, error) {
	rows, err := r.db.Query(ctx,
		`SELECT `+itemColumns+` FROM media_items WHERE user_id=$1 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("query all items: %w", err)
	}
	defer rows.Close()

	var items []*Item
	for rows.Next() {
		item, err := scanItem(rows)
		if err != nil {
			return nil, fmt.Errorf("scan item: %w", err)
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
