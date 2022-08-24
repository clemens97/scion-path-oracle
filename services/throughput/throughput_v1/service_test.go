package throughput_v1

import (
	"github.com/clemens97/scion-path-oracle"
	"github.com/scionproto/scion/go/lib/addr"
	"github.com/scionproto/scion/go/lib/common"
	"github.com/scionproto/scion/go/lib/snet"
	"github.com/scionproto/scion/go/lib/snet/path"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestOnStatsReceived(t *testing.T) {
	anySrc := addr.IA{I: 0, A: 0}
	anyDst := addr.IA{I: 1, A: 0}

	anyPathFp := oracle.PathFingerprint("123")
	bw1 := 5000
	bw2 := 10000
	report := getTestReportWithBw(bw1, anySrc, anyDst, anyPathFp)
	anotherReport := getTestReportWithBw(bw2, anySrc, anyDst, anyPathFp)

	if1 := snet.PathInterface{ID: common.IFIDType(1), IA: anySrc}
	if2 := snet.PathInterface{ID: common.IFIDType(2), IA: anyDst}

	pathForReport := path.Path{Dst: anyDst, Meta: snet.PathMetadata{Interfaces: []snet.PathInterface{if1, if2}}}

	service := New(nil, zap.S())
	service.OnStatsReceived(report, pathForReport)
	service.OnStatsReceived(anotherReport, pathForReport)

	// expected average throughput (avg: 10k, 5k, and 0 (no static bw configured))
	assert.Equal(t, 5000., service.GetScores(report.DstIA)[report.PathFp])
	assert.Equal(t, 2, service.graph.Nodes().Len())
}

func getTestReportWithBw(bw int, src, dst addr.IA, fp oracle.PathFingerprint) oracle.Report {
	stats := make(map[string]interface{})
	stats["throughput"] = bw

	return oracle.Report{
		Properties: stats,
		SrcIA:      src,
		DstIA:      dst,
		PathFp:     fp,
	}
}
