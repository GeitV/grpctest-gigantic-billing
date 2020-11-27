package timings

import "time"

type StreamTimingResponse struct {
	TotalRequestTime    time.Duration
	AverageRequestTime  time.Duration
	TotalResponseTime   time.Duration
	AverageResponseTime time.Duration
}
