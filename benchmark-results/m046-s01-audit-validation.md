# M046 S01 Audit Validation Evidence

Captured: 2026-06-14
Source: GitHub issue #3 full-codebase audit.

## Goal

Validate whether issue #3 P0/P1 claims exist in current code before making source changes.

## Static Probe

Command evidence: `gsd_exec:cafb3f84-b852-4d21-b71f-c13c9f5afd77`.

Result:

```text
v1_batch_route_only_lifecycle: CONFIRMED
legacy_batch_bare_handler: CONFIRMED
auth_empty_key_allows_all: CONFIRMED
metrics_public: CONFIRMED
localcache_no_close: CONFIRMED
localcache_size_counter: CONFIRMED
legacy_batch_per_input_embed: CONFIRMED
v1batch_per_input_embed: CONFIRMED
```

## Confirmed P0 Signals

### P0 #1: default-open auth posture

- `api/middleware/auth.go`: empty `FD_API_KEY` disables auth for all endpoints.
- `docker-compose.yaml`: API publishes `8000:8000`.
- PR #2 correctly keeps health/readiness public, but does not address public inference/admin/diagnostics exposure when no key is configured.

Classification: confirmed policy risk. Fix in S04.

### P0 #2: `/embeddings/batch` guardrail bypass

- `api/main.go`: `/embeddings/batch` is mounted directly to `batchHandler.CreateBatchEmbeddings`.
- `api/handlers/batch.go`: handler has partial inline validation but no shared body/input/capacity/rate-limit policy.
- Per-input loop can create backend work per input.

Classification: confirmed defect. Fix in S02; N+1 optimization in S03.

### P0 #3: `/v1/batch` guardrail bypass

- `api/main.go`: `/v1/batch` has lifecycle capacity gate only.
- `api/handlers/v1batch.go`: has inner batch count checks, but route lacks request body cap and rate limit.

Classification: confirmed defect. Fix in S02; N+1 optimization in S03.

## Confirmed P1 Signals

- P1 #4 and #5: batch handlers call TEI once per cache miss input. Confirmed in `batch.go` and `v1batch.go`. Fix in S03.
- P1 #6: `/v1/embeddings` sequential cache peek is bounded and lower priority. Defer to S06 or future optimization unless S03 introduces a shared batch cache API.
- P1 #7: `/metrics` auth exemption confirmed. Fix policy in S04.
- P1 #8: rate-limit proxy/map risks require exposure-policy decision. Handle in S04 or S06 depending on chosen edge posture.
- P1 #9: 405 registry/WriteError cleanup needs exact contract check. Handle in S06.
- P1 #10: LocalCache lacks Close and has split `sync.Map` plus manual size accounting. Fix in S05.

## Safety Notes

No destructive load test was run. Abuse paths were validated with code evidence and static probes. S02 should add failing unit/handler tests that prove rejection before backend work, then implement the minimal guardrails.
