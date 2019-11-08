package time2

import (
	"time"

	"github.com/liov/hoper/go/v2/utils/log"
)

func Format(t time.Time) string {
	return t.Format("2006-01-02 15:04:05.999")
}

func TimeCost(start time.Time) {
	log.Info(time.Since(start))
}