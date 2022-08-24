package graph

import (
	"github.com/scionproto/scion/go/lib/addr"
	"github.com/scionproto/scion/go/lib/common"
	"github.com/scionproto/scion/go/lib/snet"
	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/graph/encoding"
	"testing"
)

func TestLinkAttributes(t *testing.T) {
	from := interfaceNodeOf(0, 1, 2)
	to := interfaceNodeOf(3, 4, 5)
	linkStateMock := LinkStateMock{}
	linkStateMock.On("Score").Once().Return(0.123)

	interfaceLink := NewInterfaceLink(&from, &to, &linkStateMock)
	attrs := interfaceLink.Attributes()

	assert.Contains(t, attrs, encoding.Attribute{Key: "score", Value: "0.12"})
	linkStateMock.AssertExpectations(t)
}

func interfaceNodeOf(ifId common.IFIDType, i addr.ISD, a addr.AS) InterfaceNode {
	return newInterfaceNode(snet.PathInterface{
		ID: ifId,
		IA: addr.IA{I: i, A: a},
	})
}
