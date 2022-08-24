package graph

import (
	"github.com/scionproto/scion/go/lib/addr"
	"github.com/scionproto/scion/go/lib/snet"
	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/graph/encoding"
	"testing"
)

func TestNodeId(t *testing.T) {
	pi := snet.PathInterface{
		ID: 1,
		IA: addr.IA{I: 1, A: 1},
	}
	in := newInterfaceNode(pi)
	assert.Equal(t, in.ID(), NodeID(pi))
}

func TestNodeAttributes(t *testing.T) {
	in := newInterfaceNode(snet.PathInterface{
		ID: 1,
		IA: addr.IA{I: 2, A: 3},
	})
	attrs := in.Attributes()

	assert.Contains(t, attrs, encoding.Attribute{Key: "interface", Value: "1"})
	assert.Contains(t, attrs, encoding.Attribute{Key: "isd", Value: "2"})
	assert.Contains(t, attrs, encoding.Attribute{Key: "as", Value: "3"})
}
