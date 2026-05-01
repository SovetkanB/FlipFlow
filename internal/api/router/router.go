package router

import (
	"fmt"
	"net/http"
	"time"

	fliflowMiddleware "github.com/SovetkanB/FlipFlow/internal/api/middleware"
	"github.com/SovetkanB/FlipFlow/internal/domain/auth"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Router struct {
	auth           *auth.Handler
	authMiddleware *fliflowMiddleware.AuthMiddleware
}

func NewRouter(auth *auth.Handler, authMiddleware *fliflowMiddleware.AuthMiddleware) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Timeout(30 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"ok","service":"flipper-backend"}`)
	})

	r.Route("/api/v1", func(r chi.Router) {
		//Public endpoints
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", auth.Register)
			r.Post("/login", auth.Login)
			r.Post("/refresh", auth.Refresh)
		})

		//Private endpoints
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.Authenticate)

			r.Get("/auth/me", auth.Me)
		})
	})

	return r
}
