package controller

import (
	"context"
	"fmt"
	"time"

	"trafficlab/internal/config"
	"trafficlab/internal/engine"
	"trafficlab/internal/httpclient"
	"trafficlab/internal/metrics"
)

type Result struct {
	Metrics metrics.Snapshot
	Elapsed time.Duration
}

func Run(ctx context.Context, cfg config.Config) (Result, error) {
	if err := cfg.Validate(); err != nil {
		return Result{}, fmt.Errorf("validate config: %w", err)
	}
	start := time.Now()
	client := httpclient.New(30 * time.Second)
	snapshot := engine.New(client).Run(ctx, cfg)
	return Result{Metrics: snapshot, Elapsed: time.Since(start)}, nil
}

func (r Result) Format(cfg config.Config) string {
	actualRPS := 0.0
	if r.Elapsed > 0 {
		actualRPS = float64(r.Metrics.Requests) / r.Elapsed.Seconds()
	}
	return fmt.Sprintf(`TrafficLab v0.1

Target       %s
Method       %s
Duration     %s
Target RPS   %.2f
Concurrency  %d

Requests     %d
Success      %d
Errors       %d
Error rate   %.2f%%

Latency
  P50        %.2f ms
  P95        %.2f ms
  P99        %.2f ms

Actual RPS   %.2f
`, cfg.URL, cfg.Method, cfg.Duration, cfg.RPS, cfg.Concurrency, r.Metrics.Requests, r.Metrics.Success, r.Metrics.Errors, r.Metrics.ErrorRate(), r.Metrics.P50, r.Metrics.P95, r.Metrics.P99, actualRPS)
}
