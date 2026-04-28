package router

import (
	"net/http"

	"github.com/SovetkanB/FlipFlow/internal/domain/auth"
	"github.com/go-chi/chi"
)

type Router struct {
	auth *auth.Handler
}

func NewRouter(auth *auth.Handler) http.Handler {
	r := chi.NewRouter()

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", auth.Register)
		})
	})

	return r
}
