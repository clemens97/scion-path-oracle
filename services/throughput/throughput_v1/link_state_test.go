package throughput_v1

import (
	"github.com/clemens97/scion-path-oracle"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type pathStateMock struct {
	mock.Mock
}

func (l *pathStateMock) Score() float64 {
	args := l.Called()
	return args.Get(0).(float64)
}

func (l *pathStateMock) AddBw(_ float64) {
	l.Called()
}

func TestLinkScoreWithoutPath(t *testing.T) {
	lstate := newLinkState()
	fp := oracle.PathFingerprint("123")

	pstate := pathStateMock{}
	pstate.On("Score").Once().Return(0.)

	lstate.addPath(fp, &pstate)
	assert.Equal(t, 0., lstate.Score())
}

func TestLinkScore(t *testing.T) {
	lstate := newLinkState()

	pstate1 := pathStateMock{}
	pstate1.On("Score").Once().Return(2.)

	pstate2 := pathStateMock{}
	pstate2.On("Score").Once().Return(4.)

	lstate.addPath("123", &pstate1)
	lstate.addPath("456", &pstate2)
	assert.Equal(t, 4., lstate.Score())
}
