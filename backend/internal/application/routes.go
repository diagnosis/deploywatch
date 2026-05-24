package application

import (
	"time"

	"github.com/diagnosis/go-toolkit/middleware"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(app *Application) *chi.Mux {
	allowedOrigins := []string{
		"http://localhost:5173",
	}

	r := chi.NewRouter()
	r.Use(middleware.CorrelationID())
	r.Use(middleware.RequestLogger())
	r.Use(middleware.CORS(allowedOrigins))
	r.Use(middleware.RateLimit(10, 20, 5*time.Minute))

	// health
	r.Get("/api/health", app.healthHandler.HandleHealth)

	//auth
	r.Get("/api/auth/github/login", app.authHandler.HandleLogin)
	r.Get("/api/auth/github/callback", app.authHandler.HandleCallBack)
	r.Post("/api/auth/refresh", app.authHandler.HandleRefresh)
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth(app.authHandler.AuthFunc))
		r.Get("/api/auth/me", app.authHandler.HandleMe)
		r.Get("/api/sse", app.sseHandler.HandleSSE)
		r.Post("/api/repos/watch", app.watchedRepoHandler.HandleAddRepo)
		r.Delete("/api/repos/watch/{id}", app.watchedRepoHandler.HandleRemoveRepo)
		r.Get("/api/repos", app.watchedRepoHandler.HandleListRepo)
		r.Post("/api/auth/logout", app.authHandler.HandleLogout)
		r.Get("/api/events", app.eventHandler.HandleListEvent)
		r.Get("/api/github/repos", app.gitHubClientHandler.HandleListRepos)
	})

	r.Post("/webhooks/github", app.webhookHandler.HandleWebhook)

	return r
}
