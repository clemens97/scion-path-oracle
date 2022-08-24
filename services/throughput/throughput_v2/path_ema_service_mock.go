package throughput_v2

import (
	"github.com/clemens97/scion-path-oracle"
	"github.com/stretchr/testify/mock"
)

type PathEMAMock struct {
	mock.Mock
}

func (p *PathEMAMock) Update(val float64, fp oracle.PathFingerprint) {
	p.Called(val, fp)
}

func (p *PathEMAMock) GetEMAOrDefault(fp oracle.PathFingerprint, def float64) float64 {
	args := p.Called(fp, def)
	return args.Get(0).(float64)
}
