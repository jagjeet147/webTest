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
	Mode        string        `yaml:"mode"`
	URL         string        `yaml:"url"`
	Method      string        `yaml:"method"`
	RPS         float64       `yaml:"rps"`
	Duration    time.Duration `yaml:"duration"`
	Concurrency int           `yaml:"concurrency"`
	Body        string        `yaml:"body,omitempty"`
	ProxyFile   string        `yaml:"proxy_file,omitempty"`
	Scrolls     int           `yaml:"scrolls,omitempty"`
	Headless    bool          `yaml:"headless"`
}

func (c *Config) UnmarshalYAML(value *yaml.Node) error {
	var raw struct {
		Mode        string  `yaml:"mode"`
		URL         string  `yaml:"url"`
		Method      string  `yaml:"method"`
		RPS         float64 `yaml:"rps"`
		Duration    string  `yaml:"duration"`
		Concurrency int     `yaml:"concurrency"`
		Body        string  `yaml:"body,omitempty"`
		ProxyFile   string  `yaml:"proxy_file,omitempty"`
		Scrolls     int     `yaml:"scrolls,omitempty"`
		Headless    *bool   `yaml:"headless"`
	}
	if err := value.Decode(&raw); err != nil {
		return err
	}
	duration, err := time.ParseDuration(raw.Duration)
	if err != nil {
		return fmt.Errorf("parse duration: %w", err)
	}
	mode := raw.Mode
	if mode == "" {
		mode = "http"
	}
	headless := true
	if raw.Headless != nil {
		headless = *raw.Headless
	}
	*c = Config{Mode: mode, URL: raw.URL, Method: raw.Method, RPS: raw.RPS, Duration: duration, Concurrency: raw.Concurrency, Body: raw.Body, ProxyFile: raw.ProxyFile, Scrolls: raw.Scrolls, Headless: headless}
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
	if c.Mode != "http" && c.Mode != "browser" {
		return fmt.Errorf("mode must be http or browser")
	}
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
	if c.Mode == "browser" && c.Scrolls < 0 {
		return fmt.Errorf("scrolls cannot be negative")
	}
	return nil
}
