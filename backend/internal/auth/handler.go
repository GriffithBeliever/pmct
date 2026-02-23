package auth

import (
	"encoding/json"
	"net/http"

	"github.com/your-org/ems/internal/httputil"
)

// Handler handles HTTP requests for auth endpoints.
type Handler struct {
	svc *Service
}

// NewHandler creates a new auth Handler.
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// Register handles POST /api/auth/register.
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, token, err := h.svc.Register(r.Context(), req)
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	httputil.WriteJSON(w, http.StatusCreated, map[string]any{
		"user":  user,
		"token": token,
	})
}

// Login handles POST /api/auth/login.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, token, err := h.svc.Login(r.Context(), req)
	if err != nil {
		httputil.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	httputil.WriteJSON(w, http.StatusOK, map[string]any{
		"user":  user,
		"token": token,
	})
}

// Me handles GET /api/auth/me.
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	claims := ClaimsFromCtx(r.Context())
	if claims == nil {
		httputil.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, err := h.svc.GetByID(r.Context(), claims.UserID)
	if err != nil {
		httputil.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	httputil.WriteJSON(w, http.StatusOK, user)
}
