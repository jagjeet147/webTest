package browser

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"trafficlab/internal/config"
	"trafficlab/internal/identity"
	"trafficlab/internal/metrics"

	"github.com/mxschmitt/playwright-go"
)

func Run(ctx context.Context, cfg config.Config) (metrics.Snapshot, error) {
	proxies, err := loadProxies(cfg.ProxyFile)
	if err != nil {
		return metrics.Snapshot{}, err
	}

	playwrightInstance, err := playwright.Run()
	if err != nil {
		return metrics.Snapshot{}, fmt.Errorf("start playwright: %w (run 'playwright install chromium' first)", err)
	}
	defer playwrightInstance.Stop()

	headless := cfg.Headless
	browserInstance, err := playwrightInstance.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless})
	if err != nil {
		return metrics.Snapshot{}, fmt.Errorf("launch chromium: %w", err)
	}
	defer browserInstance.Close()

	jobs := make(chan struct{})
	results := &metrics.Metrics{}
	workerCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	var workers sync.WaitGroup
	for index := 0; index < cfg.Concurrency; index++ {
		workers.Add(1)
		go func(index int) {
			defer workers.Done()
			runWorker(workerCtx, jobs, browserInstance, proxies, index, cfg, results)
		}(index)
	}

	schedule(workerCtx, cfg.RPS, cfg.Duration, jobs)
	cancel()
	workers.Wait()
	return results.Snapshot(), nil
}

func runWorker(ctx context.Context, jobs <-chan struct{}, browserInstance playwright.Browser, proxies []string, index int, cfg config.Config, results *metrics.Metrics) {
	random := rand.New(rand.NewSource(time.Now().UnixNano() + int64(index)))
	contextOptions := playwright.BrowserNewContextOptions{UserAgent: playwright.String(identity.RandomUserAgent(random))}
	if len(proxies) > 0 {
		contextOptions.Proxy = &playwright.Proxy{Server: proxies[index%len(proxies)]}
	}

	browserContext, err := browserInstance.NewContext(contextOptions)
	if err != nil {
		return
	}
	defer browserContext.Close()

	for {
		select {
		case <-jobs:
			page, err := browserContext.NewPage()
			if err != nil {
				results.Record(0, 0, err)
				continue
			}
			start := time.Now()
			response, navigationErr := page.Goto(cfg.URL)
			if navigationErr != nil {
				results.Record(0, time.Since(start), navigationErr)
			} else {
				for cycle := 0; cycle < cfg.Scrolls; cycle++ {
					if err := page.Mouse().Wheel(0, 800); err != nil {
						break
					}
					if err := page.Mouse().Wheel(0, -800); err != nil {
						break
					}
				}
				results.Record(response.Status(), time.Since(start), nil)
			}
			_ = page.Close()
		case <-ctx.Done():
			return
		}
	}
}

func loadProxies(path string) ([]string, error) {
	if strings.TrimSpace(path) == "" {
		return nil, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open proxy file: %w", err)
	}
	defer file.Close()

	var proxies []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		proxy := strings.TrimSpace(scanner.Text())
		if proxy != "" && !strings.HasPrefix(proxy, "#") {
			proxies = append(proxies, proxy)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read proxy file: %w", err)
	}
	return proxies, nil
}

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
