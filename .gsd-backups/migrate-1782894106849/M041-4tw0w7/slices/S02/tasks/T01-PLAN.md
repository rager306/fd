---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T01: Added lifecycle State package with warmup/readiness/shutdown flags, in-flight request tracking, drain timeout, last error, and context helpers.

api/lifecycle/state.go: type State struct { warmupDone atomic.Bool; shuttingDown atomic.Bool; inflight sync.WaitGroup; lastError atomic.Value }. Methods: MarkWarmupDone(), IsReady() bool, BeginShutdown(), IsShuttingDown() bool, TrackRequest(start, done), WaitDrain(timeout) error. State синглтон, передаётся в handlers через context.

## Inputs

- None specified.

## Expected Output

- `api/lifecycle/state.go`
- `api/lifecycle/state_test.go`

## Verification

Unit tests: MarkWarmupDone → IsReady true. BeginShutdown → IsShuttingDown true. TrackRequest properly tracks inflight, WaitDrain(0) returns immediately when empty, WaitDrain(100ms) blocks. Test concurrent shutdown while requests inflight.
