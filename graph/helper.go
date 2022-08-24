package graph

import (
	"bytes"
	"encoding/base64"
	"github.com/goccy/go-graphviz"
)

func SvgB64Buffer(graph InterfaceGraph) (string, error) {
	dotg, err := graph.SerializeDot()
	if err != nil {
		return "", err
	}
	viz, err := graphviz.ParseBytes(dotg)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = graphviz.New().Render(viz, graphviz.SVG, &buf)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}
