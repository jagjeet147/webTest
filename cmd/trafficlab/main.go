package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"trafficlab/internal/config"
	"trafficlab/internal/controller"
)

const version = "0.1"

func main() {
	if len(os.Args) < 2 || os.Args[1] == "--help" || os.Args[1] == "-h" {
		printUsage()
		return
	}

	if os.Args[1] != "run" {
		fmt.Fprintf(os.Stderr, "unknown command %q\n\n", os.Args[1])
		printUsage()
		os.Exit(2)
	}

	if err := run(os.Args[2:]); err != nil {
		fmt.Fprintf(os.Stderr, "trafficlab: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	flags := flag.NewFlagSet("run", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)
	url := flags.String("url", "", "target URL")
	method := flags.String("method", "GET", "HTTP method")
	rps := flags.Float64("rps", 1, "target requests per second")
	duration := flags.Duration("duration", 30*time.Second, "test duration")
	concurrency := flags.Int("concurrency", 1, "maximum concurrent workers")
	body := flags.String("body", "", "request body")
	configPath := flags.String("config", "", "optional YAML configuration file")

	if err := flags.Parse(args); err != nil {
		return err
	}

	settings := config.Config{
		URL:         *url,
		Method:      *method,
		RPS:         *rps,
		Duration:    *duration,
		Concurrency: *concurrency,
		Body:        *body,
	}
	if *configPath != "" {
		loaded, err := config.Load(*configPath)
		if err != nil {
			return err
		}
		settings = loaded
	}
	if err := settings.Validate(); err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	result, err := controller.Run(ctx, settings)
	if err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	fmt.Print(result.Format(settings))
	return nil
}

func printUsage() {
	fmt.Printf("TrafficLab v%s\n\n", version)
	fmt.Println("Usage:")
	fmt.Println("  trafficlab run --url URL [flags]")
	fmt.Println("\nRun flags:")
	flags := flag.NewFlagSet("run", flag.ContinueOnError)
	flags.String("url", "", "target URL")
	flags.String("method", "GET", "HTTP method")
	flags.Float64("rps", 1, "target requests per second")
	flags.Duration("duration", 30*time.Second, "test duration")
	flags.Int("concurrency", 1, "maximum concurrent workers")
	flags.String("body", "", "request body")
	flags.String("config", "", "optional YAML configuration file")
	flags.PrintDefaults()
}
