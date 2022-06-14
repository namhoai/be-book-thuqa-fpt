package book_server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm"
	"github.com/library/middleware"
	"github.com/library/models"
	"github.com/sirupsen/logrus"
)

func GetAuthInfoFromContext(ctx context.Context) *models.AuthInfo {
	return ctx.Value(middleware.ContextAuthInfo).(*models.AuthInfo)
}

func (srv *Server) addBook(wr http.ResponseWriter, r *http.Request) {
	w := &middleware.LogResponseWriter{ResponseWriter: wr}
	ctx := r.Context()
	authInfo := GetAuthInfoFromContext(ctx)
	if authInfo.Role != models.AdminAccount {
		handleError(w, ctx, srv, "add_book", errors.New("permission denied"), http.StatusUnauthorized)
		return
	}
	book := &models.Book{}
	err := json.NewDecoder(r.Body).Decode(book)
	if err != nil {
		handleError(w, ctx, srv, "add_book", err, http.StatusInternalServerError)
		return
	}
	book.AvailableDate = time.Now()
	book.Available = true
	err = srv.DB.CreateBook(*book)
	if err != nil {
		if strings.Contains(err.Error(), "1062") {
			handleError(w, ctx, srv, "add_book", err, http.StatusBadRequest)
			return
		}
		handleError(w, ctx, srv, "add_book", err, http.StatusInternalServerError)
	}
	logrus.WithFields(logrus.Fields{
		"statusCode": http.StatusOK,
	}).Info(fmt.Sprintf("new book added: %v", book.Name))
	err = json.NewEncoder(w).Encode(book)
	if err != nil {
		handleError(w, ctx, srv, "adding book", err, http.StatusInternalServerError)
		return
	}
}

func (srv *Server) getBooks(wr http.ResponseWriter, r *http.Request) {
	w := &middleware.LogResponseWriter{ResponseWriter: wr}
	ctx := r.Context()
	books, err := srv.DB.GetBooks()
	if err != nil {
		if err == gorm.ErrRecordNotFound || books == nil {
			handleError(w, ctx, srv, "get_books", errors.New("no record found"), http.StatusOK)
			return
		}
		handleError(w, ctx, srv, "get_books", err, http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(books)
	if err != nil {
		handleError(w, ctx, srv, "get_books", err, http.StatusInternalServerError)
	}
}

func (srv *Server) getBooksByTitle(wr http.ResponseWriter, r *http.Request) {
	w := &middleware.LogResponseWriter{ResponseWriter: wr}
	ctx := r.Context()
	title := chi.URLParam(r, "title")
	books, err := srv.DB.GetBooksByTitle(title)
	if err != nil {
		if err == gorm.ErrRecordNotFound || len(*books) == 0 {
			handleError(w, ctx, srv, "get_books_by_title", errors.New("no record found"), http.StatusOK)
			return
		}
		handleError(w, ctx, srv, "get_books_by_title", err, http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(books)
	if err != nil {
		handleError(w, ctx, srv, "get_books_by_title", err, http.StatusInternalServerError)
	}
}

func (srv *Server) getBooksByISBN(wr http.ResponseWriter, r *http.Request) {
	w := &middleware.LogResponseWriter{ResponseWriter: wr}
	ctx := r.Context()
	isbn := chi.URLParam(r, "isbn")
	books, err := srv.DB.GetBooksByISBN(isbn)
	if err != nil {
		if err == gorm.ErrRecordNotFound || len(*books) == 0 {
			handleError(w, ctx, srv, "get_books_by_isbn", errors.New("no record found"), http.StatusOK)
			return
		}
		handleError(w, ctx, srv, "get_books_by_isbn", err, http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(books)
	if err != nil {
		handleError(w, ctx, srv, "get_books_by_isbn", err, http.StatusInternalServerError)
	}
}

func (srv *Server) getBookByBookID(wr http.ResponseWriter, r *http.Request) {
	w := &middleware.LogResponseWriter{ResponseWriter: wr}
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	bookID, err := strconv.Atoi(id)
	if err != nil {
		handleError(w, ctx, srv, "get_book_by_id", err, http.StatusInternalServerError)
		return
	}
	book, err := srv.DB.GetBookByID(uint(bookID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleError(w, ctx, srv, "get_book_by_id", errors.New("no record found"), http.StatusOK)
			return
		}
		handleError(w, ctx, srv, "get_book_by_id", err, http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(book)
	if err != nil {
		handleError(w, ctx, srv, "get_book_by_id", err, http.StatusInternalServerError)
	}
}

func (srv *Server) getBooksByStock(wr http.ResponseWriter, r *http.Request) {
	w := &middleware.LogResponseWriter{ResponseWriter: wr}
	ctx := r.Context()
	stock := chi.URLParam(r, "stock")
	stockInt, _ := strconv.Atoi(stock)
	books, err := srv.DB.GetBooksByStock(uint(stockInt))
	if err != nil {
		if err == gorm.ErrRecordNotFound || len(*books) == 0 {
			handleError(w, ctx, srv, "get_books_by_stock", errors.New("no record found"), http.StatusOK)
			return
		}
		handleError(w, ctx, srv, "get_books_by_stock", err, http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(books)
	if err != nil {
		handleError(w, ctx, srv, "get_books_by_stock", err, http.StatusInternalServerError)
	}
}

func (srv *Server) getBooksByAuthor(wr http.ResponseWriter, r *http.Request) {
	w := &middleware.LogResponseWriter{ResponseWriter: wr}
	ctx := r.Context()
	author := chi.URLParam(r, "author")
	books, err := srv.DB.GetBooksByAuthor(author)
	if err != nil {
		if err == gorm.ErrRecordNotFound || len(*books) == 0 {
			handleError(w, ctx, srv, "get_books_by_author", errors.New("no record found"), http.StatusOK)
			return
		}
		handleError(w, ctx, srv, "get_books_by_author", err, http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(books)
	if err != nil {
		handleError(w, ctx, srv, "get_books_by_author", err, http.StatusInternalServerError)
	}
}

func (srv *Server) getBooksByYear(wr http.ResponseWriter, r *http.Request) {
	w := &middleware.LogResponseWriter{ResponseWriter: wr}
	ctx := r.Context()
	year := chi.URLParam(r, "year")
	books, err := srv.DB.GetBooksByYear(year)
	if err != nil {
		if err == gorm.ErrRecordNotFound || len(*books) == 0 {
			handleError(w, ctx, srv, "get_books_by_year", errors.New("no record found"), http.StatusOK)
			return
		}
		handleError(w, ctx, srv, "get_books_by_year", err, http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(books)
	if err != nil {
		handleError(w, ctx, srv, "get_books_by_year", err, http.StatusInternalServerError)
	}
}

func (srv *Server) getBooksByEdition(wr http.ResponseWriter, r *http.Request) {
	w := &middleware.LogResponseWriter{ResponseWriter: wr}
	ctx := r.Context()
	edition := chi.URLParam(r, "edition")
	editionInt, _ := strconv.Atoi(edition)
	books, err := srv.DB.GetBooksByEdition(uint(editionInt))
	if err != nil {
		if err == gorm.ErrRecordNotFound || len(*books) == 0 {
			handleError(w, ctx, srv, "get_books_by_edition", errors.New("no record found"), http.StatusOK)
			return
		}
		handleError(w, ctx, srv, "get_books_by_edition", err, http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(books)
	if err != nil {
		handleError(w, ctx, srv, "get_books_by_edition", err, http.StatusInternalServerError)
	}
}

func (srv *Server) getBooksByAvailable(wr http.ResponseWriter, r *http.Request) {
	w := &middleware.LogResponseWriter{ResponseWriter: wr}
	ctx := r.Context()
	available := chi.URLParam(r, "available")
	if available == "false" {
		return
	}
	books, err := srv.DB.GetBooksByAvailable(true)
	if err != nil {
		if err == gorm.ErrRecordNotFound || len(*books) == 0 {
			handleError(w, ctx, srv, "get_books_by_available", errors.New("no record found"), http.StatusOK)
			return
		}
		handleError(w, ctx, srv, "get_books_by_available", err, http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(books)
	if err != nil {
		handleError(w, ctx, srv, "get_books_by_available", err, http.StatusInternalServerError)
	}
}

func (srv *Server) getBooksByBorrow(wr http.ResponseWriter, r *http.Request) {
	w := &middleware.LogResponseWriter{ResponseWriter: wr}
	ctx := r.Context()
	available := chi.URLParam(r, "available")
	if available == "true" {
		return
	}
	books, err := srv.DB.GetBooksByAvailable(false)
	if err != nil {
		if err == gorm.ErrRecordNotFound || len(*books) == 0 {
			handleError(w, ctx, srv, "get_books_by_borrow", errors.New("no record found"), http.StatusOK)
			return
		}
		handleError(w, ctx, srv, "get_books_by_borrow", err, http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(books)
	if err != nil {
		handleError(w, ctx, srv, "get_books_by_borrow", err, http.StatusInternalServerError)
	}
}

func (srv *Server) health() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

func handleError(w *middleware.LogResponseWriter, ctx context.Context, srv *Server, task string, err error, statusCode int) {
	if !srv.TestRun {
		srv.TracingID = ctx.Value(middleware.RequestTracingID).(string)
	}
	http.Error(w, err.Error(), statusCode)

	logrus.WithFields(logrus.Fields{
		"tracingID":  srv.TracingID,
		"statusCode": statusCode,
		"error":      err,
	}).Error(task)
}
