package activity

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository handles activity event persistence.
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new activity Repository.
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// Record inserts a new activity event.
func (r *Repository) Record(ctx context.Context, userID uuid.UUID, mediaItemID *uuid.UUID, eventType EventType, payload map[string]any) error {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	_, err = r.db.Exec(ctx, `
		INSERT INTO activity_events (user_id, media_item_id, event_type, payload)
		VALUES ($1, $2, $3, $4)
	`, userID, mediaItemID, eventType, payloadJSON)
	if err != nil {
		return fmt.Errorf("insert activity: %w", err)
	}
	return nil
}

// List returns recent activity events for a user.
func (r *Repository) List(ctx context.Context, userID uuid.UUID, limit int) ([]*Event, error) {
	if limit <= 0 {
		limit = 20
	}

	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, media_item_id, event_type, payload, created_at
		FROM activity_events
		WHERE user_id=$1
		ORDER BY created_at DESC
		LIMIT $2
	`, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("list activity: %w", err)
	}
	defer rows.Close()

	var events []*Event
	for rows.Next() {
		var e Event
		var payloadJSON []byte
		if err := rows.Scan(&e.ID, &e.UserID, &e.MediaItemID, &e.EventType, &payloadJSON, &e.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan event: %w", err)
		}
		if err := json.Unmarshal(payloadJSON, &e.Payload); err != nil {
			return nil, fmt.Errorf("unmarshal payload: %w", err)
		}
		events = append(events, &e)
	}
	return events, rows.Err()
}
