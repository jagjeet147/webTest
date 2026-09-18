package engine

import (
	"context"
	"time"
)

func schedule(ctx context.Context, rps float64, duration time.Duration, jobs chan<- struct{}) {
	interval := time.Duration(float64(time.Second) / rps)
	if interval < time.Nanosecond {
		interval = time.Nanosecond
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	timer := time.NewTimer(duration)
	defer timer.Stop()

	for {
		select {
		case <-ticker.C:
			select {
			case jobs <- struct{}{}:
			case <-ctx.Done():
				return
			}
		case <-timer.C:
			return
		case <-ctx.Done():
			return
		}
	}
}
