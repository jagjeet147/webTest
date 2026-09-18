package engine

import (
	"context"
	"math/rand"

	"trafficlab/internal/httpclient"
	"trafficlab/internal/metrics"
)

func worker(ctx context.Context, jobs <-chan struct{}, client *httpclient.Client, method, target, body string, results *metrics.Metrics, seed int64) {
	random := rand.New(rand.NewSource(seed))
	for {
		select {
		case <-jobs:
			result := client.Do(ctx, method, target, body, random)
			results.Record(result.StatusCode, result.Duration, result.Err)
		case <-ctx.Done():
			return
		}
	}
}
