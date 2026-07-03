---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T03: Verify batch hit benchmark evidence

Run full benchmark verification with Docker stack, Go tests/lint, artifact parser for new sections and Redis deltas, and GitNexus detect_changes. Decide whether S05 should proceed or be skipped based on evidence.

## Inputs

- `benchmark.py`

## Expected Output

- `benchmark-results/fd-benchmark-m009-s04.txt`
- `.gsd/milestones/M009-zjrq6j/slices/S04/tasks/T03-SUMMARY.md`

## Verification

`docker compose config`; `cd api && go test ./... -short`; pinned lint; `uv run --python 3.13 --with requests --with redis python benchmark.py`; parser checks; GitNexus detect_changes.

## Observability Impact

Produces S04 evidence for S05 go/no-go decision.
