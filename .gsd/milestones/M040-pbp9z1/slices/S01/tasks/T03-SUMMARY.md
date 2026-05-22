---
id: T03
parent: S01
milestone: M040-pbp9z1
key_files:
  - api/handlers/embeddings_integration_test.go
  - docs/same-host-embedding-service-contract.md
  - .gsd/DECISIONS.md
key_decisions:
  - D041: Treat `/v1/embeddings` request `model` as OpenAI-compatibility metadata and keep response model plus `/health.runtime.model` authoritative.
duration: 
verification_result: passed
completed_at: 2026-05-22T05:32:34.159Z
blocker_discovered: false
---

# T03: Resolved `/v1/embeddings` request `model` as compatibility metadata while making response model and `/health.runtime.model` authoritative.

**Resolved `/v1/embeddings` request `model` as compatibility metadata while making response model and `/health.runtime.model` authoritative.**

## What Happened

Verified the task plan and surrounding embeddings handler, integration tests, service contract, and S01 research. Existing tests already exercised requests whose `model` differed from the configured handler model, so hard-rejecting mismatches would have broadened behavior change and risked breaking OpenAI-compatible clients that send placeholders. I chose the smallest safe contract: leave handler behavior unchanged, add an explicit regression test proving mismatched request model values still generate embeddings, and tighten the same-host contract language to state that request `model` is not used for routing or validation.

## Verification

Formatted the updated Go test file and ran the planned short Go suite with `cd api && go test ./... -short`; all packages passed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `gofmt -w api/handlers/embeddings_integration_test.go && cd api && go test ./... -short` | 0 | ✅ pass | 1763ms |

## Deviations

Did not modify `api/handlers/embeddings.go` because the selected safest behavior is the existing implementation behavior; T03 was resolved through an explicit regression test and contract documentation.

## Known Issues

None.

## Files Created/Modified

- `api/handlers/embeddings_integration_test.go`
- `docs/same-host-embedding-service-contract.md`
- `.gsd/DECISIONS.md`
