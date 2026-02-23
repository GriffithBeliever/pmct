package profile

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/ems/internal/auth"
	"github.com/your-org/ems/internal/httputil"
	"github.com/your-org/ems/internal/media"
)

// Handler handles HTTP requests for profile endpoints.
type Handler struct {
	db        *pgxpool.Pool
	mediaRepo *media.Repository
}

// NewHandler creates a new profile Handler.
func NewHandler(db *pgxpool.Pool, mediaRepo *media.Repository) *Handler {
	return &Handler{db: db, mediaRepo: mediaRepo}
}

// GetPublic handles GET /api/profile/:username.
func (h *Handler) GetPublic(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	var profile Profile
	err := h.db.QueryRow(r.Context(), `
		SELECT id, username, display_name, bio, avatar_url, is_public, created_at
		FROM users WHERE username = $1
	`, username).Scan(
		&profile.ID, &profile.Username, &profile.DisplayName,
		&profile.Bio, &profile.AvatarURL, &profile.IsPublic, &profile.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			httputil.WriteError(w, http.StatusNotFound, "profile not found")
			return
		}
		httputil.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !profile.IsPublic {
		httputil.WriteError(w, http.StatusForbidden, "profile is private")
		return
	}

	items, _, err := h.mediaRepo.List(r.Context(), media.ListFilter{
		UserID:   profile.ID,
		PageSize: 100,
	})
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputil.WriteJSON(w, http.StatusOK, map[string]any{
		"profile": profile,
		"items":   items,
	})
}

// Update handles PUT /api/profile.
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromCtx(r.Context())

	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	sets := []string{}
	args := []any{}
	argIdx := 1

	if req.DisplayName != nil {
		sets = append(sets, fmt.Sprintf("display_name=$%d", argIdx))
		args = append(args, *req.DisplayName)
		argIdx++
	}
	if req.Bio != nil {
		sets = append(sets, fmt.Sprintf("bio=$%d", argIdx))
		args = append(args, *req.Bio)
		argIdx++
	}
	if req.AvatarURL != nil {
		sets = append(sets, fmt.Sprintf("avatar_url=$%d", argIdx))
		args = append(args, *req.AvatarURL)
		argIdx++
	}
	if req.IsPublic != nil {
		sets = append(sets, fmt.Sprintf("is_public=$%d", argIdx))
		args = append(args, *req.IsPublic)
		argIdx++
	}

	if len(sets) == 0 {
		httputil.WriteError(w, http.StatusBadRequest, "no fields to update")
		return
	}

	args = append(args, claims.UserID)
	var profile Profile
	err := h.db.QueryRow(r.Context(),
		fmt.Sprintf(`UPDATE users SET %s WHERE id=$%d
			RETURNING id, username, display_name, bio, avatar_url, is_public, created_at`,
			strings.Join(sets, ","), argIdx),
		args...,
	).Scan(
		&profile.ID, &profile.Username, &profile.DisplayName,
		&profile.Bio, &profile.AvatarURL, &profile.IsPublic, &profile.CreatedAt,
	)
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputil.WriteJSON(w, http.StatusOK, profile)
}

// GetMe handles GET /api/profile/me.
func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromCtx(r.Context())

	var profile Profile
	err := h.db.QueryRow(r.Context(), `
		SELECT id, username, display_name, bio, avatar_url, is_public, created_at
		FROM users WHERE id = $1
	`, claims.UserID).Scan(
		&profile.ID, &profile.Username, &profile.DisplayName,
		&profile.Bio, &profile.AvatarURL, &profile.IsPublic, &profile.CreatedAt,
	)
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httputil.WriteJSON(w, http.StatusOK, profile)
}
