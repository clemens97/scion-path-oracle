package oracle

import (
	"fmt"
	"github.com/netsec-ethz/scion-apps/pkg/pan"
	"github.com/scionproto/scion/go/lib/snet"
	"strings"
)

// PathFingerprint is a unique identifier for a path
type PathFingerprint string

type FingerprintSet map[PathFingerprint]bool

func OracleFingerprint(p snet.Path) PathFingerprint {
	panInterfaces := convertPathInterfaceSlice(p.Metadata().Interfaces)
	if len(panInterfaces) == 0 {
		return ""
	}

	b := &strings.Builder{}
	_, err := fmt.Fprintf(b, "%d", panInterfaces[0].IfID)
	if err != nil {
		return ""
	}

	for _, i := range panInterfaces[1:] {
		_, err = fmt.Fprintf(b, " %d", i.IfID)
		if err != nil {
			return ""
		}
	}
	return PathFingerprint(b.String())
}

func convertPathInterfaceSlice(interfaces []snet.PathInterface) []pan.PathInterface {
	pis := make([]pan.PathInterface, len(interfaces))
	for i, spi := range interfaces {
		pis[i] = pan.PathInterface{
			IA:   pan.IA(spi.IA),
			IfID: pan.IfID(spi.ID),
		}
	}
	return pis
}
