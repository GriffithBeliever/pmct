// Package media provides media item domain types, repository, service and handler.
package media

import (
	"time"

	"github.com/google/uuid"
)

// MediaType represents the type of media.
type MediaType string

const (
	MediaTypeMovie MediaType = "movie"
	MediaTypeMusic MediaType = "music"
	MediaTypeGame  MediaType = "game"
)

// Status represents the user's ownership/usage status for an item.
type Status string

const (
	StatusOwned          Status = "owned"
	StatusWishlist       Status = "wishlist"
	StatusCurrentlyUsing Status = "currently_using"
	StatusCompleted      Status = "completed"
)

// Item represents a media item in the collection.
type Item struct {
	ID            uuid.UUID      `json:"id"`
	UserID        uuid.UUID      `json:"user_id"`
	Title         string         `json:"title"`
	MediaType     MediaType      `json:"media_type"`
	Status        Status         `json:"status"`
	Creator       string         `json:"creator"`
	Genre         []string       `json:"genre"`
	ReleaseYear   *int           `json:"release_year,omitempty"`
	CoverURL      string         `json:"cover_url"`
	Notes         string         `json:"notes"`
	Rating        *float64       `json:"rating,omitempty"`
	TMDBId        *string        `json:"tmdb_id,omitempty"`
	MusicbrainzID *string        `json:"musicbrainz_id,omitempty"`
	IGDBId        *string        `json:"igdb_id,omitempty"`
	Metadata      map[string]any `json:"metadata"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

// CreateRequest is the payload for creating a new media item.
type CreateRequest struct {
	Title          string    `json:"title"`
	MediaType      MediaType `json:"media_type"`
	Status         Status    `json:"status"`
	Creator        string    `json:"creator"`
	Genre          []string  `json:"genre"`
	ReleaseYear    *int      `json:"release_year,omitempty"`
	CoverURL       string    `json:"cover_url"`
	Notes          string    `json:"notes"`
	Rating         *float64  `json:"rating,omitempty"`
	EnrichMetadata bool      `json:"enrich_metadata"`
}

// UpdateRequest is the payload for updating a media item.
type UpdateRequest struct {
	Title       *string  `json:"title,omitempty"`
	Status      *Status  `json:"status,omitempty"`
	Creator     *string  `json:"creator,omitempty"`
	Genre       []string `json:"genre,omitempty"`
	ReleaseYear *int     `json:"release_year,omitempty"`
	CoverURL    *string  `json:"cover_url,omitempty"`
	Notes       *string  `json:"notes,omitempty"`
	Rating      *float64 `json:"rating,omitempty"`
}

// StatusUpdateRequest is the payload for patching just the status.
type StatusUpdateRequest struct {
	Status Status `json:"status"`
}

// ListFilter holds query parameters for listing media items.
type ListFilter struct {
	UserID    uuid.UUID
	MediaType *MediaType
	Status    *Status
	Genre     *string
	Page      int
	PageSize  int
}
