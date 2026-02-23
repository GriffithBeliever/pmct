// Package search provides full-text search over the media collection.
package search

import (
	"net/http"
	"strconv"

	"github.com/your-org/ems/internal/auth"
	"github.com/your-org/ems/internal/httputil"
	"github.com/your-org/ems/internal/media"
)

// Handler handles HTTP requests for search.
type Handler struct {
	repo *media.Repository
}

// NewHandler creates a new search Handler.
func NewHandler(repo *media.Repository) *Handler {
	return &Handler{repo: repo}
}

// Search handles GET /api/search.
func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromCtx(r.Context())

	q := r.URL.Query().Get("q")
	if q == "" {
		httputil.WriteError(w, http.StatusBadRequest, "query parameter 'q' is required")
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	var mediaType *media.MediaType
	if t := r.URL.Query().Get("type"); t != "" {
		mt := media.MediaType(t)
		mediaType = &mt
	}

	items, total, err := h.repo.Search(r.Context(), claims.UserID, q, mediaType, page, pageSize)
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputil.WriteJSON(w, http.StatusOK, map[string]any{
		"items": items,
		"total": total,
		"page":  page,
		"query": q,
	})
}
