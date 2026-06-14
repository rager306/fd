package buildinfo

import (
	"testing"
	"time"
)

func TestNewAppliesDefaults(t *testing.T) {
	info := New(Info{})

	if info.Service != DefaultService {
		t.Fatalf("Service = %q, want %q", info.Service, DefaultService)
	}
	if info.Version != DefaultVersion {
		t.Fatalf("Version = %q, want %q", info.Version, DefaultVersion)
	}
	if info.BuildHash != DefaultBuildHash {
		t.Fatalf("BuildHash = %q, want %q", info.BuildHash, DefaultBuildHash)
	}
	if info.BuildDate != DefaultBuildDate {
		t.Fatalf("BuildDate = %q, want %q", info.BuildDate, DefaultBuildDate)
	}
	if info.StartedAt.IsZero() {
		t.Fatal("StartedAt should be set when missing")
	}
}

func TestNewPreservesProvidedValues(t *testing.T) {
	startedAt := time.Now().Add(-time.Second)
	info := New(Info{
		Service:      "custom-service",
		Version:      "2.0.0",
		Model:        "deepvk/USER-bge-m3",
		ModelVersion: "2026-06-13",
		BuildHash:    "abc1234",
		BuildDate:    "2026-06-13T00:00:00Z",
		StartedAt:    startedAt,
	})

	if info.Service != "custom-service" {
		t.Fatalf("Service = %q", info.Service)
	}
	if info.Version != "2.0.0" {
		t.Fatalf("Version = %q", info.Version)
	}
	if info.Model != "deepvk/USER-bge-m3" {
		t.Fatalf("Model = %q", info.Model)
	}
	if info.ModelVersion != "2026-06-13" {
		t.Fatalf("ModelVersion = %q", info.ModelVersion)
	}
	if info.BuildHash != "abc1234" {
		t.Fatalf("BuildHash = %q", info.BuildHash)
	}
	if info.BuildDate != "2026-06-13T00:00:00Z" {
		t.Fatalf("BuildDate = %q", info.BuildDate)
	}
	if !info.StartedAt.Equal(startedAt) {
		t.Fatalf("StartedAt = %s, want %s", info.StartedAt, startedAt)
	}
}

func TestUptimeIncreases(t *testing.T) {
	info := New(Info{StartedAt: time.Now().Add(-time.Millisecond)})
	first := info.Uptime()
	time.Sleep(time.Millisecond)
	second := info.Uptime()

	if first <= 0 {
		t.Fatalf("first uptime = %s, want > 0", first)
	}
	if second <= first {
		t.Fatalf("second uptime = %s, want > first %s", second, first)
	}
}

func TestZeroInfoUptimeIsZero(t *testing.T) {
	var info Info
	if got := info.Uptime(); got != 0 {
		t.Fatalf("zero uptime = %s, want 0", got)
	}
}
