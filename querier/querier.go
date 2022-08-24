package querier

import (
	"context"
	"fmt"
	"github.com/clemens97/scion-path-oracle"
	"github.com/netsec-ethz/scion-apps/pkg/pan"
	"github.com/scionproto/scion/go/lib/addr"
	"github.com/scionproto/scion/go/lib/daemon"
	"github.com/scionproto/scion/go/lib/snet"
	"os"
	"strings"
	"time"
)

type PathQuerier interface {
	snet.PathQuerier
	GetPath(dst addr.IA, fp oracle.PathFingerprint) (snet.Path, error)
	LocalIA() addr.IA
}

type Querier struct {
	daemon.Querier
}

func New() (PathQuerier, error) {
	daemonAddr := "127.0.0.1:30255"
	if customDaemon := os.Getenv("SCION_DAEMON_ADDRESS"); len(customDaemon) > 0 {
		daemonAddr = customDaemon
	}
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

	sdConn, err := daemon.NewService(daemonAddr).Connect(ctx)
	if err != nil {
		return nil, err
	}
	ia, err := sdConn.LocalIA(context.Background())
	if err != nil {
		return nil, err
	}
	return &Querier{daemon.Querier{Connector: sdConn, IA: ia}}, nil
}

// GetPath returns the path to a destination dst with a given fingerprint fp, when it exists
func (f Querier) GetPath(dst addr.IA, fp oracle.PathFingerprint) (snet.Path, error) {
	// should we cache this on application level? ~local daemon should be caching already
	spaths, err := f.Query(context.Background(), dst)
	if err != nil {
		return nil, err
	}
	for _, sp := range spaths {
		fps := getFingerprintSnet(sp)
		if fps == fp {
			return sp, nil
		}
	}
	return nil, nil
}

func (q Querier) LocalIA() addr.IA {
	return q.IA
}

// same fingerprint as scion-apps/pkg/pan/path.go uses
func getFingerprintSnet(path snet.Path) oracle.PathFingerprint {
	panInterfaces := convertPathInterfaceSlice(path.Metadata().Interfaces)
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
	return oracle.PathFingerprint(b.String())
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
