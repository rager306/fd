---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T02: Wire FD_ASYNC_CHUNKS env в handler + main.go

api/handlers/embeddings.go: в CreateEmbedding, если FD_ASYNC_CHUNKS=true → использовать async orchestrator (chunks of 32 sent in parallel), иначе current sequential loop. api/main.go: parse FD_ASYNC_CHUNKS env на startup, log mode. Sync mode default (off) для backward compat.

## Inputs

- None specified.

## Expected Output

- `api/handlers/embeddings.go (modified)`
- `api/main.go (env read)`

## Verification

FD_ASYNC_CHUNKS=true → 4 parallel TEI calls per request (verify via TEI logs overlapping timestamps). FD_ASYNC_CHUNKS=false (default) → sequential (no regression vs M041). Integration test asserts both modes pass.
