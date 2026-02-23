package media

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/your-org/ems/internal/auth"
	"github.com/your-org/ems/internal/httputil"
)

// Handler handles HTTP requests for media endpoints.
type Handler struct {
	svc *Service
}

// NewHandler creates a new media Handler.
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// List handles GET /api/media.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromCtx(r.Context())

	f := ListFilter{
		UserID:   claims.UserID,
		Page:     queryInt(r, "page", 1),
		PageSize: queryInt(r, "page_size", 20),
	}

	if t := r.URL.Query().Get("type"); t != "" {
		mt := MediaType(t)
		f.MediaType = &mt
	}
	if s := r.URL.Query().Get("status"); s != "" {
		st := Status(s)
		f.Status = &st
	}
	if g := r.URL.Query().Get("genre"); g != "" {
		f.Genre = &g
	}

	items, total, err := h.svc.List(r.Context(), f)
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputil.WriteJSON(w, http.StatusOK, map[string]any{
		"items": items,
		"total": total,
		"page":  f.Page,
	})
}

// Create handles POST /api/media.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromCtx(r.Context())

	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Title == "" {
		httputil.WriteError(w, http.StatusBadRequest, "title is required")
		return
	}
	if req.Status == "" {
		req.Status = StatusOwned
	}

	item, err := h.svc.Create(r.Context(), claims.UserID, req)
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputil.WriteJSON(w, http.StatusCreated, item)
}

// Get handles GET /api/media/:id.
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromCtx(r.Context())
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	item, err := h.svc.GetByID(r.Context(), id, claims.UserID)
	if err != nil {
		httputil.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	httputil.WriteJSON(w, http.StatusOK, item)
}

// Update handles PUT /api/media/:id.
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromCtx(r.Context())
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	item, err := h.svc.Update(r.Context(), id, claims.UserID, req)
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputil.WriteJSON(w, http.StatusOK, item)
}

// Delete handles DELETE /api/media/:id.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromCtx(r.Context())
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.svc.Delete(r.Context(), id, claims.UserID); err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdateStatus handles PATCH /api/media/:id/status.
func (h *Handler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromCtx(r.Context())
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req StatusUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	item, err := h.svc.UpdateStatus(r.Context(), id, claims.UserID, req.Status)
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputil.WriteJSON(w, http.StatusOK, item)
}

func queryInt(r *http.Request, key string, defaultVal int) int {
	v := r.URL.Query().Get(key)
	if v == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return defaultVal
	}
	return n
}
