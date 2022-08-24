package server

import (
	"github.com/scionproto/scion/go/lib/addr"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseIAFromNumeric(t *testing.T) {
	ia, err := parseIAFromNumeric("10", "20")
	assert.NoError(t, err)
	assert.Equal(t, ia, addr.IA{I: 10, A: 20})
}

func TestParseIAWithError(t *testing.T) {
	_, err := parseIAFromNumeric("19", "ffaa:1:e4b")
	assert.Error(t, err)
}
