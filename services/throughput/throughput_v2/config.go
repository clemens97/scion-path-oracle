package throughput_v2

import (
	"os"
	"strconv"
	"strings"
)

const (
	maxReportAmountDefault  = 100
	bottleNeckRatioDefault  = 0.5
	linkEMASmoothingDefault = 0.3
	pathEmaSmoothingDefault = 0.3
)

type config struct {
	// maxReportAmount is the amount of reports we store for a link. We drop the oldest report when we are about to exceed this amount.
	maxReportAmount int
	// bottleNeckRatio is the ratio of paths we consider to be bottlenecked by another link
	bottleNeckRatio float64
	// linkEMASmoothing is the EMA Smoothing to applied to all Bandwidth Reports of Paths who are not considered a bottleneck
	linkEMASmoothing float64
	pathEMASmoothing float64
}

func newConfig() *config {
	envBase := strings.ToUpper(name) + "_"
	maxRepAmount, err := strconv.ParseInt(os.Getenv(envBase+"MAX_REP"), 10, 32)
	if err != nil {
		maxRepAmount = maxReportAmountDefault
	}
	bottleRatio, err := strconv.ParseFloat(os.Getenv(envBase+"BOT_RAT"), 64)
	if err != nil {
		bottleRatio = bottleNeckRatioDefault
	}
	linkSmoothing, err := strconv.ParseFloat(os.Getenv(envBase+"L_SMOOTHING"), 64)
	if err != nil {
		linkSmoothing = linkEMASmoothingDefault
	}
	pathSmoothing, err := strconv.ParseFloat(os.Getenv(envBase+"P_SMOOTHING"), 64)
	if err != nil {
		pathSmoothing = pathEmaSmoothingDefault
	}

	return &config{
		maxReportAmount:  int(maxRepAmount),
		bottleNeckRatio:  bottleRatio,
		linkEMASmoothing: linkSmoothing,
		pathEMASmoothing: pathSmoothing,
	}
}
