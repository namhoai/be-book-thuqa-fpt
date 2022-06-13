package user_server

import (
	"github.com/go-chi/chi"
	"github.com/library/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetupRouter(srv *Server, prom *prometheus.Registry) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.ChainMiddlewares(false, promMetrics, srv.Env)...)
	r.Post("/register", srv.register())
	r.Post("/login", srv.login())
	r.Get("/health", srv.health())
	r.Handle("/metrics", promhttp.HandlerFor(prom, promhttp.HandlerOpts{}))

	r.Route("/get", func(r chi.Router) {
		r.Use(middleware.ChainMiddlewares(true, promMetrics, srv.Env)...)
		r.Get("/users", srv.getUsers)
		r.Get("/users-by-name/{name}", srv.getUserByName)
		r.Get("/users-by-email/{email}", srv.getUserByEmail)
		r.Get("/users-by-id/{id}", srv.getUserByID)
	})

	return r
}
