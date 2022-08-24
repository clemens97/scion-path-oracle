package graph

import (
	"fmt"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
)

type InterfaceLink struct {
	from, to *InterfaceNode
	State    LinkState
}

func NewInterfaceLink(from, to *InterfaceNode, state LinkState) InterfaceLink {
	return InterfaceLink{from: from, to: to, State: state}
}

func (e InterfaceLink) From() graph.Node {
	return *e.from
}

func (e InterfaceLink) To() graph.Node {
	return *e.to
}

func (e InterfaceLink) ReversedEdge() graph.Edge {
	return InterfaceLink{to: e.from, from: e.to}
}

func (e InterfaceLink) Attributes() []encoding.Attribute {
	score := fmt.Sprintf("%.2f", (e.State).Score())
	return []encoding.Attribute{
		{Key: "label", Value: fmt.Sprintf(" score:%s", score)},
		{Key: "score", Value: score},
	}
}
