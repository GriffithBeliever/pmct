package metadata

import (
	"encoding/json"
	"net/http"

	"github.com/your-org/ems/internal/httputil"
	"github.com/your-org/ems/internal/media"
)

// Handler handles HTTP requests for metadata search.
type Handler struct {
	svc *Service
}

// NewHandler creates a new metadata Handler.
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// Search handles POST /api/metadata/search.
func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title     string          `json:"title"`
		MediaType media.MediaType `json:"media_type"`
		Year      *int            `json:"year,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Title == "" {
		httputil.WriteError(w, http.StatusBadRequest, "title is required")
		return
	}

	results, err := h.svc.Search(r.Context(), req.Title, req.MediaType, req.Year)
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputil.WriteJSON(w, http.StatusOK, results)
}
