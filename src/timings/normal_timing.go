package timings

import "time"

type TimingResponse struct {
	TotalTime   time.Duration
	AverageTime time.Duration
}
