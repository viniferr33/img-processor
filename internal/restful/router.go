package restful

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/viniferr33/img-processor/pkg/logger"
)

func NewRouter(authHandler *authHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(logger.Middleware())
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("healthy"))
	})

	r.Route("/api", func(r chi.Router) {
		// Public routes
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authHandler.handleUserRegistration)
			r.Post("/login", authHandler.handleUserLogin)
		})

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(authHandler.AuthMiddleware)

			r.Route("/image", func(r chi.Router) {
			})
		})
	})

	return r
}
