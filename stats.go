package oracle

import (
	"github.com/scionproto/scion/go/lib/addr"
)

// MonitoredProperties can contain any property a client measured of a connection, e.g. throughput, latency or packet loss.
type MonitoredProperties map[string]interface{}

type MetadataProperties map[string]interface{}

type Metadata struct {
	// Application which Connection MonitoredProperties are reported for, e.g. bittorrent.
	Application string `json:"application"`
	// Duration of the underlying Connection in seconds
	Duration float64 `json:"duration"`
	// Properties contain further information on the application's networking requirements and behaviour
	Properties MetadataProperties `json:"properties"`
}

type Report struct {
	// HTTP Request Body
	Metadata   Metadata            `json:"meta"`
	Properties MonitoredProperties `json:"stats"`
	// Path Parameters
	SrcIA  addr.IA         `json:"-"`
	DstIA  addr.IA         `json:"-"`
	PathFp PathFingerprint `json:"-"`
}
