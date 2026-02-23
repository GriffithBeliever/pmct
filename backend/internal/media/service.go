package media

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

// MetadataEnricher enriches media metadata by external provider.
type MetadataEnricher interface {
	Enrich(ctx context.Context, title string, mediaType MediaType, releaseYear *int) (map[string]any, error)
}

// Service orchestrates media operations with optional enrichment.
type Service struct {
	repo     *Repository
	enricher MetadataEnricher
}

// NewService creates a new media Service.
func NewService(repo *Repository, enricher MetadataEnricher) *Service {
	return &Service{repo: repo, enricher: enricher}
}

// Create adds a new media item, optionally enriching with external metadata.
func (s *Service) Create(ctx context.Context, userID uuid.UUID, req CreateRequest) (*Item, error) {
	var metaOverride map[string]any

	if req.EnrichMetadata && s.enricher != nil {
		g, gCtx := errgroup.WithContext(ctx)
		var enriched map[string]any

		g.Go(func() error {
			var err error
			enriched, err = s.enricher.Enrich(gCtx, req.Title, req.MediaType, req.ReleaseYear)
			if err != nil {
				// Non-fatal: log and continue without metadata
				slog.Warn("metadata enrichment failed", "title", req.Title, "error", err)
			}
			return nil
		})

		if err := g.Wait(); err != nil {
			return nil, fmt.Errorf("create with enrichment: %w", err)
		}

		if enriched != nil {
			metaOverride = enriched
			if req.CoverURL == "" {
				if url, ok := enriched["cover_url"].(string); ok {
					req.CoverURL = url
				}
			}
			if req.Creator == "" {
				if creator, ok := enriched["creator"].(string); ok {
					req.Creator = creator
				}
			}
			if len(req.Genre) == 0 {
				if genres, ok := enriched["genres"].([]string); ok {
					req.Genre = genres
				}
			}
		}
	}

	item, err := s.repo.Create(ctx, userID, req, metaOverride)
	if err != nil {
		return nil, fmt.Errorf("create item: %w", err)
	}
	return item, nil
}

// GetByID returns a single item owned by userID.
func (s *Service) GetByID(ctx context.Context, id, userID uuid.UUID) (*Item, error) {
	return s.repo.GetByID(ctx, id, userID)
}

// List returns paginated items matching the filter.
func (s *Service) List(ctx context.Context, f ListFilter) ([]*Item, int, error) {
	return s.repo.List(ctx, f)
}

// Update modifies an existing item.
func (s *Service) Update(ctx context.Context, id, userID uuid.UUID, req UpdateRequest) (*Item, error) {
	return s.repo.Update(ctx, id, userID, req)
}

// Delete removes an item.
func (s *Service) Delete(ctx context.Context, id, userID uuid.UUID) error {
	return s.repo.Delete(ctx, id, userID)
}

// UpdateStatus changes just the status of an item.
func (s *Service) UpdateStatus(ctx context.Context, id, userID uuid.UUID, status Status) (*Item, error) {
	return s.repo.UpdateStatus(ctx, id, userID, status)
}

// GetAllForUser returns all items for AI features.
func (s *Service) GetAllForUser(ctx context.Context, userID uuid.UUID) ([]*Item, error) {
	return s.repo.GetAllForUser(ctx, userID)
}
