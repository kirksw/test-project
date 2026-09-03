package config

import (
	"testing"
)

func TestNewDefaults(t *testing.T) {
	t.Setenv(EnvGreeting, "")
	t.Setenv(EnvLogLevel, "")

	cfg, err := New()
	if err != nil {
		t.Fatalf("New() error = %v, want nil", err)
	}
	if cfg.Greeting != DefaultGreeting {
		t.Errorf("Greeting = %q, want %q", cfg.Greeting, DefaultGreeting)
	}
	if cfg.LogLevel != DefaultLogLevel {
		t.Errorf("LogLevel = %q, want %q", cfg.LogLevel, DefaultLogLevel)
	}
}

func TestNewFromEnv(t *testing.T) {
	t.Setenv(EnvGreeting, "Hi %s, welcome!")
	t.Setenv(EnvLogLevel, "debug")

	cfg, err := New()
	if err != nil {
		t.Fatalf("New() error = %v, want nil", err)
	}
	if cfg.Greeting != "Hi %s, welcome!" {
		t.Errorf("Greeting = %q, want %q", cfg.Greeting, "Hi %s, welcome!")
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("LogLevel = %q, want %q", cfg.LogLevel, "debug")
	}
}

func TestNewInvalidGreeting(t *testing.T) {
	tests := []struct {
		name     string
		greeting string
	}{
		{"no verb", "Hello!"},
		{"two verbs", "Hello %s and %s!"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(EnvGreeting, tt.greeting)

			if _, err := New(); err == nil {
				t.Errorf("New() with greeting %q succeeded, want error", tt.greeting)
			}
		})
	}
}
