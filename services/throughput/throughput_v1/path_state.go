package throughput_v1

import "sync"

const maxReportAmount = 20

type pathStateI interface {
	Score() float64
	AddBw(bw float64)
}

type pathState struct {
	reportedBws []float64 // index 0: oldest report, last index: most recent report
	m           sync.Mutex
}

func newPathState(staticBw float64) pathState {
	return pathState{reportedBws: []float64{staticBw}}
}

// AddBw appends a newly reported throughput to the linkState
func (s *pathState) AddBw(bw float64) {
	s.m.Lock()
	defer s.m.Unlock()
	if len(s.reportedBws) < maxReportAmount {
		s.reportedBws = append(s.reportedBws, bw)
		return
	}
	s.reportedBws = append(s.reportedBws[1:], bw)
}

func (s *pathState) Score() float64 {
	s.m.Lock()
	defer s.m.Unlock()

	sum := 0.
	for _, bw := range s.reportedBws {
		sum += bw
	}
	return sum / float64(len(s.reportedBws))
}
