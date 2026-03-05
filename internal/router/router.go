package router

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/zostay/bobsbinder/internal/config"
	"github.com/zostay/bobsbinder/internal/handlers"
	"github.com/zostay/bobsbinder/internal/middleware"
)

func New(db *sql.DB, cfg *config.Config, logger *zap.Logger) *chi.Mux {
	r := chi.NewRouter()

	// Global middleware
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(middleware.RequestLogging(logger))
	r.Use(chimw.Recoverer)

	authHandler := &handlers.AuthHandler{
		DB:     db,
		Config: cfg,
		Logger: logger,
	}

	docHandler := &handlers.DocumentHandler{
		DB:     db,
		Logger: logger,
	}

	// Public routes
	r.Get("/api/health", handlers.HealthCheck)
	r.Post("/api/auth/register", authHandler.Register)
	r.Post("/api/auth/login", authHandler.Login)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuth(cfg.JWTSecret, logger))

		r.Post("/api/auth/refresh", authHandler.Refresh)

		r.Get("/api/documents", docHandler.List)
		r.Post("/api/documents", docHandler.Create)
		r.Get("/api/documents/{id}", docHandler.Get)
		r.Put("/api/documents/{id}", docHandler.Update)
		r.Delete("/api/documents/{id}", docHandler.Delete)
	})

	return r
}
