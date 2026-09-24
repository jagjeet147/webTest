.PHONY: build build-windows test fmt run

build:
	go build ./cmd/trafficlab

build-windows:
	go build -o trafficlab.exe ./cmd/trafficlab

test:
	go test ./...

fmt:
	gofmt -w cmd internal tests

run:
	go run ./cmd/trafficlab run --config configs/example.yaml
