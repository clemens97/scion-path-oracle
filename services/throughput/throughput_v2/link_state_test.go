package throughput_v2

import (
	oracle "github.com/clemens97/scion-path-oracle"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"testing"
)

func TestExpire(t *testing.T) {
	pathEmaM := PathEMAMock{}
	pathEmaM.On("GetEMAOrDefault", mock.Anything, mock.Anything).Return(0.)

	lstate := newLinkState(testConfig(), 100., &pathEmaM, zap.S())

	for i := 0; i < lstate.config.maxReportAmount; i++ {
		lstate.addThroughput(1, "")
	}

	assert.Equal(t, 1., lstate.Score())
	assert.Equal(t, 1., lstate.GetConfidence())
}

func TestFreshLState(t *testing.T) {
	lstate := &linkState{config: testConfig(), logger: zap.S()}

	assert.Equal(t, 0., lstate.Score())
	assert.Equal(t, 0., lstate.GetConfidence())
}

func TestGetPathsToConsider(t *testing.T) {
	pathEmaM := PathEMAMock{}
	pathEmaM.On("GetEMAOrDefault", oracle.PathFingerprint("p1"), mock.Anything).Return(1.).Twice()
	pathEmaM.On("GetEMAOrDefault", oracle.PathFingerprint("p2"), mock.Anything).Return(2.).Twice()
	lstate := newLinkState(testConfig(), 0, &pathEmaM, zap.S())

	lstate.reportedBws = []throughputEntry{{100, "p1"}, {200, "p1"}, {300, "p2"}, {300, "p2"}}
	ptc := lstate.getPathsToConsider()

	assert.False(t, ptc["p1"])
	assert.True(t, ptc["p2"])
	pathEmaM.AssertExpectations(t)
}

func testConfig() *config {
	return &config{maxReportAmount: 10, linkEMASmoothing: 0.5, pathEMASmoothing: 0.5, bottleNeckRatio: 0.5}
}
