package main

import (
	"log/slog"
	"os"
	"strconv"
	"testing"

	"pgregory.net/rapid"
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

func TestGetEnvIntParsesDigits_Rapid(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		value := rapid.IntRange(0, 1_000_000).Draw(t, "value")
		if err := os.Setenv("FD_TEST_ENV_INT", strconv.Itoa(value)); err != nil {
			t.Fatalf("Setenv failed: %v", err)
		}
		if got := getEnvInt("FD_TEST_ENV_INT", -1); got != value {
			t.Fatalf("getEnvInt = %d, want %d", got, value)
		}
	})
}

func TestGetEnvIntFallsBackForInvalidValues(t *testing.T) {
	for _, value := range []string{"", "-1", "12ms", " 12", "12 "} {
		t.Run(value, func(t *testing.T) {
			t.Setenv("FD_TEST_ENV_INT", value)
			if got := getEnvInt("FD_TEST_ENV_INT", 7); got != 7 {
				t.Fatalf("getEnvInt(%q) = %d, want fallback", value, got)
			}
		})
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
