# TrafficLab

TrafficLab is a Go load-testing tool with two traffic modes: fast HTTP requests and browser sessions that render pages, use isolated user agents and optional proxies, and scroll through the page.

## Quick start

```sh
go run ./cmd/trafficlab run --url http://127.0.0.1:8080/api/test --method GET --rps 100 --duration 30s --concurrency 20
```

Configuration can also be loaded from YAML:

```sh
go run ./cmd/trafficlab run --config configs/example.yaml
```

## Browser mode

### Predefined Windows executable

Build the Windows executable:

```powershell
go build -o trafficlab.exe ./cmd/trafficlab
```

Running `trafficlab.exe` without arguments, or running `trafficlab.exe preset`, starts the predefined browser test:

```text
URL:         https://codespira.com
RPS:         114
Duration:    5 minutes
Concurrency: 40 browser workers
Scrolls:     2 down-and-up cycles
Headless:    true
```

The equivalent command is:

```powershell
trafficlab.exe run --mode browser --url https://codespira.com --rps 114 --duration 5m --concurrency 40 --scrolls 2 --headless=true
```

Press `Ctrl+C` to stop the run early. The executable still requires the Playwright driver and browser runtime on the Windows machine.

Browser mode uses Playwright and Chromium. Install the matching Playwright driver and browser once:

```sh
go run github.com/mxschmitt/playwright-go/cmd/playwright@v0.6201.1 install chromium
```

Run isolated browser contexts with three down-and-up scroll cycles:

```sh
go run ./cmd/trafficlab run --mode browser --url https://example.com --rps 2 --duration 30s --concurrency 20 --scrolls 3
```

Use one proxy URL per line to assign proxy-backed network identities to browser contexts. Blank lines and lines beginning with `#` are ignored:

```sh
go run ./cmd/trafficlab run --mode browser --url https://example.com --concurrency 50 --proxy-file proxies.txt --headless=true
```

The proxy list is required for distinct public source IPs. User-Agent changes alone do not change the source IP address.

Browser mode uses isolated Playwright contexts in one Chromium process. This gives each worker a separate browser profile while using substantially less memory than starting 50 full browser processes.

Network, MAC-address control, scenario flows, sessions, and ramp profiles are not included yet.

## Development

```sh
go test ./...
go build ./cmd/trafficlab
```
