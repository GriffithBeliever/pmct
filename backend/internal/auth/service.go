package auth

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// User represents the authenticated user record.
type User struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	DisplayName string    `json:"display_name"`
	Bio         string    `json:"bio"`
	AvatarURL   string    `json:"avatar_url"`
	IsPublic    bool      `json:"is_public"`
}

// RegisterRequest holds registration parameters.
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest holds login parameters.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Service handles authentication business logic.
type Service struct {
	db         *pgxpool.Pool
	tokenSvc   *TokenService
	bcryptCost int
}

// NewService creates a new auth Service.
func NewService(db *pgxpool.Pool, tokenSvc *TokenService, bcryptCost int) *Service {
	return &Service{
		db:         db,
		tokenSvc:   tokenSvc,
		bcryptCost: bcryptCost,
	}
}

// Register creates a new user account.
func (s *Service) Register(ctx context.Context, req RegisterRequest) (*User, string, error) {
	if len(req.Username) < 3 {
		return nil, "", fmt.Errorf("username must be at least 3 characters")
	}
	if len(req.Password) < 8 {
		return nil, "", fmt.Errorf("password must be at least 8 characters")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), s.bcryptCost)
	if err != nil {
		return nil, "", fmt.Errorf("hash password: %w", err)
	}

	var user User
	err = s.db.QueryRow(ctx, `
		INSERT INTO users (username, email, password_hash, display_name)
		VALUES ($1, $2, $3, $1)
		RETURNING id, username, email, display_name, bio, avatar_url, is_public
	`, req.Username, req.Email, string(hash)).Scan(
		&user.ID, &user.Username, &user.Email,
		&user.DisplayName, &user.Bio, &user.AvatarURL, &user.IsPublic,
	)
	if err != nil {
		return nil, "", fmt.Errorf("insert user: %w", err)
	}

	token, err := s.tokenSvc.Sign(user.ID, user.Username)
	if err != nil {
		return nil, "", fmt.Errorf("sign token: %w", err)
	}

	return &user, token, nil
}

// Login authenticates a user and returns a JWT token.
func (s *Service) Login(ctx context.Context, req LoginRequest) (*User, string, error) {
	var user User
	var passwordHash string

	err := s.db.QueryRow(ctx, `
		SELECT id, username, email, password_hash, display_name, bio, avatar_url, is_public
		FROM users WHERE email = $1
	`, req.Email).Scan(
		&user.ID, &user.Username, &user.Email, &passwordHash,
		&user.DisplayName, &user.Bio, &user.AvatarURL, &user.IsPublic,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, "", fmt.Errorf("invalid credentials")
		}
		return nil, "", fmt.Errorf("query user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		return nil, "", fmt.Errorf("invalid credentials")
	}

	token, err := s.tokenSvc.Sign(user.ID, user.Username)
	if err != nil {
		return nil, "", fmt.Errorf("sign token: %w", err)
	}

	return &user, token, nil
}

// GetByID fetches a user by ID.
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	var user User
	err := s.db.QueryRow(ctx, `
		SELECT id, username, email, display_name, bio, avatar_url, is_public
		FROM users WHERE id = $1
	`, id).Scan(
		&user.ID, &user.Username, &user.Email,
		&user.DisplayName, &user.Bio, &user.AvatarURL, &user.IsPublic,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("query user: %w", err)
	}
	return &user, nil
}
