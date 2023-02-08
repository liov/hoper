package rate

import (
	"math/rand"
	"time"
)

type SpeedLimiter struct {
	*time.Timer
	randSpeedLimitBase, randSpeedLimitRange time.Duration
}

func NewSpeedLimiter(interval time.Duration) *SpeedLimiter {
	return &SpeedLimiter{
		Timer:              time.NewTimer(interval),
		randSpeedLimitBase: interval,
	}
}

// start:最小等待时间
// stop：最大等待时间
// stop-start: 等待范围
func NewRandSpeedLimiter(start, stop time.Duration) *SpeedLimiter {
	return &SpeedLimiter{
		Timer:               time.NewTimer(time.Duration(rand.Intn(int(start))) + stop - start),
		randSpeedLimitBase:  start,
		randSpeedLimitRange: stop - start,
	}
}

func (s *SpeedLimiter) Reset() {
	if s.randSpeedLimitRange == 0 {
		s.Timer.Reset(s.randSpeedLimitBase)
	} else {
		s.Timer.Reset(time.Duration(rand.Intn(int(s.randSpeedLimitBase))) + s.randSpeedLimitRange)
	}
}
