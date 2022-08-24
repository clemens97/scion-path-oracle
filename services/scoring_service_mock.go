package services

import (
	"github.com/clemens97/scion-path-oracle"
	"github.com/scionproto/scion/go/lib/addr"
	"github.com/scionproto/scion/go/lib/snet"
	"github.com/stretchr/testify/mock"
)

type ScoringServiceMock struct {
	mock.Mock
}

func (s *ScoringServiceMock) OnStatsReceived(report oracle.Report, path snet.Path) {
	s.Called(report, path)
}

func (s *ScoringServiceMock) GetScores(dst addr.IA) PathScorings {
	s.Called(dst)
	return make(PathScorings)
}

func (s *ScoringServiceMock) Name() ServiceName {
	s.Called()
	return "mock"
}
