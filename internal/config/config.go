// Package config loads runtime configuration from the environment.
package config

import (
	"fmt"
	"os"
	"strings"
)

// Environment variables read by New.
const (
	EnvGreeting = "HELLO_FX_GREETING"
	EnvLogLevel = "HELLO_FX_LOG_LEVEL"
)

// Defaults applied when the corresponding environment variable is unset.
const (
	DefaultGreeting = "Hello, %s!"
	DefaultLogLevel = "info"
)

// Config holds runtime configuration for the hello-fx CLI.
type Config struct {
	// Greeting is the greeting template. It must contain exactly one %s
	// verb, which is replaced with the greeted name.
	Greeting string

	// LogLevel is the zap log level name for the application logger.
	LogLevel string
}

// New builds a Config from well-known environment variables, applying
// defaults for anything unset.
func New() (*Config, error) {
	cfg := &Config{
		Greeting: getenv(EnvGreeting, DefaultGreeting),
		LogLevel: getenv(EnvLogLevel, DefaultLogLevel),
	}
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) validate() error {
	if n := strings.Count(c.Greeting, "%s"); n != 1 {
		return fmt.Errorf("config: %s must contain exactly one %%s verb, got %d in %q", EnvGreeting, n, c.Greeting)
	}
	if c.LogLevel == "" {
		return fmt.Errorf("config: %s must not be empty", EnvLogLevel)
	}
	return nil
}

func getenv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
