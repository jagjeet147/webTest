# Architecture

The `trafficlab run` command parses CLI flags or a YAML file into `config.Config`. The controller validates the configuration, creates the HTTP client, and starts the engine.

The engine starts a fixed worker pool and a scheduler. The scheduler emits one job per target request interval until the configured duration or context cancellation. Workers execute requests through the HTTP client, which selects a User-Agent and measures request duration. Metrics records counts and latency observations and reports P50, P95, and P99 values at the end.

The engine deliberately owns cancellation and waits for all workers before returning, so a result is complete and no worker is left behind after a run.
