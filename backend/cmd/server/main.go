// Package main is the entry point for the EMS API server.
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/your-org/ems/internal/activity"
	"github.com/your-org/ems/internal/ai"
	"github.com/your-org/ems/internal/auth"
	"github.com/your-org/ems/internal/config"
	"github.com/your-org/ems/internal/db"
	"github.com/your-org/ems/internal/httputil"
	"github.com/your-org/ems/internal/media"
	"github.com/your-org/ems/internal/metadata"
	"github.com/your-org/ems/internal/profile"
	"github.com/your-org/ems/internal/search"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	cfg, err := config.Load()
	if err != nil {
		slog.Error("load config", "error", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pool, err := db.New(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("connect database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	// Auth
	tokenSvc := auth.NewTokenService(cfg.JWTSecret, cfg.JWTExpiration)
	authSvc := auth.NewService(pool.Pool, tokenSvc, cfg.BcryptCost)
	authHandler := auth.NewHandler(authSvc)

	// Metadata providers
	tmdbClient := metadata.NewTMDBClient(cfg.TMDBAPIKey)
	mbClient := metadata.NewMusicBrainzClient()
	igdbClient := metadata.NewIGDBClient(cfg.IGDBClientID, cfg.IGDBClientSecret)
	metaSvc := metadata.NewService(tmdbClient, mbClient, igdbClient)
	metaHandler := metadata.NewHandler(metaSvc)

	// Media
	mediaRepo := media.NewRepository(pool.Pool)
	mediaSvc := media.NewService(mediaRepo, metaSvc)
	mediaHandler := media.NewHandler(mediaSvc)

	// AI
	aiClient := ai.NewClient(cfg.AnthropicAPIKey)
	aiCache := ai.NewLRUCache(100)
	aiSvc, err := ai.NewService(aiClient, aiCache)
	if err != nil {
		slog.Error("init ai service", "error", err)
		os.Exit(1)
	}
	aiHandler := ai.NewHandler(aiSvc, mediaSvc)

	// Search
	searchHandler := search.NewHandler(mediaRepo)

	// Profile
	profileHandler := profile.NewHandler(pool.Pool, mediaRepo)

	// Activity
	activityRepo := activity.NewRepository(pool.Pool)

	// Router
	r := httputil.NewRouter(httputil.RouterConfig{FrontendURL: cfg.FrontendURL})

	r.Get("/health", httputil.HealthHandler())

	r.Route("/api", func(r chi.Router) {
		r.Post("/auth/register", authHandler.Register)
		r.Post("/auth/login", authHandler.Login)
		r.Get("/profile/{username}", profileHandler.GetPublic)

		r.Group(func(r chi.Router) {
			r.Use(authSvc.RequireAuth)

			r.Get("/auth/me", authHandler.Me)

			r.Get("/media", mediaHandler.List)
			r.Post("/media", mediaHandler.Create)
			r.Get("/media/{id}", mediaHandler.Get)
			r.Put("/media/{id}", mediaHandler.Update)
			r.Delete("/media/{id}", mediaHandler.Delete)
			r.Patch("/media/{id}/status", mediaHandler.UpdateStatus)

			r.Get("/search", searchHandler.Search)
			r.Post("/metadata/search", metaHandler.Search)

			r.Get("/ai/recommendations", aiHandler.Recommendations)
			r.Get("/ai/insights", aiHandler.Insights)
			r.Post("/ai/nl-search", aiHandler.NLSearch)
			r.Post("/ai/mood", aiHandler.MoodDiscovery)
			r.Post("/ai/duplicates", aiHandler.DetectDuplicates)

			r.Get("/profile/me", profileHandler.GetMe)
			r.Put("/profile", profileHandler.Update)

			r.Get("/activity", func(w http.ResponseWriter, r *http.Request) {
				claims := auth.ClaimsFromCtx(r.Context())
				events, err := activityRepo.List(r.Context(), claims.UserID, 50)
				if err != nil {
					httputil.WriteError(w, http.StatusInternalServerError, err.Error())
					return
				}
				httputil.WriteJSON(w, http.StatusOK, events)
			})
		})
	})

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("server starting", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	<-quit
	slog.Info("shutting down server")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("server shutdown", "error", err)
	}
}
