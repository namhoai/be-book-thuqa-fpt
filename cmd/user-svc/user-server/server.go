package user_server

import (
	"net/http"

	datastore "github.com/library/data-store"
	"github.com/library/envConfig"
	"github.com/library/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

var (
	prom        *prometheus.Registry
	promMetrics *metrics.Metrics
)

type Server struct {
	DB        datastore.DbUtil
	Env       *envConfig.Env
	TracingID string
	TestRun   bool
}

func NewServer(env *envConfig.Env, db datastore.DbUtil) *Server {
	return &Server{
		DB:        db,
		Env:       env,
		TracingID: "",
		TestRun:   false,
	}
}

func (srv *Server) ListenAndServe(service string, port string) error {
	prom = prometheus.NewRegistry()
	promMetrics = metrics.NewMetrics("user_svc")
	prom.MustRegister(promMetrics.RequestCounter)
	prom.MustRegister(promMetrics.LatencyCalculator)

	r := SetupRouter(srv, prom)
	logrus.WithFields(logrus.Fields{
		"service": service,
	}).Info(service+" binding on ", ":"+port)

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		return err
	}
	return nil
}
