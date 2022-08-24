package throughput_v2

import (
	"github.com/clemens97/scion-path-oracle"
	"github.com/clemens97/scion-path-oracle/querier"
	"github.com/scionproto/scion/go/lib/addr"
	"github.com/scionproto/scion/go/lib/common"
	"github.com/scionproto/scion/go/lib/snet"
	"github.com/scionproto/scion/go/lib/snet/path"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"testing"
)

func TestOnStatsReceived(t *testing.T) {
	//    IA SRC           IA I          IA DST         reported
	// A: ([src] -> E1) -> (I1 -> E1) -> (I1 -> [dst])  100
	// B: ([src] -> E2) -> (I1 -> E1) -> (I1 -> [dst])   50
	// C: ([src] -> E1) -> (I1 -> E2) -> (I1 -> [dst])    1
	// D: ([src] -> E2) -> (I1 -> E2) -> (I1 -> [dst])    -

	iaSrc := addr.IA{I: 0, A: 0}
	srcE1 := snet.PathInterface{ID: common.IFIDType(1), IA: iaSrc}
	srcE2 := snet.PathInterface{ID: common.IFIDType(2), IA: iaSrc}

	iaI := addr.IA{I: 0, A: 1}
	iI1 := snet.PathInterface{ID: common.IFIDType(1), IA: iaI}
	iE1 := snet.PathInterface{ID: common.IFIDType(2), IA: iaI}
	iE2 := snet.PathInterface{ID: common.IFIDType(3), IA: iaI}

	iaDst := addr.IA{I: 0, A: 2}
	dstI1 := snet.PathInterface{ID: common.IFIDType(1), IA: iaDst}

	pathA := path.Path{Dst: iaDst, Meta: snet.PathMetadata{Interfaces: []snet.PathInterface{srcE1, iI1, iE1, dstI1}}}
	pathB := path.Path{Dst: iaDst, Meta: snet.PathMetadata{Interfaces: []snet.PathInterface{srcE2, iI1, iE1, dstI1}}}
	pathC := path.Path{Dst: iaDst, Meta: snet.PathMetadata{Interfaces: []snet.PathInterface{srcE1, iI1, iE2, dstI1}}}
	pathD := path.Path{Dst: iaDst, Meta: snet.PathMetadata{Interfaces: []snet.PathInterface{srcE2, iI1, iE2, dstI1}}}

	fpA := oracle.OracleFingerprint(pathA)
	fpB := oracle.OracleFingerprint(pathB)
	fpC := oracle.OracleFingerprint(pathC)
	fpD := oracle.OracleFingerprint(pathD)

	reportA := getTestReportWithBw(100, iaSrc, iaDst, fpA)
	reportB := getTestReportWithBw(50, iaSrc, iaDst, fpB)
	reportC := getTestReportWithBw(1, iaSrc, iaDst, fpC)

	querier := querier.QuerierMock{}
	querier.On("Query", mock.Anything, mock.Anything).
		Return([]snet.Path{pathA, pathB, pathC, pathD}, nil).
		Once()

	service := New(nil, &querier, zap.S())
	//service.querier = &querier
	service.OnStatsReceived(reportA, pathA)
	service.OnStatsReceived(reportB, pathB)
	service.OnStatsReceived(reportC, pathC)
	scores := service.GetScores(iaDst)

	assert.Equal(t, service.graph.Nodes().Len(), 6)
	assert.Greater(t, scores[fpA], scores[fpB])
	assert.Greater(t, scores[fpB], scores[fpC])
	assert.NotZero(t, scores[fpD])
	querier.AssertExpectations(t)
}

func getTestReportWithBw(bw int, src, dst addr.IA, fp oracle.PathFingerprint) oracle.Report {
	stats := make(map[string]interface{})
	stats["throughput"] = bw

	metaProps := make(oracle.MetadataProperties)
	metaProps["taps-capacity-profile"] = "capacity-seeking"

	meta := oracle.Metadata{Properties: metaProps}

	return oracle.Report{
		Metadata:   meta,
		Properties: stats,
		SrcIA:      src,
		DstIA:      dst,
		PathFp:     fp,
	}
}
