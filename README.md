# TrafficLab

TrafficLab is a small Go load-testing foundation: CLI -> controller -> scheduler -> HTTP workers -> metrics.

## Quick start

```sh
go run ./cmd/trafficlab run --url http://127.0.0.1:8080/api/test --method GET --rps 100 --duration 30s --concurrency 20
```

Configuration can also be loaded from YAML:

```sh
go run ./cmd/trafficlab run --config configs/example.yaml
```

This first version intentionally focuses on HTTP load generation. Network, IP, MAC, scenario, session, and ramp-profile functionality are not included yet.

## Development

```sh
go test ./...
go build ./cmd/trafficlab
```
