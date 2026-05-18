package router

import (
	"fmt"
	"net/http"
	"time"

	flipMiddle "github.com/SovetkanB/FlipFlow/internal/api/middleware"
	"github.com/SovetkanB/FlipFlow/internal/domain/auth"
	"github.com/SovetkanB/FlipFlow/internal/domain/expense"
	"github.com/SovetkanB/FlipFlow/internal/domain/project"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Router struct {
	authS   *auth.Service
	authH   *auth.Handler
	project *project.Handler
	expense *expense.Handler
}

func NewRouter(authS *auth.Service, auth *auth.Handler, project *project.Handler, expense *expense.Handler) http.Handler {
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
			r.Use(flipMiddle.JWT(authS))

			r.Get("/auth/me", auth.Me)

			r.Route("/projects", func(r chi.Router) {
				r.Get("/", project.List)
				r.Post("/", project.Create)

				r.Route("/{projectID}", func(r chi.Router) {
					r.Get("/", project.GetByID)
					r.Patch("/", project.Update)
					r.Delete("/", project.Delete)
					r.Get("/summary", project.FinancialSummary)

					r.Route("/expenses", func(r chi.Router) {
						r.Get("/", expense.List)
						r.Post("/", expense.Create)

						r.Route("/{expenseID}", func(r chi.Router) {
							r.Delete("/", expense.Delete)
						})
					})
				})
			})
		})
	})

	return r
}
