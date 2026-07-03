---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T04: Link contract and verify S01 safety

Link the new contract from README. Run Go tests if code changed, run a docs/artifact leak check, and record verification evidence in the task summary. Ensure S01 remains scoped to contract/readiness metadata and does not perform S02 benchmark or S04 runtime recommendation work.

## Inputs

- `README.md`
- `docs/same-host-embedding-service-contract.md`
- `.gsd/milestones/M040-pbp9z1/M040-pbp9z1-CONTEXT.md`

## Expected Output

- `README.md`

## Verification

cd api && go test ./... -short
rg -n "signed|token=|X-Amz|BEGIN|PRIVATE|юридическая справка" docs README.md benchmark-results .gsd/milestones/M040-pbp9z1/slices/S01

## Observability Impact

Leaves a discoverable contract entry point and verifies docs/artifacts do not leak secrets or prohibited raw evidence.
