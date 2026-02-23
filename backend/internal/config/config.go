// Package config provides application configuration loaded from environment variables.
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all application configuration.
type Config struct {
	// Server
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
	FrontendURL     string

	// Database
	DatabaseURL string

	// Auth
	JWTSecret     string
	JWTExpiration time.Duration
	BcryptCost    int

	// External APIs
	AnthropicAPIKey  string
	TMDBAPIKey       string
	IGDBClientID     string
	IGDBClientSecret string
}

// Option is a functional option for Config.
type Option func(*Config)

// WithPort sets the server port.
func WithPort(port string) Option {
	return func(c *Config) { c.Port = port }
}

// WithFrontendURL sets the allowed frontend origin for CORS.
func WithFrontendURL(url string) Option {
	return func(c *Config) { c.FrontendURL = url }
}

// Load reads configuration from environment variables and applies options.
func Load(opts ...Option) (*Config, error) {
	cfg := &Config{
		Port:            getEnvOrDefault("PORT", "8080"),
		ReadTimeout:     15 * time.Second,
		WriteTimeout:    60 * time.Second,
		ShutdownTimeout: 10 * time.Second,
		JWTExpiration:   7 * 24 * time.Hour,
		BcryptCost:      12,
		FrontendURL:     getEnvOrDefault("FRONTEND_URL", "http://localhost:3000"),
	}

	cfg.DatabaseURL = os.Getenv("DATABASE_URL")
	cfg.JWTSecret = os.Getenv("JWT_SECRET")
	cfg.AnthropicAPIKey = os.Getenv("ANTHROPIC_API_KEY")
	cfg.TMDBAPIKey = os.Getenv("TMDB_API_KEY")
	cfg.IGDBClientID = os.Getenv("IGDB_CLIENT_ID")
	cfg.IGDBClientSecret = os.Getenv("IGDB_CLIENT_SECRET")

	if costStr := os.Getenv("BCRYPT_COST"); costStr != "" {
		cost, err := strconv.Atoi(costStr)
		if err != nil {
			return nil, fmt.Errorf("parse BCRYPT_COST: %w", err)
		}
		cfg.BcryptCost = cost
	}

	for _, opt := range opts {
		opt(cfg)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config validation: %w", err)
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	if c.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	return nil
}

func getEnvOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
