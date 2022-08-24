package graph

import (
	"crypto/sha1"
	"encoding/binary"
	"github.com/scionproto/scion/go/lib/snet"
	"gonum.org/v1/gonum/graph/encoding"
	"strconv"
)

type InterfaceNode struct {
	i snet.PathInterface
}

func newInterfaceNode(i snet.PathInterface) InterfaceNode {
	return InterfaceNode{i: i}
}

func NodeID(inf snet.PathInterface) int64 {
	buf := make([]byte, 16)
	binary.LittleEndian.PutUint64(buf, uint64(inf.IA.IAInt()))
	binary.LittleEndian.PutUint64(buf[8:], uint64(inf.ID))
	h := sha1.New()
	h.Write(buf)
	return int64(binary.LittleEndian.Uint64(h.Sum(nil)))
}

func (ni InterfaceNode) ID() int64 {
	return NodeID(ni.i)
}

func (ni InterfaceNode) DOTID() string {
	return ni.i.String()
}

func (ni InterfaceNode) Attributes() []encoding.Attribute {
	return []encoding.Attribute{
		{Key: "as", Value: ni.i.IA.A.String()},
		{Key: "isd", Value: ni.i.IA.I.String()},
		{Key: "interface", Value: strconv.FormatUint(uint64(ni.i.ID), 10)},
	}
}
