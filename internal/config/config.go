package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	URL         string        `yaml:"url"`
	Method      string        `yaml:"method"`
	RPS         float64       `yaml:"rps"`
	Duration    time.Duration `yaml:"duration"`
	Concurrency int           `yaml:"concurrency"`
	Body        string        `yaml:"body,omitempty"`
}

func (c *Config) UnmarshalYAML(value *yaml.Node) error {
	var raw struct {
		URL         string  `yaml:"url"`
		Method      string  `yaml:"method"`
		RPS         float64 `yaml:"rps"`
		Duration    string  `yaml:"duration"`
		Concurrency int     `yaml:"concurrency"`
		Body        string  `yaml:"body,omitempty"`
	}
	if err := value.Decode(&raw); err != nil {
		return err
	}
	duration, err := time.ParseDuration(raw.Duration)
	if err != nil {
		return fmt.Errorf("parse duration: %w", err)
	}
	*c = Config{URL: raw.URL, Method: raw.Method, RPS: raw.RPS, Duration: duration, Concurrency: raw.Concurrency, Body: raw.Body}
	return nil
}

func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read config: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse config: %w", err)
	}
	return cfg, nil
}

func (c Config) Validate() error {
	parsedURL, err := url.ParseRequestURI(c.URL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return fmt.Errorf("url must be an absolute HTTP(S) URL")
	}
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("url scheme must be http or https")
	}
	if strings.TrimSpace(c.Method) == "" {
		return fmt.Errorf("method is required")
	}
	if c.RPS <= 0 {
		return fmt.Errorf("rps must be greater than zero")
	}
	if c.Duration <= 0 {
		return fmt.Errorf("duration must be greater than zero")
	}
	if c.Concurrency <= 0 {
		return fmt.Errorf("concurrency must be greater than zero")
	}
	return nil
}
