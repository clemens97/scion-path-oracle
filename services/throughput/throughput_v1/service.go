package throughput_v1

import (
	"fmt"
	"github.com/clemens97/scion-path-oracle"
	"github.com/clemens97/scion-path-oracle/graph"
	"github.com/clemens97/scion-path-oracle/services"
	"github.com/clemens97/scion-path-oracle/services/throughput"
	"github.com/gorilla/mux"
	"github.com/scionproto/scion/go/lib/addr"
	"github.com/scionproto/scion/go/lib/snet"
	"go.uber.org/zap"
	"net/http"
)

type ScoringService struct {
	graph  graph.InterfaceGraph
	paths  map[addr.IA]*oracle.FingerprintSet
	states map[oracle.PathFingerprint]*pathState
	logger *zap.SugaredLogger
}

func New(moniRouter *mux.Router, logger *zap.SugaredLogger) *ScoringService {
	s := ScoringService{
		graph:  graph.NewInterfaceGraph(),
		paths:  make(map[addr.IA]*oracle.FingerprintSet),
		states: make(map[oracle.PathFingerprint]*pathState),
	}
	s.logger = logger.With("service_name", s.Name())

	if moniRouter != nil {
		monitoringPath := string("/" + s.Name())
		s.logger.Debugw("attaching service status monitoring", "path", monitoringPath)
		moniRouter.HandleFunc(monitoringPath, s.internalStatus)
	}

	return &s
}

func (s *ScoringService) internalStatus(w http.ResponseWriter, _ *http.Request) {
	const monitorTemplate = `<html><img src="data:image/svg+xml;base64,%s" /></html>`
	svgGraph, err := graph.SvgB64Buffer(s.graph)
	if err != nil {
		s.logger.Errorw("could not create monitoring graph", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, monitorTemplate, svgGraph)
}

func (s *ScoringService) OnStatsReceived(report oracle.Report, path snet.Path) {
	bw, err := throughput.ParseThroughput(report.Properties)
	if err != nil {
		s.logger.Debugw("could not parse throughput from stats of report", "stats", report.Properties)
		return
	}

	knownPaths := s.paths[report.DstIA]
	if knownPaths == nil {
		fps := make(oracle.FingerprintSet)
		knownPaths = &fps
		s.paths[report.DstIA] = &fps
	}
	(*knownPaths)[report.PathFp] = true

	state := s.states[report.PathFp]
	if state == nil {
		ps := newPathState(float64(throughput.Min(path.Metadata().Bandwidth) * 1000))
		state = &ps
		s.states[report.PathFp] = state
	}
	state.AddBw(bw)

	meta := path.Metadata()
	previous, _ := s.graph.GetOrCreateNode(meta.Interfaces[0])
	for i := 1; i < len(meta.Interfaces); i++ {
		current, _ := s.graph.GetOrCreateNode(meta.Interfaces[i])
		link := s.graph.Edge(previous.ID(), current.ID())
		if link == nil {
			lState := newLinkState()
			iLink := graph.NewInterfaceLink(previous, current, &lState)
			s.graph.SetEdge(iLink)
			link = &iLink
		}
		link.State.(*linkState).addPath(report.PathFp, state)
	}

}

func (s *ScoringService) GetScores(dst addr.IA) services.PathScorings {
	scorings := make(services.PathScorings)
	knownPaths := *s.paths[dst]

	for fp := range knownPaths {
		state := s.states[fp]
		scorings[fp] = state.Score()
	}

	return scorings
}

func (s *ScoringService) Name() services.ServiceName {
	return "throughput_v1"
}
