package querier

import (
	"context"
	"github.com/clemens97/scion-path-oracle"
	"github.com/scionproto/scion/go/lib/addr"
	"github.com/scionproto/scion/go/lib/snet"
	"github.com/stretchr/testify/mock"
)

type QuerierMock struct {
	mock.Mock
}

func (q *QuerierMock) Query(ctx context.Context, ia addr.IA) ([]snet.Path, error) {
	args := q.Called(ctx, ia)
	paths := args.Get(0).([]snet.Path)
	err := args.Get(1)
	if err != nil {
		return paths, err.(error)
	}
	return paths, nil
}

func (q *QuerierMock) GetPath(ia addr.IA, fp oracle.PathFingerprint) (snet.Path, error) {
	args := q.Called(ia, fp)
	path := args.Get(0).(snet.Path)
	err := args.Get(1)
	if err != nil {
		return path, err.(error)
	}
	return path, nil
}

func (q *QuerierMock) LocalIA() addr.IA {
	args := q.Called()
	return args.Get(0).(addr.IA)
}
