package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/VinVorteX/NoBurn/internal/handlers/auth"
	"github.com/VinVorteX/NoBurn/internal/handlers/employee"
	"github.com/VinVorteX/NoBurn/internal/handlers/settings"
	"github.com/VinVorteX/NoBurn/internal/handlers/survey"
	"github.com/VinVorteX/NoBurn/internal/handlers/analytics"
	"github.com/VinVorteX/NoBurn/internal/handlers/health"
	"github.com/VinVorteX/NoBurn/internal/handlers/webhook"
	middlewareAuth "github.com/VinVorteX/NoBurn/internal/middleware/auth"
)

func New() http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.Throttle(100))

	// CORS
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
			if r.Method == "OPTIONS" {
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	// Public routes
	r.Get("/health", health.HealthCheck)
	r.Post("/auth/login", auth.Login)
	r.Post("/auth/register", auth.Register)
	
	// Public survey routes (no auth required)
	r.Get("/api/surveys/{surveyID}/public", survey.GetPublicSurvey)
	r.Post("/api/surveys/responses/public", survey.SubmitPublicResponse)
	
	// Webhook routes
	r.Post("/webhooks/slack", webhook.HandleSlackWebhook)
	r.Post("/webhooks/survey", webhook.HandleSurveyWebhook)
	r.Post("/webhooks/email", webhook.HandleEmailWebhook)

	// Protected routes
	r.Route("/api", func(r chi.Router) {
		r.Use(middlewareAuth.RequireAuth)
		
		// Employee routes
		r.Post("/employees", employee.AddEmployee)
		r.Post("/employees/bulk", employee.BulkUploadEmployees)
		r.Get("/employees", employee.GetEmployees)
		
		// Survey routes
		r.Post("/surveys", survey.CreateSurvey)
		r.Get("/surveys", survey.GetSurveys)
		r.Post("/surveys/responses", survey.SubmitResponse)
		
		// Analytics routes
		r.Get("/dashboard", analytics.GetDashboard)
		r.Get("/attrition-risks", analytics.GetAttritionRisks)
		r.Get("/retention-suggestions/{userID}", analytics.GetRetentionSuggestions)
		
		// Settings routes
		r.Get("/settings/smtp", settings.GetSMTPSettings)
		r.Put("/settings/smtp", settings.UpdateSMTPSettings)
	})

	return r
}