package graph

import (
	"github.com/scionproto/scion/go/lib/snet"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"
)

type InterfaceGraph struct {
	g *simple.DirectedGraph
}

func NewInterfaceGraph() InterfaceGraph {
	return InterfaceGraph{g: simple.NewDirectedGraph()}
}

func (i *InterfaceGraph) GetOrCreateNode(p snet.PathInterface) (*InterfaceNode, bool) {
	if node := i.Node(NodeID(p)); node != nil {
		return node, false
	}
	node := newInterfaceNode(p)
	i.AddNode(node)
	return &node, true
}

func (i *InterfaceGraph) AddNode(node InterfaceNode) {
	i.g.AddNode(node)
}

func (i *InterfaceGraph) Node(id int64) *InterfaceNode {
	node := i.g.Node(id)
	if node == nil {
		return nil
	}

	if n, ok := node.(InterfaceNode); ok {
		return &n
	}
	panic("node of InterfaceGraph holds unsupported node")
}

func (i *InterfaceGraph) Edge(from, to int64) *InterfaceLink {
	edge := i.g.Edge(from, to)
	if edge == nil {
		return nil
	}
	if l, ok := edge.(InterfaceLink); ok {
		return &l
	}
	panic("node of InterfaceGraph holds unsupported edge")
}

func (i *InterfaceGraph) Nodes() graph.Nodes {
	return i.g.Nodes()
}

func (i *InterfaceGraph) SetEdge(link InterfaceLink) {
	i.g.SetEdge(link)
}

func (i *InterfaceGraph) SerializeDot() ([]byte, error) {
	return dot.Marshal(i.g, "Network Graph", "", "\t")
}
