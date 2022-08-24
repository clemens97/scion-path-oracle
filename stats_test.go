package oracle

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecoding(t *testing.T) {
	jsonStatsReport := `{
  		"meta": {
    		"application": "bittorrent",
    		"duration": 120,
			"properties": {
			  "protocols": [
				"SCION",
				"UDP",
				"QUIC"
			  ],
			"taps-capacity-profile": "capacity-seeking"
    		}
		},
		"stats": {
    		"throughput": 5000,
    		"loss": 0.02
  		}
	}`
	var report Report

	assert.NoError(t, json.Unmarshal([]byte(jsonStatsReport), &report))
	assert.Equal(t, 5000., report.Properties["throughput"])
	assert.Equal(t, "bittorrent", report.Metadata.Application)
	assert.Contains(t, report.Metadata.Properties["protocols"], "UDP")
	assert.Equal(t, "capacity-seeking", report.Metadata.Properties["taps-capacity-profile"])
}
