package main

import (
	"flag"
	"os"

	"github.com/fluent/fluent-logger-golang/fluent"
	"github.com/golang/glog"
	"github.com/kelseyhightower/envconfig"
	management_server "github.com/library/cmd/management-svc/management-server"
	data_store "github.com/library/data-store"
	"github.com/library/envConfig"
	"github.com/library/middleware"
	"github.com/sirupsen/logrus"
)

var (
	dataStore *data_store.DataStore
	env       *envConfig.Env
	logger    *fluent.Fluent
	srv       *management_server.Server
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

	middleware.SetJwtSigningKey(env.JwtSigningKey)
	dataStore = data_store.DbConnect(env, testRun)

	srv = management_server.NewServer(env, dataStore, logger)
	err = srv.ListenAndServe("management-service", env.ManagementSvcPort)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("management-server start")
	}
}
