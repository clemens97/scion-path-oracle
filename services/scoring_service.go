package services

import (
	"github.com/clemens97/scion-path-oracle"
	"github.com/scionproto/scion/go/lib/addr"
	"github.com/scionproto/scion/go/lib/snet"
)

type PathScorings map[oracle.PathFingerprint]float64

type ServiceName string

type ScoringService interface {
	OnStatsReceived(report oracle.Report, path snet.Path)
	GetScores(dst addr.IA) PathScorings
	Name() ServiceName
}
