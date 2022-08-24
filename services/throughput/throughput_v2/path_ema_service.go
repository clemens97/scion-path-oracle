package throughput_v2

import (
	"github.com/clemens97/scion-path-oracle"
	"github.com/clemens97/scion-path-oracle/services/throughput"
	"go.uber.org/zap"
	"sync"
)

type PathEMAServiceI interface {
	Update(val float64, fp oracle.PathFingerprint)
	GetEMAOrDefault(fp oracle.PathFingerprint, def float64) float64
}

type PathEMAService struct {
	state     map[oracle.PathFingerprint]float64
	lock      sync.RWMutex
	smoothing float64
	logger    *zap.SugaredLogger
}

func NewPathEMAService(smoothing float64, logger *zap.SugaredLogger) PathEMAService {
	return PathEMAService{smoothing: smoothing,
		state:  make(map[oracle.PathFingerprint]float64),
		logger: logger.With("path_ema_smoothing", smoothing)}
}

func (p *PathEMAService) Update(val float64, fp oracle.PathFingerprint) {
	p.lock.Lock()
	defer p.lock.Unlock()

	if prev, loaded := p.state[fp]; loaded {
		p.state[fp] = throughput.EMA(prev, val, p.smoothing)
		p.logger.Debugw("updated path ema",
			"fingerprint", fp, "value_old", prev, "value_new", p.state[fp])
	} else {
		p.state[fp] = val
		p.logger.Debugw("initialized path ema", "fingerprint", fp, "value", val)
	}
}

func (p *PathEMAService) GetEMAOrDefault(fp oracle.PathFingerprint, def float64) float64 {
	p.lock.RLock()
	defer p.lock.RUnlock()
	if val, loaded := p.state[fp]; loaded {
		return val
	}
	return def
}
