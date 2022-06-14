package management_server

import (
	"github.com/go-chi/chi"
	"github.com/library/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetupRouter(srv *Server, prom *prometheus.Registry) *chi.Mux {
	r := chi.NewRouter()
	r.Route("/admin", func(r chi.Router) {
		r.Use(middleware.ChainMiddlewares(true, promMetrics, srv.Env)...)
		r.Get("/complete-history", srv.getCompleteHistory)
		r.Get("/get-history/{id}", srv.getHistory)
		r.Get("/borrowed-history", srv.getBorrowHistory)
		r.Get("/returned-history", srv.getReturnHistory)
		r.Get("/overdue-history", srv.getOverdueHistory)
		r.Delete("/delete-book/{id}", srv.deleteBook)
		r.Patch("/update-book/{id}", srv.updateBook)
	})
	r.Route("/user", func(r chi.Router) {
		r.Use(middleware.ChainMiddlewares(true, promMetrics, srv.Env)...)
		r.Post("/reserve-book/{id}", srv.reserveBook)
		r.Get("/check-availability/{id}", srv.checkAvailability)
	})
	r.Get("/health", srv.health())
	r.Handle("/metrics", promhttp.HandlerFor(prom, promhttp.HandlerOpts{}))

	return r
}
