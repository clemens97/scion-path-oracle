package graph

type LinkState interface {
	Score() float64
	GetConfidence() float64
}
