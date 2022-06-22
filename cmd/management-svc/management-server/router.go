package management_server

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/library/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetupRouter(srv *Server, prom *prometheus.Registry) *chi.Mux {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Route("/admin", func(r chi.Router) {
		r.Use(middleware.ChainMiddlewares(true, promMetrics, srv.Env)...)
		r.Get("/complete-history", srv.getCompleteHistory)
		r.Get("/get-history/{id}", srv.getHistory)
		r.Get("/borrowed-history", srv.getBorrowHistory)
		r.Get("/returned-history", srv.getReturnHistory)
		r.Get("/overdue-history", srv.getOverdueHistory)
		r.Get("/student-return-books", srv.getAllBooksStudentReturned)
		r.Get("/student-return-book/{id}", srv.getBooksStudentReturned)
		r.Get("/confirm-return-book/{id}", srv.adminConfirmReturnBook)
		r.Delete("/delete-book/{id}", srv.deleteBook)
		r.Put("/update-book/{id}", srv.updateBook)
		r.Get("/update-book-overdue", srv.updateBookOverdue)
	})
	r.Route("/user", func(r chi.Router) {
		r.Use(middleware.ChainMiddlewares(true, promMetrics, srv.Env)...)
		r.Get("/get-book-overdue-by-student/{id}", srv.getBooksStudentOverdue)
		r.Get("/get-book-reserved-by-student/{id}", srv.getBooksStudentReserved)
		r.Post("/reserve-book/{id}", srv.reserveBook)
		r.Post("/return-book/{id}", srv.studentReturnBook)
		r.Get("/check-availability/{id}", srv.checkAvailability)
	})
	r.Get("/health", srv.health())
	r.Handle("/metrics", promhttp.HandlerFor(prom, promhttp.HandlerOpts{}))

	return r
}
