.PHONY: build test fmt run

build:
	go build ./cmd/trafficlab

test:
	go test ./...

fmt:
	gofmt -w cmd internal tests

run:
	go run ./cmd/trafficlab run --config configs/example.yaml
