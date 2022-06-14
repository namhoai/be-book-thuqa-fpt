package book_server

import (
	"github.com/go-chi/chi"
	"github.com/library/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetupRouter(srv *Server) *chi.Mux {
	r := chi.NewRouter()
	r.Route("/admin/add", func(r chi.Router) {
		r.Use(middleware.ChainMiddlewares(true, promMetrics, srv.Env)...)
		r.Post("/book", srv.addBook)
	})
	r.Route("/get", func(r chi.Router) {
		r.Use(middleware.ChainMiddlewares(true, promMetrics, srv.Env)...)
		r.Get("/books", srv.getBooks)
		r.Get("/books-by-title/{title}", srv.getBooksByTitle)
		r.Get("/books-by-isbn/{isbn}", srv.getBooksByISBN)
		r.Get("/book-by-id/{id}", srv.getBookByBookID)
		r.Get("/book-by-stock/{id}", srv.getBooksByStock)
		r.Get("/book-by-author/{id}", srv.getBooksByAuthor)
		r.Get("/book-by-year/{id}", srv.getBooksByYear)
		r.Get("/book-by-edition/{id}", srv.getBooksByEdition)
		r.Get("/book-by-available/{id}", srv.getBooksByAvailable)
		r.Get("/book-by-borrow/{id}", srv.getBooksByBorrow)
	})
	r.Get("/health", srv.health())
	r.Handle("/metrics", promhttp.HandlerFor(prom, promhttp.HandlerOpts{}))

	return r
}
