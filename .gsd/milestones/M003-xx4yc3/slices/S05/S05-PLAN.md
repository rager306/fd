# S05: Runtime fixes and optimization plan

**Goal:** Fix discovered runtime defects if any and produce evidence-backed optimization plan.
**Demo:** Known errors are fixed or performance improvements are prioritized from evidence.

## Must-Haves

- Any runtime defect has root cause, impact analysis, minimal fix, tests, and rerun gate.
- If no defect fix is needed, explicitly document no-code-change decision.
- Performance improvement proposals include evidence, expected measurable outcome, files, validation command, and rollback path.

## Proof Level

- This slice proves: rerun failed gates or written assessment

## Integration Closure

Closes validation loop and defines next performance work from measurements.

## Verification

- Records root causes, fixes, and benchmark-backed recommendations.

## Tasks

- [x] **T01: Classify findings and fixes** `est:small`
  Classify runtime findings from S01-S04 into fixed defects, non-blocking host/runtime notes, and future optimizations.
  - Files: `docker-compose.override.yaml`, `benchmark-results/fd-benchmark-baseline-py313.txt`, `benchmark-results/fd-runtime-stats-logs.txt`
  - Verify: Written classification in task summary.

- [x] **T02: Write optimization assessment** `est:medium`
  Write an evidence-backed performance optimization assessment with prioritized improvements, expected outcomes, touched files, validation commands, and rollback paths.
  - Files: `benchmark-results/fd-benchmark-baseline-py313.txt`, `benchmark-results/fd-runtime-stats-logs.txt`
  - Verify: Assessment artifact saved and references benchmark evidence.

- [x] **T03: Verify and close runtime baseline milestone** `est:small`
  Run final verification: compose config, Go tests, GitNexus change detection for repo fd, complete milestone, and commit local changes.
  - Files: `docker-compose.override.yaml`, `benchmark-results/`, ` .gsd/milestones/M003-xx4yc3`
  - Verify: docker compose config && cd api && go test ./... -short && gitnexus_detect_changes(repo=fd).

## Files Likely Touched

- docker-compose.override.yaml
- benchmark-results/fd-benchmark-baseline-py313.txt
- benchmark-results/fd-runtime-stats-logs.txt
- benchmark-results/
-  .gsd/milestones/M003-xx4yc3
