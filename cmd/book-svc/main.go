package main

import (
	"flag"
	"os"

	"github.com/fluent/fluent-logger-golang/fluent"
	"github.com/golang/glog"
	"github.com/kelseyhightower/envconfig"
	book_server "github.com/library/cmd/book-svc/book-server"
	data_store "github.com/library/data-store"
	"github.com/library/envConfig"
	"github.com/library/middleware"
	"github.com/sirupsen/logrus"
)

var (
	dataStore *data_store.DataStore
	env       *envConfig.Env
	logger    *fluent.Fluent
	srv       *book_server.Server
	testRun   bool
)

func init() {
	testRun = false
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
}

func main() {
	flag.Parse()
	env = &envConfig.Env{}
	err := envconfig.Process("LIBRARY", env)
	if err != nil {
		glog.Fatal(err)
	}

	dataStore = data_store.DbConnect(env, testRun)
	middleware.SetJwtSigningKey(env.JwtSigningKey)

	srv = book_server.NewServer(env, dataStore, logger)
	err = srv.ListenAndServe("book-service", env.BookSvcPort)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("book-server start")
	}
}
