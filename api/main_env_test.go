package main

import (
	"log/slog"
	"testing"
)

func TestGetEnvReturnsDefaultForUnsetOrEmpty(t *testing.T) {
	t.Setenv("FD_TEST_ENV_VALUE", "")
	if got := getEnv("FD_TEST_ENV_VALUE", "fallback"); got != "fallback" {
		t.Fatalf("getEnv empty = %q, want fallback", got)
	}
	if got := getEnv("FD_TEST_ENV_MISSING", "fallback"); got != "fallback" {
		t.Fatalf("getEnv missing = %q, want fallback", got)
	}
}

func TestGetEnvReturnsConfiguredValue(t *testing.T) {
	t.Setenv("FD_TEST_ENV_VALUE", "configured")
	if got := getEnv("FD_TEST_ENV_VALUE", "fallback"); got != "configured" {
		t.Fatalf("getEnv configured = %q, want configured", got)
	}
}

func TestGetLogLevel(t *testing.T) {
	cases := map[string]slog.Level{
		"debug":   slog.LevelDebug,
		"DEBUG":   slog.LevelDebug,
		"warn":    slog.LevelWarn,
		"warning": slog.LevelWarn,
		"error":   slog.LevelError,
		"info":    slog.LevelInfo,
		"":        slog.LevelInfo,
		"trace":   slog.LevelInfo,
	}
	for input, want := range cases {
		t.Run(input, func(t *testing.T) {
			if got := getLogLevel(input); got != want {
				t.Fatalf("getLogLevel(%q) = %v, want %v", input, got, want)
			}
		})
	}
}
