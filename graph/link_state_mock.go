package graph

import "github.com/stretchr/testify/mock"

type LinkStateMock struct {
	mock.Mock
}

func (l *LinkStateMock) Score() float64 {
	args := l.Called()
	return args.Get(0).(float64)
}

func (l *LinkStateMock) GetConfidence() float64 {
	args := l.Called()
	return args.Get(0).(float64)
}
