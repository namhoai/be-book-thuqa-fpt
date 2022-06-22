package book_server

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/library/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetupRouter(srv *Server) *chi.Mux {
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

	r.Route("/admin/add", func(r chi.Router) {
		r.Use(middleware.ChainMiddlewares(true, promMetrics, srv.Env)...)
		r.Post("/book", srv.addBook)
	})
	r.Route("/admin", func(r chi.Router) {
		r.Use(middleware.ChainMiddlewares(true, promMetrics, srv.Env)...)
		r.Post("/updoad-image", srv.uploadImageToS3)
		r.Post("/download-image", srv.downloadImageFromS3)
	})
	r.Route("/get", func(r chi.Router) {
		r.Use(middleware.ChainMiddlewares(false, promMetrics, srv.Env)...)
		r.Get("/books", srv.getBooks)
		r.Get("/books-by-title/{title}", srv.getBooksByTitle)
		r.Get("/books-by-isbn/{isbn}", srv.getBooksByISBN)
		r.Get("/book-by-id/{id}", srv.getBookByBookID)
		r.Get("/book-by-stock/{stock}", srv.getBooksByStock)
		r.Get("/book-by-author/{author}", srv.getBooksByAuthor)
		r.Get("/book-by-year/{year}", srv.getBooksByYear)
		r.Get("/book-by-edition/{edition}", srv.getBooksByEdition)
		r.Get("/book-available", srv.getBooksByAvailable)
		r.Get("/book-borrow", srv.getBooksByBorrow)
	})
	r.Get("/health", srv.health())
	r.Handle("/metrics", promhttp.HandlerFor(prom, promhttp.HandlerOpts{}))

	return r
}
