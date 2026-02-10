package httpserver

import (
	"net/http"

	"rest_api/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(h *handlers.Handler) http.Handler {
	r := chi.NewRouter()

	// middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// public auth routes
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", h.Register)
		r.Post("/login", h.Login)
		r.Post("/refresh", h.Refresh)
	})

	// private tasks routes
	r.Route("/tasks", func(r chi.Router) {
		r.Use(h.RequireAuth)

		r.Get("/", h.GetAllTasks)
		r.Post("/", h.CreateTask)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.GetTaskByID)
			r.Patch("/", h.UpdateTask)
			r.Delete("/", h.DeleteTask)
		})
	})

	return r
}
