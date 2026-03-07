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

	contactHandler := &handlers.ContactHandler{
		DB:     db,
		Logger: logger,
	}

	locationHandler := &handlers.LocationHandler{
		DB:     db,
		Logger: logger,
	}

	digitalAccessHandler := &handlers.DigitalAccessHandler{
		DB:     db,
		Logger: logger,
	}

	insurancePolicyHandler := &handlers.InsurancePolicyHandler{
		DB:     db,
		Logger: logger,
	}

	serviceAccountHandler := &handlers.ServiceAccountHandler{
		DB:     db,
		Logger: logger,
	}

	obituaryInfoHandler := &handlers.ObituaryInfoHandler{
		DB:     db,
		Logger: logger,
	}

	partyHandler := &handlers.PartyHandler{
		DB:     db,
		Logger: logger,
	}

	documentCategoryHandler := &handlers.DocumentCategoryHandler{
		DB:     db,
		Logger: logger,
	}

	letterHandler := &handlers.SurvivorLetterHandler{
		DB:     db,
		Logger: logger,
	}

	confidentialHandler := &handlers.ConfidentialHandler{
		DB:     db,
		Logger: logger,
	}

	checklistHandler := &handlers.ChecklistHandler{
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

		r.Get("/api/parties", partyHandler.List)
		r.Post("/api/parties", partyHandler.Create)
		r.Put("/api/parties/{id}", partyHandler.Update)
		r.Delete("/api/parties/{id}", partyHandler.Delete)
		r.Get("/api/document-categories", documentCategoryHandler.List)

		r.Get("/api/checklist", checklistHandler.ListAll)
		r.Get("/api/parties/{partyId}/checklist", checklistHandler.ListForParty)
		r.Put("/api/parties/{partyId}/checklist/{categoryId}", checklistHandler.UpdateStatus)

		r.Get("/api/documents", docHandler.List)
		r.Post("/api/documents", docHandler.Create)
		r.Get("/api/documents/{id}", docHandler.Get)
		r.Put("/api/documents/{id}", docHandler.Update)
		r.Delete("/api/documents/{id}", docHandler.Delete)

		r.Get("/api/contacts", contactHandler.List)
		r.Post("/api/contacts", contactHandler.Create)
		r.Get("/api/contacts/{id}", contactHandler.Get)
		r.Put("/api/contacts/{id}", contactHandler.Update)
		r.Delete("/api/contacts/{id}", contactHandler.Delete)

		r.Get("/api/locations", locationHandler.List)
		r.Post("/api/locations", locationHandler.Create)
		r.Get("/api/locations/{id}", locationHandler.Get)
		r.Put("/api/locations/{id}", locationHandler.Update)
		r.Delete("/api/locations/{id}", locationHandler.Delete)

		r.Get("/api/digital-access", digitalAccessHandler.List)
		r.Post("/api/digital-access", digitalAccessHandler.Create)
		r.Get("/api/digital-access/{id}", digitalAccessHandler.Get)
		r.Put("/api/digital-access/{id}", digitalAccessHandler.Update)
		r.Delete("/api/digital-access/{id}", digitalAccessHandler.Delete)

		r.Get("/api/insurance-policies", insurancePolicyHandler.List)
		r.Post("/api/insurance-policies", insurancePolicyHandler.Create)
		r.Get("/api/insurance-policies/{id}", insurancePolicyHandler.Get)
		r.Put("/api/insurance-policies/{id}", insurancePolicyHandler.Update)
		r.Delete("/api/insurance-policies/{id}", insurancePolicyHandler.Delete)

		r.Get("/api/service-accounts", serviceAccountHandler.List)
		r.Post("/api/service-accounts", serviceAccountHandler.Create)
		r.Get("/api/service-accounts/{id}", serviceAccountHandler.Get)
		r.Put("/api/service-accounts/{id}", serviceAccountHandler.Update)
		r.Delete("/api/service-accounts/{id}", serviceAccountHandler.Delete)

		r.Get("/api/parties/{partyId}/obituary-info", obituaryInfoHandler.List)
		r.Post("/api/parties/{partyId}/obituary-info", obituaryInfoHandler.Create)
		r.Put("/api/parties/{partyId}/obituary-info/{id}", obituaryInfoHandler.Update)
		r.Delete("/api/parties/{partyId}/obituary-info/{id}", obituaryInfoHandler.Delete)

		r.Get("/api/confidential", confidentialHandler.GetConfidential)

		r.Get("/api/survivor-letter", letterHandler.GetLetter)
		r.Put("/api/survivor-letter", letterHandler.UpdateBoilerplate)
		r.Put("/api/survivor-letter/sections/reorder", letterHandler.ReorderSections)
		r.Put("/api/survivor-letter/sections/{sectionId}", letterHandler.UpdateSection)
		r.Post("/api/survivor-letter/sections/{sectionId}/items", letterHandler.AddItem)
		r.Put("/api/survivor-letter/items/reorder", letterHandler.ReorderItems)
		r.Put("/api/survivor-letter/items/{itemId}", letterHandler.EditItem)
		r.Delete("/api/survivor-letter/items/{itemId}", letterHandler.DeleteItem)
		r.Post("/api/survivor-letter/items/{itemId}/unsuppress", letterHandler.UnsuppressItem)
	})

	return r
}
