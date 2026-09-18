package engine

import (
	"context"
	"sync"
	"time"

	"trafficlab/internal/config"
	"trafficlab/internal/httpclient"
	"trafficlab/internal/metrics"
)

type Engine struct {
	client *httpclient.Client
}

func New(client *httpclient.Client) *Engine {
	return &Engine{client: client}
}

func (e *Engine) Run(ctx context.Context, cfg config.Config) metrics.Snapshot {
	jobs := make(chan struct{})
	results := &metrics.Metrics{}
	workerCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	var workers sync.WaitGroup
	for index := 0; index < cfg.Concurrency; index++ {
		workers.Add(1)
		go func(index int) {
			defer workers.Done()
			worker(workerCtx, jobs, e.client, cfg.Method, cfg.URL, cfg.Body, results, time.Now().UnixNano()+int64(index))
		}(index)
	}

	schedule(workerCtx, cfg.RPS, cfg.Duration, jobs)
	cancel()
	workers.Wait()
	return results.Snapshot()
}
