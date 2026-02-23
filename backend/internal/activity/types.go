// Package activity tracks user activity events for the media collection.
package activity

import (
	"time"

	"github.com/google/uuid"
)

// EventType identifies what kind of activity occurred.
type EventType string

const (
	EventItemAdded     EventType = "item_added"
	EventItemUpdated   EventType = "item_updated"
	EventItemDeleted   EventType = "item_deleted"
	EventStatusChanged EventType = "status_changed"
	EventRatingUpdated EventType = "rating_updated"
)

// Event represents a single activity event.
type Event struct {
	ID          uuid.UUID      `json:"id"`
	UserID      uuid.UUID      `json:"user_id"`
	MediaItemID *uuid.UUID     `json:"media_item_id,omitempty"`
	EventType   EventType      `json:"event_type"`
	Payload     map[string]any `json:"payload"`
	CreatedAt   time.Time      `json:"created_at"`
}
