package server

import (
	"github.com/clemens97/scion-path-oracle/querier"
	"github.com/clemens97/scion-path-oracle/services"
	"github.com/clemens97/scion-path-oracle/services/throughput/throughput_v2"
	"github.com/gorilla/mux"
	"github.com/netsec-ethz/scion-apps/pkg/shttp"
	"go.uber.org/zap"
	"net/http"
)

type OracleServer struct {
	querier               querier.PathQuerier
	router, routerMonitor *mux.Router
	scoringServices       map[services.ServiceName]*services.ScoringService
	listen, listenMonitor string
	logger                *zap.SugaredLogger
}

const (
	dstIsdKey = "dst_isd"
	dstAsKey  = "dst_as"
	pathKey   = "path_fp"
)

func New(listen, listenMonitor string, logger *zap.SugaredLogger) (*OracleServer, error) {
	q, err := querier.New()
	if err != nil {
		return nil, err
	}

	os := OracleServer{
		router:        mux.NewRouter(),
		routerMonitor: mux.NewRouter(),
		listen:        listen,
		listenMonitor: listenMonitor,
		querier:       q,
		logger:        logger,
	}

	// os.AddScoringService(throughput_v1.New(os.routerMonitor, logger))
	os.AddScoringService(throughput_v2.New(os.routerMonitor, q, logger))
	os.RegisterScoringAPI()
	os.RegisterStatsAPI()

	return &os, nil
}

func (o *OracleServer) Start() error {
	go func() {
		err := http.ListenAndServe(o.listenMonitor, o.routerMonitor)
		o.logger.Errorw("monitoring http error", "error", err)
	}()
	return shttp.ListenAndServe(o.listen, o.router)
}

func (o *OracleServer) AddScoringService(service services.ScoringService) {
	o.logger.Infow("adding scoring service", "service name", service.Name())
	if o.scoringServices == nil {
		o.scoringServices = make(map[services.ServiceName]*services.ScoringService)
	}
	if o.scoringServices[service.Name()] != nil {
		o.logger.Warnw("attempting to re-register service, old service instance will be overwritten",
			"service name", service.Name())
	}
	o.scoringServices[service.Name()] = &service
}
