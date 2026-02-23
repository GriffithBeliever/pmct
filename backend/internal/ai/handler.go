package ai

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/your-org/ems/internal/auth"
	"github.com/your-org/ems/internal/httputil"
	"github.com/your-org/ems/internal/media"
)

// Handler handles HTTP requests for AI endpoints.
type Handler struct {
	svc      *Service
	mediaSvc *media.Service
}

// NewHandler creates a new AI Handler.
func NewHandler(svc *Service, mediaSvc *media.Service) *Handler {
	return &Handler{svc: svc, mediaSvc: mediaSvc}
}

// Recommendations handles GET /api/ai/recommendations.
func (h *Handler) Recommendations(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromCtx(r.Context())
	items, err := h.mediaSvc.GetAllForUser(r.Context(), claims.UserID)
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	recs, err := h.svc.Recommend(r.Context(), claims.UserID, items)
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.WriteJSON(w, http.StatusOK, recs)
}

// Insights handles GET /api/ai/insights with SSE streaming.
func (h *Handler) Insights(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromCtx(r.Context())
	items, err := h.mediaSvc.GetAllForUser(r.Context(), claims.UserID)
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	flusher, ok := w.(http.Flusher)
	if !ok {
		httputil.WriteError(w, http.StatusInternalServerError, "streaming not supported")
		return
	}

	out := make(chan string, 32)
	errCh := make(chan error, 1)

	go func() {
		errCh <- h.svc.StreamInsights(r.Context(), items, out)
		close(out)
	}()

	for {
		select {
		case token, ok := <-out:
			if !ok {
				fmt.Fprintf(w, "event: done\ndata: {}\n\n")
				flusher.Flush()
				return
			}
			fmt.Fprintf(w, "data: %s\n\n", jsonStr(token))
			flusher.Flush()
		case err := <-errCh:
			if err != nil {
				fmt.Fprintf(w, "event: error\ndata: %s\n\n", jsonStr(err.Error()))
				flusher.Flush()
			}
			return
		case <-r.Context().Done():
			return
		}
	}
}

// NLSearch handles POST /api/ai/nl-search.
func (h *Handler) NLSearch(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Query string `json:"query"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	result, err := h.svc.NLSearch(r.Context(), req.Query)
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.WriteJSON(w, http.StatusOK, result)
}

// MoodDiscovery handles POST /api/ai/mood.
func (h *Handler) MoodDiscovery(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromCtx(r.Context())
	var req struct {
		Mood string `json:"mood"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	items, err := h.mediaSvc.GetAllForUser(r.Context(), claims.UserID)
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	result, err := h.svc.MoodDiscovery(r.Context(), req.Mood, items)
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.WriteJSON(w, http.StatusOK, result)
}

// DetectDuplicates handles POST /api/ai/duplicates.
func (h *Handler) DetectDuplicates(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromCtx(r.Context())
	var req struct {
		Title     string `json:"title"`
		MediaType string `json:"media_type"`
		Creator   string `json:"creator"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	items, err := h.mediaSvc.GetAllForUser(r.Context(), claims.UserID)
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	result, err := h.svc.DetectDuplicates(r.Context(), req.Title, req.MediaType, req.Creator, items)
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.WriteJSON(w, http.StatusOK, result)
}

func jsonStr(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}
