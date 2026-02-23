package metadata

import (
	"context"
	"fmt"

	"github.com/your-org/ems/internal/media"
)

// Service dispatches metadata searches by media type.
type Service struct {
	tmdb        *TMDBClient
	musicbrainz *MusicBrainzClient
	igdb        *IGDBClient
}

// NewService creates a new metadata Service.
func NewService(tmdb *TMDBClient, mb *MusicBrainzClient, igdb *IGDBClient) *Service {
	return &Service{tmdb: tmdb, musicbrainz: mb, igdb: igdb}
}

// Search dispatches to the correct provider based on media type.
func (s *Service) Search(ctx context.Context, title string, mediaType media.MediaType, year *int) ([]*Result, error) {
	switch mediaType {
	case media.MediaTypeMovie:
		return s.tmdb.Search(ctx, title, year)
	case media.MediaTypeMusic:
		return s.musicbrainz.Search(ctx, title, year)
	case media.MediaTypeGame:
		return s.igdb.Search(ctx, title, year)
	default:
		return nil, fmt.Errorf("unknown media type: %s", mediaType)
	}
}

// Enrich fetches the best metadata match and returns it as a map.
func (s *Service) Enrich(ctx context.Context, title string, mediaType media.MediaType, releaseYear *int) (map[string]any, error) {
	results, err := s.Search(ctx, title, mediaType, releaseYear)
	if err != nil {
		return nil, fmt.Errorf("metadata search: %w", err)
	}
	if len(results) == 0 {
		return nil, nil
	}

	r := results[0]
	return map[string]any{
		"external_id":  r.ExternalID,
		"cover_url":    r.CoverURL,
		"creator":      r.Creator,
		"genres":       r.Genres,
		"release_year": r.ReleaseYear,
		"overview":     r.Overview,
		"source":       string(mediaType),
	}, nil
}
