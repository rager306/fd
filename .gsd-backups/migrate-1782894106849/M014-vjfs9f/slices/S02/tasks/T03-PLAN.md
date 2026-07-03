---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T03: Verify TEI benchmark artifact

Parse TEI artifact for required sections, snapshot fields, Redis/cache sections, PASS/summary markers, and raw probe text absence.

## Inputs

- `benchmark-results/fd-benchmark-m014-tei-baseline.txt`

## Expected Output

- `Task summary with artifact verification evidence`

## Verification

Parser/leak checks and GitNexus detect_changes pass.

## Observability Impact

Confirms S03 can compare against a valid TEI artifact.
