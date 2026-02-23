// Package metadata provides external API metadata enrichment.
package metadata

import "context"

// Result holds normalized metadata from external providers.
type Result struct {
	ExternalID  string         `json:"external_id"`
	Title       string         `json:"title"`
	Creator     string         `json:"creator"`
	Genres      []string       `json:"genres"`
	CoverURL    string         `json:"cover_url"`
	ReleaseYear int            `json:"release_year"`
	Overview    string         `json:"overview"`
	Raw         map[string]any `json:"raw,omitempty"`
}

// Provider is the interface for metadata providers.
type Provider interface {
	Search(ctx context.Context, title string, year *int) ([]*Result, error)
	GetByID(ctx context.Context, id string) (*Result, error)
}
