package management_server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
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

func (srv *Server) getCompleteHistory(wr http.ResponseWriter, r *http.Request) {
	w := middleware.NewLogResponseWriter(wr)
	ctx := r.Context()
	authInfo := GetAuthInfoFromContext(ctx)
	if authInfo.Role != models.AdminAccount {
		handleError(w, ctx, srv, "get_complete_history", errors.New("permission denied"), http.StatusUnauthorized)
		return
	}
	history, err := srv.DB.GetCompleteHistory()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleError(w, ctx, srv, "get_complete_history", errors.New("no record found"), http.StatusOK)
		} else {
			handleError(w, ctx, srv, "get_complete_history", err, http.StatusInternalServerError)
		}
		return
	}
	err = json.NewEncoder(w).Encode(history)
	if err != nil {
		handleError(w, ctx, srv, "get_complete_history", err, http.StatusInternalServerError)
	}
}

func (srv *Server) getHistory(wr http.ResponseWriter, r *http.Request) {
	w := middleware.NewLogResponseWriter(wr)
	ctx := r.Context()
	authInfo := GetAuthInfoFromContext(ctx)
	if authInfo.Role != models.AdminAccount {
		handleError(w, ctx, srv, "get_history", errors.New("permission denied"), http.StatusUnauthorized)
		return
	}
	id := chi.URLParam(r, "id")
	bookID, err := strconv.Atoi(id)
	if err != nil {
		handleError(w, ctx, srv, "get_history", err, http.StatusInternalServerError)
		return
	}
	history, err := srv.DB.GetHistory(uint(bookID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleError(w, ctx, srv, "get_history", errors.New("no record found"), http.StatusOK)
		} else {
			handleError(w, ctx, srv, "get_history", err, http.StatusInternalServerError)
		}
		return
	}
	err = json.NewEncoder(w).Encode(history)
	if err != nil {
		handleError(w, ctx, srv, "get_history", err, http.StatusInternalServerError)
	}
}

func (srv *Server) getBorrowHistory(wr http.ResponseWriter, r *http.Request) {
	w := middleware.NewLogResponseWriter(wr)
	ctx := r.Context()
	authInfo := GetAuthInfoFromContext(ctx)
	if authInfo.Role != models.AdminAccount {
		handleError(w, ctx, srv, "get_borrowed_history", errors.New("permission denied"), http.StatusUnauthorized)
		return
	}
	history, err := srv.DB.GetBooksbyStatus("borrowed")
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleError(w, ctx, srv, "get_borrowed_history", errors.New("no record found"), http.StatusOK)
		} else {
			handleError(w, ctx, srv, "get_borrowed_history", err, http.StatusInternalServerError)
		}
		return
	}
	err = json.NewEncoder(w).Encode(history)
	if err != nil {
		handleError(w, ctx, srv, "get_borrowed_history", err, http.StatusInternalServerError)
	}
}

func (srv *Server) getReturnHistory(wr http.ResponseWriter, r *http.Request) {
	w := middleware.NewLogResponseWriter(wr)
	ctx := r.Context()
	authInfo := GetAuthInfoFromContext(ctx)
	if authInfo.Role != models.AdminAccount {
		handleError(w, ctx, srv, "get_returned_history", errors.New("permission denied"), http.StatusUnauthorized)
		return
	}
	history, err := srv.DB.GetBooksbyStatus("returned")
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleError(w, ctx, srv, "get_returned_history", errors.New("no record found"), http.StatusOK)
		} else {
			handleError(w, ctx, srv, "get_returned_history", err, http.StatusInternalServerError)
		}
		return
	}
	err = json.NewEncoder(w).Encode(history)
	if err != nil {
		handleError(w, ctx, srv, "get_returned_history", err, http.StatusInternalServerError)
	}
}

func (srv *Server) getOverdueHistory(wr http.ResponseWriter, r *http.Request) {
	w := middleware.NewLogResponseWriter(wr)
	ctx := r.Context()
	authInfo := GetAuthInfoFromContext(ctx)
	if authInfo.Role != models.AdminAccount {
		handleError(w, ctx, srv, "get_overdue_history", errors.New("permission denied"), http.StatusUnauthorized)
		return
	}
	history, err := srv.DB.GetBooksbyStatus("overdue")
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleError(w, ctx, srv, "get_overdue_history", errors.New("no record found"), http.StatusOK)
		} else {
			handleError(w, ctx, srv, "get_overdue_history", err, http.StatusInternalServerError)
		}
		return
	}
	err = json.NewEncoder(w).Encode(history)
	if err != nil {
		handleError(w, ctx, srv, "get_overdue_history", err, http.StatusInternalServerError)
	}
}

func (srv *Server) checkAvailability(wr http.ResponseWriter, r *http.Request) {
	w := middleware.NewLogResponseWriter(wr)
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	bookID, err := strconv.Atoi(id)
	if err != nil {
		handleError(w, ctx, srv, "check_availability", err, http.StatusInternalServerError)
		return
	}
	avail, err := srv.DB.CheckAvailability(uint(bookID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleError(w, ctx, srv, "check_availability", errors.New("no record found"), http.StatusOK)
			return
		}
		handleError(w, ctx, srv, "check_availability", err, http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(avail)
	if err != nil {
		handleError(w, ctx, srv, "check_availability", err, http.StatusInternalServerError)
	}
}

func (srv *Server) reserveBook(wr http.ResponseWriter, r *http.Request) {
	w := middleware.NewLogResponseWriter(wr)
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	bookID, err := strconv.Atoi(id)
	if err != nil {
		return
	}
	user := r.FormValue("userId")
	reservedDateString := r.FormValue("reservedDate")
	returnDateString := r.FormValue("returnDate")
	reservedDate, err := time.Parse("2006-01-02", reservedDateString)
	if err != nil {
		fmt.Println(err)
		return
	}
	returnDate, err := time.Parse("2006-01-02", returnDateString)
	if err != nil {
		fmt.Println(err)
		return
	}
	userID, err := strconv.Atoi(user)
	if err != nil {
		handleError(w, ctx, srv, "reserve_book", err, http.StatusInternalServerError)
		return
	}
	days := returnDate.Sub(reservedDate).Hours() / 24
	// fmt.Println(days)
	if days > 42.0 {
		// handleError(w, ctx, srv, "Book_cannot_be_reserved_for_more_than_6_weeks", http.StatusBadRequest)
		json.NewEncoder(w).Encode("Book cannot be reserved for more than 6 weeks!")
		return
	}
	err = srv.DB.ReserveBook(uint(bookID), uint(userID), &reservedDate, &returnDate)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleError(w, ctx, srv, "reserve_book", errors.New("no record found"), http.StatusOK)
			return
		}
		handleError(w, ctx, srv, "reserve_book", err, http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode("Book reserved successfully!")
	if err != nil {
		handleError(w, ctx, srv, "reserve_book", err, http.StatusInternalServerError)
	}
}

func (srv *Server) adminConfirmReturnBook(wr http.ResponseWriter, r *http.Request) {
	w := middleware.NewLogResponseWriter(wr)
	ctx := r.Context()
	authInfo := GetAuthInfoFromContext(ctx)
	if authInfo.Role != models.AdminAccount {
		handleError(w, ctx, srv, "return_book", errors.New("permission denied"), http.StatusUnauthorized)
		return
	}
	id := chi.URLParam(r, "id")
	bookID, err := strconv.Atoi(id)
	user := r.FormValue("userId")
	userID, _ := strconv.Atoi(user)
	if err != nil {
		handleError(w, ctx, srv, "return_book", err, http.StatusInternalServerError)
		return
	}
	err = srv.DB.AdminConfirmReturnBook(uint(bookID), uint(userID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleError(w, ctx, srv, "return_book", errors.New("no record found"), http.StatusOK)
			return
		}
		handleError(w, ctx, srv, "return_book", err, http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode("Book return processed successfully!")
	if err != nil {
		handleError(w, ctx, srv, "return_book", err, http.StatusInternalServerError)
	}
	err = srv.DB.DeleteRecordStudentReturnBook(uint(bookID))
	if err != nil {
		handleError(w, ctx, srv, "return_book", err, http.StatusInternalServerError)
		return
	}
}

func (srv *Server) studentReturnBook(wr http.ResponseWriter, r *http.Request) {
	w := middleware.NewLogResponseWriter(wr)
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	user := r.FormValue("userId")
	bookID, err := strconv.Atoi(id)
	if err != nil {
		handleError(w, ctx, srv, "student_return_book", err, http.StatusInternalServerError)
		return
	}
	userID, err := strconv.Atoi(user)
	if err != nil {
		handleError(w, ctx, srv, "student_return_book", err, http.StatusInternalServerError)
		return
	}
	returnDate := time.Now()
	books, err := srv.DB.GetHistory(uint(bookID))
	reservedDate := (*books)[0].ReservedDate
	err = srv.DB.StudentReturnBook(uint(bookID), uint(userID), reservedDate, &returnDate)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleError(w, ctx, srv, "student_return_book", errors.New("no record found"), http.StatusOK)
			return
		}
		handleError(w, ctx, srv, "student_return_book", err, http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode("Wait until the return is accepted!")
	if err != nil {
		handleError(w, ctx, srv, "student_return_book", err, http.StatusInternalServerError)
	}
}

func (srv *Server) updateBookOverdue(wr http.ResponseWriter, r *http.Request) {
	w := middleware.NewLogResponseWriter(wr)
	ctx := r.Context()
	authInfo := GetAuthInfoFromContext(ctx)
	if authInfo.Role != models.AdminAccount {
		handleError(w, ctx, srv, "update_book_overdue", errors.New("permission denied"), http.StatusUnauthorized)
		return
	}
	currentTime := time.Now()
	err := srv.DB.UpdateBookOverdue(&currentTime)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleError(w, ctx, srv, "update_book_overdue", errors.New("no record found"), http.StatusOK)
			return
		}
		handleError(w, ctx, srv, "update_book_overdue", err, http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode("Updated overdue books successfully!")
	if err != nil {
		handleError(w, ctx, srv, "update_book_overdue", err, http.StatusInternalServerError)
	}
}

func (srv *Server) getBooksStudentOverdue(wr http.ResponseWriter, r *http.Request) {
	w := middleware.NewLogResponseWriter(wr)
	ctx := r.Context()
	// authInfo := GetAuthInfoFromContext(ctx)
	// if authInfo.Role != models.AdminAccount {
	// 	handleError(w, ctx, srv, "get_book_overdue_of_student", errors.New("permission denied"), http.StatusUnauthorized)
	// 	return
	// }
	userId := chi.URLParam(r, "id")
	userID, err := strconv.Atoi(userId)
	if err != nil {
		handleError(w, ctx, srv, "get_book_overdue_of_student", err, http.StatusInternalServerError)
		return
	}
	history, err := srv.DB.GetBooksStudentOverdue(uint(userID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleError(w, ctx, srv, "get_book_overdue_of_student", errors.New("no record found"), http.StatusOK)
			return
		}
		handleError(w, ctx, srv, "get_book_overdue_of_student", err, http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(history)
	if err != nil {
		handleError(w, ctx, srv, "get_book_overdue_of_student", err, http.StatusInternalServerError)
	}
}

func (srv *Server) getBooksStudentReserved(wr http.ResponseWriter, r *http.Request) {
	w := middleware.NewLogResponseWriter(wr)
	ctx := r.Context()
	// authInfo := GetAuthInfoFromContext(ctx)
	// if authInfo.Role != models.AdminAccount {
	// 	handleError(w, ctx, srv, "get_book_reserved_of_student", errors.New("permission denied"), http.StatusUnauthorized)
	// 	return
	// }
	userId := chi.URLParam(r, "id")
	userID, err := strconv.Atoi(userId)
	if err != nil {
		handleError(w, ctx, srv, "get_book_reserved_of_student", err, http.StatusInternalServerError)
		return
	}
	history, err := srv.DB.GetBooksStudentReserved(uint(userID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleError(w, ctx, srv, "get_book_reserved_of_student", errors.New("no record found"), http.StatusOK)
			return
		}
		handleError(w, ctx, srv, "get_book_reserved_of_student", err, http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(history)
	if err != nil {
		handleError(w, ctx, srv, "get_book_reserved_of_student", err, http.StatusInternalServerError)
	}
}

func (srv *Server) deleteBook(wr http.ResponseWriter, r *http.Request) {
	w := middleware.NewLogResponseWriter(wr)
	ctx := r.Context()
	authInfo := GetAuthInfoFromContext(ctx)
	if authInfo.Role != models.AdminAccount {
		handleError(w, ctx, srv, "delete_book", errors.New("permission denied"), http.StatusUnauthorized)
		return
	}
	id := chi.URLParam(r, "id")
	bookID, err := strconv.Atoi(id)
	if err != nil {
		handleError(w, ctx, srv, "delete_book", err, http.StatusInternalServerError)
		return
	}
	err = srv.DB.DeleteBook(uint(bookID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleError(w, ctx, srv, "delete_book", errors.New("no record found"), http.StatusOK)
			return
		}
		handleError(w, ctx, srv, "delete_book", err, http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode("Book deleted successfully!")
	if err != nil {
		handleError(w, ctx, srv, "delete_book", err, http.StatusInternalServerError)
	}
}

func (srv *Server) updateBook(wr http.ResponseWriter, r *http.Request) {
	w := middleware.NewLogResponseWriter(wr)
	ctx := r.Context()
	authInfo := GetAuthInfoFromContext(ctx)
	if authInfo.Role != models.AdminAccount {
		handleError(w, ctx, srv, "update_name_of_book", errors.New("permission denied"), http.StatusUnauthorized)
		return
	}
	id := chi.URLParam(r, "id")
	bookID, _ := strconv.Atoi(id)
	bookName := r.FormValue("name")
	isbn := r.FormValue("isbn")
	stock := r.FormValue("stock")
	stockInt, _ := strconv.Atoi(stock)
	author := r.FormValue("author")
	year := r.FormValue("year")
	edition := r.FormValue("edition")
	editionInt, _ := strconv.Atoi(edition)
	cover := r.FormValue("cover")
	abstract := r.FormValue("abstract")
	category := r.FormValue("category")
	rating := r.FormValue("rating")
	ratingInt, _ := strconv.Atoi(rating)

	err := srv.DB.UpdateBook(uint(bookID), bookName, isbn, uint(stockInt), author, year, uint(editionInt), cover, abstract, category, uint(ratingInt))
	if err != nil {
		handleError(w, ctx, srv, "update_name_of_book", err, http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode("Book updated successfully!")
	if err != nil {
		handleError(w, ctx, srv, "update_name_of_book", err, http.StatusInternalServerError)
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

func (srv *Server) getAllBooksStudentReturned(wr http.ResponseWriter, r *http.Request) {
	w := middleware.NewLogResponseWriter(wr)
	ctx := r.Context()
	authInfo := GetAuthInfoFromContext(ctx)
	if authInfo.Role != models.AdminAccount {
		handleError(w, ctx, srv, "get_all_book_student_return", errors.New("permission denied"), http.StatusUnauthorized)
		return
	}

	bookReturnByStudent, err := srv.DB.GetAllBooksStudentReturned()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleError(w, ctx, srv, "get_all_book_student_return", errors.New("no record found"), http.StatusOK)
			return
		}
		handleError(w, ctx, srv, "get_all_book_student_return", err, http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(bookReturnByStudent)
	if err != nil {
		handleError(w, ctx, srv, "get_all_book_student_return", err, http.StatusInternalServerError)
	}
}

func (srv *Server) getBooksStudentReturned(wr http.ResponseWriter, r *http.Request) {
	w := middleware.NewLogResponseWriter(wr)
	ctx := r.Context()
	authInfo := GetAuthInfoFromContext(ctx)
	if authInfo.Role != models.AdminAccount {
		handleError(w, ctx, srv, "get_all_book_student_return", errors.New("permission denied"), http.StatusUnauthorized)
		return
	}
	id := chi.URLParam(r, "id")
	bookID, _ := strconv.Atoi(id)

	bookReturnByStudent, err := srv.DB.GetBooksStudentReturned(uint(bookID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleError(w, ctx, srv, "get_all_book_student_return", errors.New("no record found"), http.StatusOK)
			return
		}
		handleError(w, ctx, srv, "get_all_book_student_return", err, http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(bookReturnByStudent)
	if err != nil {
		handleError(w, ctx, srv, "get_all_book_student_return", err, http.StatusInternalServerError)
	}
}
