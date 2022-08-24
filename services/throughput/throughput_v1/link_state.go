package throughput_v1

import (
	"github.com/clemens97/scion-path-oracle"
	"sync"
)

type linkState struct {
	pathStates map[oracle.PathFingerprint]pathStateI
	m          sync.Mutex
}

func newLinkState() linkState {
	return linkState{pathStates: make(map[oracle.PathFingerprint]pathStateI)}
}

// addPath adds a path which utilises this links
func (s *linkState) addPath(fp oracle.PathFingerprint, state pathStateI) {
	s.m.Lock()
	defer s.m.Unlock()
	s.pathStates[fp] = state
}

// Score returns the maximum of all path scores
func (s *linkState) Score() float64 {
	s.m.Lock()
	defer s.m.Unlock()

	lscore := 0.
	for _, state := range s.pathStates {
		if pscore := state.Score(); pscore > lscore {
			lscore = pscore
		}
	}
	return lscore
}

func (s *linkState) GetConfidence() float64 {
	return 1
}
