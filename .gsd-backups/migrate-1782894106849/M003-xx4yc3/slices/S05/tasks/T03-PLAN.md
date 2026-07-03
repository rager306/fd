---
estimated_steps: 1
estimated_files: 3
skills_used: []
---

# T03: Verify and close runtime baseline milestone

Run final verification: compose config, Go tests, GitNexus change detection for repo fd, complete milestone, and commit local changes.

## Inputs

- `S05 T02`

## Expected Output

- `M003 validation`
- `commit`

## Verification

docker compose config && cd api && go test ./... -short && gitnexus_detect_changes(repo=fd).

## Observability Impact

Leaves project with committed evidence and clean state.
