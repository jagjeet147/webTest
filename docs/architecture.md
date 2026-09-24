# Architecture

The `trafficlab run` command parses CLI flags or a YAML file into `config.Config`. The controller validates the configuration and selects either the HTTP engine or the browser runner.

In HTTP mode, the engine starts a fixed worker pool and a scheduler. The scheduler emits one job per target request interval until the configured duration or context cancellation. Workers execute requests through the HTTP client, which selects a User-Agent and measures request duration.

In browser mode, the browser runner starts Chromium through Playwright and creates one isolated browser context per worker. Each context receives a User-Agent and, when configured, a proxy selected from the proxy list. A scheduled job opens a page, navigates to the target URL, scrolls down and back up for the configured number of cycles, records the main response and latency, and closes the page.

Metrics records counts and latency observations and reports P50, P95, and P99 values at the end. A proxy is the mechanism that supplies a distinct public source IP; browser contexts and User-Agent values do not change the network source address.

The engine deliberately owns cancellation and waits for all workers before returning, so a result is complete and no worker is left behind after a run.
