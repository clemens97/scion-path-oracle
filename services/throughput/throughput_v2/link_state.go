package throughput_v2

import (
	"github.com/clemens97/scion-path-oracle"
	"github.com/clemens97/scion-path-oracle/services/throughput"
	"go.uber.org/zap"
	"math"
	"sort"
	"sync"
)

type throughputEntry struct {
	throughput float64
	fp         oracle.PathFingerprint
}

type linkState struct {
	reportedBws    []throughputEntry
	staticBw       float64
	pathEMAService *PathEMAServiceI
	m              sync.Mutex
	config         *config
	logger         *zap.SugaredLogger
}

func newLinkState(config *config, staticBw float64, pathEMAService PathEMAServiceI, logger *zap.SugaredLogger) *linkState {
	return &linkState{
		config:         config,
		reportedBws:    make([]throughputEntry, 0),
		staticBw:       staticBw,
		pathEMAService: &pathEMAService,
		logger:         logger}
}

// addThroughput appends a newly reported throughput to the linkState
func (s *linkState) addThroughput(throughput float64, fp oracle.PathFingerprint) {
	s.m.Lock()
	defer s.m.Unlock()

	e := throughputEntry{throughput: throughput, fp: fp}
	s.logger.Debugw("adding throughput to link", "throughput", throughput, "fingerprint", fp)
	if len(s.reportedBws) < s.config.maxReportAmount {
		s.reportedBws = append(s.reportedBws, e)
	} else {
		s.logger.Debugw("dropping oldest throughput from link")
		s.reportedBws = append(s.reportedBws[1:], e)
	}
}

// Score calculates the average of all non-expired throughput_v2 reports
func (s *linkState) Score() float64 {
	s.m.Lock()
	defer s.m.Unlock()

	if len(s.reportedBws) == 0 {
		s.logger.Debugw("no reports present to calculate link score - falling back to static throughput",
			"static_bandwidth", s.staticBw)
		return s.staticBw
	}

	fps := s.getPathsToConsider()
	entries := make([]float64, 0)
	for _, e := range s.reportedBws {
		if fps[e.fp] {
			entries = append(entries, e.throughput)
		}
	}

	return throughput.SliceEMA(entries, s.config.linkEMASmoothing)
}

func (s *linkState) getPathsToConsider() oracle.FingerprintSet {
	pathEMAs := make(map[oracle.PathFingerprint]float64)
	for _, e := range s.reportedBws {
		pathEMAs[e.fp] = (*s.pathEMAService).GetEMAOrDefault(e.fp, 0)
	}

	candidates := make([]throughputEntry, 0)

	for k, v := range pathEMAs {
		candidates = append(candidates, throughputEntry{throughput: v, fp: k})
	}

	sort.Slice(candidates, func(i, j int) bool {
		// descending order
		return candidates[i].throughput > candidates[j].throughput
	})

	pathsToConsider := int(math.Ceil(float64(len(candidates)) * (1 - s.config.bottleNeckRatio)))
	fpsToConsider := make(oracle.FingerprintSet, pathsToConsider)

	for i := 0; i < pathsToConsider; i++ {
		fpsToConsider[candidates[i].fp] = true
	}

	s.logger.Debugw("paths considered for link score", "amount_candidates", len(candidates),
		"considered_fingerprints", fpsToConsider, "bottleNeckRatio", s.config.bottleNeckRatio)

	return fpsToConsider
}

func (s *linkState) GetConfidence() float64 {
	s.m.Lock()
	defer s.m.Unlock()
	return float64(len(s.reportedBws)) / float64(s.config.maxReportAmount)
}
