package throughput_v1

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestScore(t *testing.T) {
	staticBw := 1.
	reportedBw1 := 8.
	reportedBw2 := 9.

	pstate := newPathState(staticBw)
	pstate.AddBw(reportedBw1)
	pstate.AddBw(reportedBw2)

	assert.Equal(t, 6., pstate.Score())
}

func TestScoreExpire(t *testing.T) {
	staticBw := 100.
	pstate := newPathState(staticBw)

	for i := 0; i < maxReportAmount; i++ {
		pstate.AddBw(1)
	}

	assert.Equal(t, 1., pstate.Score())
}
