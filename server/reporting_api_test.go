package server

import (
	"github.com/clemens97/scion-path-oracle"
	"github.com/clemens97/scion-path-oracle/querier"
	"github.com/clemens97/scion-path-oracle/services"
	"github.com/gorilla/mux"
	"github.com/scionproto/scion/go/lib/addr"
	"github.com/scionproto/scion/go/lib/snet"
	"github.com/scionproto/scion/go/lib/snet/path"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type FetcherMock struct {
	mock.Mock
}

func (f *FetcherMock) GetPath(dst addr.IA, fp oracle.PathFingerprint) (snet.Path, error) {
	args := f.Called(dst, fp)
	return args.Get(0).(snet.Path), nil
}

func (f *FetcherMock) LocalIA() addr.IA {
	f.Called()
	return addr.IA{I: 0, A: 0}
}

func TestBadRequest(t *testing.T) {
	writer := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/stats/10/11/10/12/a-b-c", nil)

	o := OracleServer{logger: zap.S()}
	o.handleQueryScoring(writer, request)

	result := writer.Result()
	defer result.Body.Close()

	assert.Equal(t, http.StatusBadRequest, result.StatusCode)
}

func TestGoodRequest(t *testing.T) {
	srcStr := "19-ffaa:1:e94,[127.0.0.1]:1212"
	srcAddr, _ := snet.ParseUDPAddr(srcStr)
	dstIA := addr.IA{I: 1, A: 13}
	fp := oracle.PathFingerprint("a->b->c")

	writer := httptest.NewRecorder()
	requestBody := "{\"throughput\": 200}"
	request := httptest.NewRequest(http.MethodPost, "/stats", strings.NewReader(requestBody))
	request.RemoteAddr = srcStr

	pathVars := map[string]string{
		"dst_isd": dstIA.I.String(),
		"dst_as":  dstIA.A.String(),
		"path_fp": string(fp),
	}
	request = mux.SetURLVars(request, pathVars)

	scoringService := services.ScoringServiceMock{}
	scoringService.On("OnStatsReceived", mock.AnythingOfType("oracle.Report"), mock.AnythingOfType("path.Path")).Once()
	scoringService.On("Name")

	querier := &querier.QuerierMock{}
	querier.On("LocalIA").Once().Return(srcAddr.IA)
	querier.On("GetPath", dstIA, fp).Once().Return(path.Path{Dst: dstIA}, nil)

	o := OracleServer{querier: querier, logger: zap.S()}
	o.AddScoringService(&scoringService)
	o.handleReport(writer, request)

	result := writer.Result()
	defer result.Body.Close()

	assert.Equal(t, http.StatusCreated, result.StatusCode)
	scoringService.AssertExpectations(t)
}
