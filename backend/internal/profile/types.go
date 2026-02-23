// Package profile provides user profile management.
package profile

import (
	"time"

	"github.com/google/uuid"
)

// Profile is the public-facing user profile.
type Profile struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	DisplayName string    `json:"display_name"`
	Bio         string    `json:"bio"`
	AvatarURL   string    `json:"avatar_url"`
	IsPublic    bool      `json:"is_public"`
	CreatedAt   time.Time `json:"created_at"`
}

// UpdateRequest holds profile update fields.
type UpdateRequest struct {
	DisplayName *string `json:"display_name,omitempty"`
	Bio         *string `json:"bio,omitempty"`
	AvatarURL   *string `json:"avatar_url,omitempty"`
	IsPublic    *bool   `json:"is_public,omitempty"`
}
