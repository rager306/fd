// Package buildinfo models fd build and runtime metadata for observability endpoints.
package buildinfo

import "time"

const (
	// DefaultService is the fd API service name reported by observability endpoints.
	DefaultService = "fd-api"
	// DefaultVersion is used when the binary was built without version ldflags.
	DefaultVersion = "dev"
	// DefaultBuildHash is used when the binary was built without a git hash ldflag.
	DefaultBuildHash = "unknown"
	// DefaultBuildDate is used when the binary was built without a build date ldflag.
	DefaultBuildDate = "unknown"
)

// Info describes immutable build metadata plus process start time.
type Info struct {
	Service      string
	Version      string
	Model        string
	ModelVersion string
	BuildHash    string
	BuildDate    string
	StartedAt    time.Time
}

// New returns build metadata with stable defaults for unset values. A zero
// StartedAt is replaced with the current time so Uptime is meaningful.
func New(info Info) Info {
	if info.Service == "" {
		info.Service = DefaultService
	}
	if info.Version == "" {
		info.Version = DefaultVersion
	}
	if info.BuildHash == "" {
		info.BuildHash = DefaultBuildHash
	}
	if info.BuildDate == "" {
		info.BuildDate = DefaultBuildDate
	}
	if info.StartedAt.IsZero() {
		info.StartedAt = time.Now()
	}
	return info
}

// Uptime returns the elapsed process lifetime represented by Info.
func (i Info) Uptime() time.Duration {
	if i.StartedAt.IsZero() {
		return 0
	}
	return time.Since(i.StartedAt)
}
