package throughput

import (
	"errors"
	"github.com/clemens97/scion-path-oracle"
	"strconv"
)

func ParseThroughput(stats oracle.MonitoredProperties) (float64, error) {
	bw, ok := stats["throughput"]
	if !ok {
		return 0, errors.New("throughput not in report")
	}

	switch v := bw.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case string:
		bwf, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, errors.New("could not convert throughput")
		}
		// TODO: handle other representation, e.g. MB, GB ...
		return bwf, nil
	}
	return 0, errors.New("could not convert throughput")
}

func Min(vals []uint64) uint64 {
	if len(vals) == 0 {
		return 0
	}

	min := vals[0]
	for _, v := range vals {
		if v < min {
			min = v
		}
	}
	return min
}

func EMA(previousEMA, newVal float64, smoothingFactor float64) float64 {
	return newVal*smoothingFactor + (1-smoothingFactor)*previousEMA
}

func SliceEMA(vals []float64, smoothingFactor float64) float64 {
	if len(vals) == 0 {
		return 0
	}

	ema := vals[0]
	for _, v := range vals {
		ema = EMA(ema, v, smoothingFactor)
	}

	return ema
}
