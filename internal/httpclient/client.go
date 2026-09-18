package httpclient

import (
	"bytes"
	"context"
	"io"
	"math/rand"
	"net/http"
	"time"

	"trafficlab/internal/identity"
)

type Client struct {
	httpClient *http.Client
}

func New(timeout time.Duration) *Client {
	return &Client{httpClient: &http.Client{Timeout: timeout}}
}

type Result struct {
	StatusCode int
	Duration   time.Duration
	Err        error
}

func (c *Client) Do(ctx context.Context, method, target, body string, random *rand.Rand) Result {
	start := time.Now()
	request, err := http.NewRequestWithContext(ctx, method, target, bytes.NewBufferString(body))
	if err != nil {
		return Result{Duration: time.Since(start), Err: err}
	}
	request.Header.Set("User-Agent", identity.RandomUserAgent(random))
	if body != "" {
		request.Header.Set("Content-Type", "application/json")
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return Result{Duration: time.Since(start), Err: err}
	}
	_, _ = io.Copy(io.Discard, response.Body)
	_ = response.Body.Close()
	return Result{StatusCode: response.StatusCode, Duration: time.Since(start)}
}
