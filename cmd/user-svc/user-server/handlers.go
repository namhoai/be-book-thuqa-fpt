package user_server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm"
	"github.com/library/efk"
	"github.com/library/middleware"
	"github.com/library/models"
	password_hash "github.com/library/password-hash"
	"github.com/sirupsen/logrus"
)

func (srv *Server) register() http.HandlerFunc {
	return func(wr http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		w := middleware.NewLogResponseWriter(wr)
		account := &models.Account{}
		err := json.NewDecoder(r.Body).Decode(account)
		if err != nil {
			handleError(w, ctx, srv, "registration", err, http.StatusInternalServerError)
			return
		}
		// account.AccountRole = models.AdminAccount
		hashedPwd, err := password_hash.HashPassword(account.Password)
		if err != nil {
			handleError(w, ctx, srv, "registration", err, http.StatusInternalServerError)
			return
		}
		account.PasswordHash = hashedPwd
		err = srv.DB.CreateUserAccount(*account)
		if err != nil {
			if strings.Contains(err.Error(), "1062") {
				handleError(w, ctx, srv, "registration", err, http.StatusBadRequest)
				return
			}
			handleError(w, ctx, srv, "registration", err, http.StatusInternalServerError)
			return
		}
		// get the created user account
		acc, err := srv.DB.VerifyUser(*&models.LoginDetails{
			Email:       account.Email,
			Password:    account.Password,
			AccountRole: account.AccountRole,
		})
		if err != nil {
			handleError(w, ctx, srv, "registration", err, http.StatusInternalServerError)
			return
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":   acc.ID,
			"role": acc.AccountRole,
		})

		tokenStr, err := token.SignedString([]byte(srv.Env.JwtSigningKey))
		if err != nil {
			handleError(w, ctx, srv, "registration", err, http.StatusInternalServerError)
			return
		}
		logrus.WithFields(logrus.Fields{
			"statusCode": http.StatusOK,
		}).Info(fmt.Sprintf("new user registered with email: %v", account.Email))

		err = json.NewEncoder(w).Encode(&models.Response{AccountRole: account.AccountRole, Token: tokenStr})
		if err != nil {
			handleError(w, ctx, srv, "registration", err, http.StatusInternalServerError)
			return
		}
	}
}

func (srv *Server) login() http.HandlerFunc {
	return func(wr http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		w := &middleware.LogResponseWriter{ResponseWriter: wr}
		details := &models.LoginDetails{}
		err := json.NewDecoder(r.Body).Decode(details)
		if err != nil {
			handleError(w, ctx, srv, "login", err, http.StatusInternalServerError)
			return
		}

		account, err := srv.DB.VerifyUser(*details)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				handleError(w, ctx, srv, "login", errors.New(fmt.Sprintf("no such %v found", details.AccountRole)), http.StatusBadRequest)
			} else {
				handleError(w, ctx, srv, "login", err, http.StatusInternalServerError)
			}
			return
		}
		ok := password_hash.ValidatePassword(details.Password, account.PasswordHash)
		if !ok {
			handleError(w, ctx, srv, "login", errors.New("invalid password"), http.StatusUnauthorized)
			return
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":   account.ID,
			"role": account.AccountRole,
		})
		tokenStr, err := token.SignedString([]byte(srv.Env.JwtSigningKey))
		if err != nil {
			handleError(w, ctx, srv, "login", err, http.StatusInternalServerError)
			return
		}
		logrus.WithFields(logrus.Fields{
			"statusCode": http.StatusOK,
		}).Info(fmt.Sprintf("user login with email: %v", account.Email))
		err = json.NewEncoder(w).Encode(&models.Response{AccountRole: details.AccountRole, Token: tokenStr})
		if err != nil {
			handleError(w, ctx, srv, "login", err, http.StatusInternalServerError)
			return
		}
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
		efk.LogError(srv.EfkLogger, srv.EfkTag, srv.TracingID, task, err, statusCode)
	}
	http.Error(w, err.Error(), statusCode)

	logrus.WithFields(logrus.Fields{
		"tracingID":  srv.TracingID,
		"statusCode": statusCode,
		"error":      err,
	}).Error(task)
}

func GetAuthInfoFromContext(ctx context.Context) *models.AuthInfo {
	return ctx.Value(middleware.ContextAuthInfo).(*models.AuthInfo)
}

func (srv *Server) getUserByName(wr http.ResponseWriter, r *http.Request) {
	w := &middleware.LogResponseWriter{ResponseWriter: wr}
	ctx := r.Context()
	authInfo := GetAuthInfoFromContext(ctx)
	if authInfo.Role != models.AdminAccount {
		handleError(w, ctx, srv, "get_user_by_name", errors.New("permission denied"), http.StatusUnauthorized)
		return
	}
	userName := chi.URLParam(r, "name")
	users, err := srv.DB.GetUserByName(userName)
	if err != nil {
		if err == gorm.ErrRecordNotFound || len(*users) == 0 {
			handleError(w, ctx, srv, "get_user_by_name", errors.New("no record found"), http.StatusOK)
			return
		}
		handleError(w, ctx, srv, "get_user_by_name", err, http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		handleError(w, ctx, srv, "get_user_by_name", err, http.StatusInternalServerError)
	}
}

func (srv *Server) getUserByEmail(wr http.ResponseWriter, r *http.Request) {
	w := &middleware.LogResponseWriter{ResponseWriter: wr}
	ctx := r.Context()
	authInfo := GetAuthInfoFromContext(ctx)
	if authInfo.Role != models.AdminAccount {
		handleError(w, ctx, srv, "get_user_by_email", errors.New("permission denied"), http.StatusUnauthorized)
		return
	}
	userEmail := chi.URLParam(r, "email")

	users, err := srv.DB.GetUserByEmail(userEmail)
	if err != nil {
		if err == gorm.ErrRecordNotFound || len(*users) == 0 {
			handleError(w, ctx, srv, "get_user_by_email", errors.New("no record found"), http.StatusOK)
			return
		}
		handleError(w, ctx, srv, "get_user_by_email", err, http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		handleError(w, ctx, srv, "get_user_by_email", err, http.StatusInternalServerError)
	}
}

func (srv *Server) getUserByID(wr http.ResponseWriter, r *http.Request) {
	w := &middleware.LogResponseWriter{ResponseWriter: wr}
	ctx := r.Context()
	authInfo := GetAuthInfoFromContext(ctx)
	if authInfo.Role != models.AdminAccount {
		handleError(w, ctx, srv, "get_user_by_id", errors.New("permission denied"), http.StatusUnauthorized)
		return
	}
	id := chi.URLParam(r, "id")
	userID, err := strconv.Atoi(id)
	user, err := srv.DB.GetUserByID(uint(userID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleError(w, ctx, srv, "get_user_by_id", errors.New("no record found"), http.StatusOK)
			return
		}
		handleError(w, ctx, srv, "get_user_by_id", err, http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		handleError(w, ctx, srv, "get_user_by_id", err, http.StatusInternalServerError)
	}
}

func (srv *Server) getUsers(wr http.ResponseWriter, r *http.Request) {
	w := &middleware.LogResponseWriter{ResponseWriter: wr}
	ctx := r.Context()
	authInfo := GetAuthInfoFromContext(ctx)
	if authInfo.Role != models.AdminAccount {
		handleError(w, ctx, srv, "get_users", errors.New("permission denied"), http.StatusUnauthorized)
		return
	}
	users, err := srv.DB.GetUsers()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleError(w, ctx, srv, "get_users", errors.New("no record found"), http.StatusOK)
			return
		}
		handleError(w, ctx, srv, "get_users", err, http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		handleError(w, ctx, srv, "get_users", err, http.StatusInternalServerError)
	}
}
