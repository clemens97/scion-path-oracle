package graph

import (
	"github.com/scionproto/scion/go/lib/addr"
	"github.com/scionproto/scion/go/lib/snet"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetOrCreate_Created(t *testing.T) {
	graph := NewInterfaceGraph()
	inf := snet.PathInterface{ID: 1, IA: addr.IA{I: 2, A: 3}}
	n, created := graph.GetOrCreateNode(inf)

	assert.NotNil(t, n)
	assert.True(t, created)
	assert.Equal(t, 1, graph.Nodes().Len())
}

func TestGetOrCreate_NotCreated(t *testing.T) {
	graph := NewInterfaceGraph()
	inf := snet.PathInterface{ID: 1, IA: addr.IA{I: 2, A: 3}}
	node := newInterfaceNode(inf)
	graph.AddNode(node)

	n, created := graph.GetOrCreateNode(inf)

	assert.Equal(t, n, &node)
	assert.False(t, created)
}

func TestEdge(t *testing.T) {
	graph := NewInterfaceGraph()
	from := newInterfaceNode(snet.PathInterface{ID: 1, IA: addr.IA{I: 2, A: 3}})
	to := newInterfaceNode(snet.PathInterface{ID: 4, IA: addr.IA{I: 5, A: 6}})
	state := LinkStateMock{}

	edge := NewInterfaceLink(&from, &to, &state)
	assert.Nil(t, graph.Edge(from.ID(), to.ID()))

	graph.SetEdge(edge)
	assert.Equal(t, &edge, graph.Edge(from.ID(), to.ID()))
}
