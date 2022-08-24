package throughput_v2

import (
	"context"
	"fmt"
	"github.com/clemens97/scion-path-oracle"
	"github.com/clemens97/scion-path-oracle/graph"
	"github.com/clemens97/scion-path-oracle/querier"
	"github.com/clemens97/scion-path-oracle/services"
	"github.com/clemens97/scion-path-oracle/services/throughput"
	"github.com/gorilla/mux"
	"github.com/scionproto/scion/go/lib/addr"
	"github.com/scionproto/scion/go/lib/snet"
	"go.uber.org/zap"
	"net/http"
	"sync"
)

const name = "throughput"

type ScoringService struct {
	querier        querier.PathQuerier
	pathEMAService PathEMAService
	// knownFingerprints keeps track of all known path (fingerprints) for each known destination AS
	knownFingerprints map[addr.IA]oracle.FingerprintSet
	lock              sync.Mutex

	knownPaths sync.Map
	graph      graph.InterfaceGraph
	config     *config
	logger     *zap.SugaredLogger
}

func New(moniRouter *mux.Router, connector querier.PathQuerier, slogger *zap.SugaredLogger) *ScoringService {
	s := ScoringService{
		querier:           connector,
		knownFingerprints: make(map[addr.IA]oracle.FingerprintSet),
		graph:             graph.NewInterfaceGraph(),
	}

	s.logger = slogger.With("service_name", s.Name())
	s.config = newConfig()
	s.logger.Infow("service config loaded", "config", fmt.Sprintf("%+v", *s.config))
	s.pathEMAService = NewPathEMAService(s.config.pathEMASmoothing, s.logger)

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
	s.lock.Lock()
	defer s.lock.Unlock()

	if !s.filter(report) {
		return
	}

	isKnownDst := s.knownFingerprints[report.DstIA] != nil
	if !isKnownDst { // fetch all paths for dest
		paths, err := s.querier.Query(context.Background(), report.DstIA)
		if err != nil {
			s.logger.Errorw("could not query paths for destination of report",
				"destination", report.DstIA,
				"error", err)
			return
		}

		s.knownFingerprints[report.DstIA] = make(oracle.FingerprintSet)
		for _, p := range paths {
			fp := oracle.OracleFingerprint(p)
			s.knownFingerprints[report.DstIA][fp] = true
			s.knownPaths.Store(fp, p)
			s.addInterfacesWithStaticBw(*p.Metadata())
		}
		s.logger.Infow("fetched paths for report destination",
			"destination", report.DstIA,
			"fingerprints", s.knownFingerprints[report.DstIA])
	}

	bw, err := throughput.ParseThroughput(report.Properties)
	if err != nil {
		s.logger.Debugw("could not parse throughput from stats of report", "stats", report.Properties)
		return
	}

	s.pathEMAService.Update(bw, report.PathFp)

	meta := path.Metadata()
	for i := 0; i < len(meta.Interfaces)-1; i++ {
		firstIf := meta.Interfaces[i]
		secondIf := meta.Interfaces[i+1]

		staticBw := 0.
		if i < len(meta.Bandwidth) {
			staticBw = float64(meta.Bandwidth[i] * 1000)
		}

		s.addToGraph(firstIf, secondIf, staticBw)
		link := s.graph.Edge(graph.NodeID(firstIf), graph.NodeID(secondIf))
		link.State.(*linkState).addThroughput(bw, report.PathFp)
	}

}

func (s *ScoringService) filter(report oracle.Report) bool {
	profileProp, profExists := report.Metadata.Properties["taps-capacity-profile"]
	if !profExists {
		return false
	}
	profile, profileOk := profileProp.(string)
	return profileOk && profile == "capacity-seeking"
}

func (s *ScoringService) Name() services.ServiceName {
	return name
}

func (s *ScoringService) GetScores(dst addr.IA) services.PathScorings {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.logger.Debugw("getting scores", "dst_ia", dst)
	scorings := make(services.PathScorings)
	for fp := range s.knownFingerprints[dst] {
		//path := s.knownPaths[fp]
		path, ok := s.knownPaths.Load(fp)
		if !ok {
			continue
		}
		sp, ok := path.(snet.Path)
		if !ok {
			continue
		}

		scorings[fp] = s.getScore(sp)
	}
	return scorings
}

// getScore returns the minimum of all LinkScores on the PathFp
func (s *ScoringService) getScore(path snet.Path) float64 {
	meta := path.Metadata()
	scoreFirstLink := s.graph.Edge(graph.NodeID(meta.Interfaces[0]), graph.NodeID(meta.Interfaces[1])).State.Score()
	currentScore := scoreFirstLink

	for i := 1; i < len(meta.Interfaces)-1; i++ {
		link := s.graph.Edge(graph.NodeID(meta.Interfaces[i]), graph.NodeID(meta.Interfaces[i+1]))
		if link == nil {
			continue
		}
		if linkScore := link.State.Score(); linkScore < currentScore {
			currentScore = linkScore
		}
	}
	return currentScore
}

func (s *ScoringService) addInterfacesWithStaticBw(meta snet.PathMetadata) {

	for i := 0; i < len(meta.Interfaces)-1; i++ {
		staticBw := 0.
		if i < len(meta.Bandwidth) {
			staticBw = float64(meta.Bandwidth[i] * 1000)
		}
		s.addToGraph(meta.Interfaces[i], meta.Interfaces[i+1], staticBw)
	}
}

func (s *ScoringService) addToGraph(first, second snet.PathInterface, staticBw float64) {
	nodeF, _ := s.graph.GetOrCreateNode(first)
	nodeS, _ := s.graph.GetOrCreateNode(second)
	link := s.graph.Edge(nodeF.ID(), nodeS.ID())
	if link == nil {
		state := newLinkState(s.config, staticBw, &s.pathEMAService, s.logger)
		s.graph.SetEdge(graph.NewInterfaceLink(nodeF, nodeS, state))
	}
}
