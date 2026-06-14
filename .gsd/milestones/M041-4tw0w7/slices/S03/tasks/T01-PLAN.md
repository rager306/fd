---
estimated_steps: 1
estimated_files: 3
skills_used: []
---

# T01: Build/version metadata injection

api/buildinfo package: type Info { Service, Version, Model, ModelVersion, BuildHash, BuildDate, StartedAt, Uptime() time.Duration }. Значения передаются через ldflags при сборке (-X main.Version=2.0.0 -X main.BuildHash=$(git rev-parse --short HEAD) -X main.BuildDate=2026-06-13T00:00:00Z). Default values если ldflags не заданы. Обновить Dockerfile (если нужно) и Makefile (если есть) для передачи ldflags.

## Inputs

- None specified.

## Expected Output

- `api/buildinfo/info.go`
- `api/buildinfo/info_test.go`
- `Dockerfile`

## Verification

go test ./api/buildinfo/...: Uptime корректно увеличивается. Build с -ldflags передаёт значения в бинарь.
