---
milestone: M043-dpr0cq
slice: S02
title: Tier 2 lint adoption and complexity refactor — phase 2 report
date: 2026-06-13
status: complete
---

# M043 S02 Phase 2 Report: Tier 2 lint adoption and complexity refactor

## Outcome

`golangci-lint run --config .golangci.yml ./...` now exits 0 with **18 enabled linters** in fail mode (12 from S01 + 6 new Tier 2). `go test ./...` exits 0 across all Go packages.

## Linters active (18)

| Tier | Linters | Status |
|---|---|---|
| M040 baseline | `errcheck`, `govet`, `ineffassign`, `staticcheck`, `unused`, `goconst`, `misspell` | fail mode |
| M043 S01 Tier 1 | `gosec`, `bodyclose`, `prealloc`, `errorlint`, `revive` | fail mode |
| M043 S02 Tier 2 | `gocyclo`, `gocritic`, `durationcheck`, `unparam`, `contextcheck`, `nilnil` | fail mode |

## Baseline → Final issue count

| Run | Issues | Breakdown |
|---|---:|---|
| T01 godoc baseline | 44 | `revive:exported` comments missing |
| T01 after godoc pass | 0 | 44 exported-symbol gaps fixed |
| T02 Tier 2 baseline | 17 | 12 `gocritic`, 4 `gocyclo`, 1 `unparam` |
| T02 after refactor | 0 | all Tier 2 issues fixed |
| T03 final fail mode | 0 | 18 linters, no severity overrides |

## T01: godoc pass

`revive:exported` was enabled after adding doc comments for public API surfaces across:

- `api/cache/local.go`: `LocalCache`, constructor, Get/Set/Delete
- `api/cache/redis.go`: `RedisCache`, options/config constructors, HashText/Get/Set/Ping/Close
- `api/embed/onnx_disabled.go`: disabled ONNX placeholder type/constructor/methods
- `api/embed/onnx_manifest.go`: public manifest constants, sentinel errors, manifest/validation types, load/validate functions
- `api/embed/onnx_types.go`: `ONNXEmbedderOptions`
- `api/embed/tei.go`: `TEIClient`, constructor, `Embed`
- `api/embed/types.go`: request/response DTOs and usage types
- `api/handlers/batch.go`: `BatchHandler`, constructor
- `api/handlers/embeddings.go`: `Embedder`, `EmbeddingsHandler`, constructor
- `api/handlers/health.go`: `RuntimeHealth`, health handler constructors

Result: `revive:exported` 44 issues → 0.

## T02: Tier 2 findings and fixes

### `gocritic` (12 issues → 0)

| Finding | Fix |
|---|---|
| `unnamedResult` in `unmarshalEmbedding` | named returns `(embedding []float32, dim int)` |
| `unnamedResult` in `RedisCache.Get` | named returns `(embedding []float32, found bool, err error)` |
| `paramTypeCombine` in ONNX/main/handler tests | combined adjacent same-typed params/returns |
| `httpNoBody` in health/recovery tests | replaced nil request bodies with `http.NoBody` |
| `exitAfterDefer` in `main.go` Redis ping failure | moved Redis close defer after successful ping; explicit close on failure |
| `exitAfterDefer` in ONNX init failure | moved Redis close defer after backend init; explicit close on ONNX init failure |
| `paramTypeCombine` in `backfillMisses` | combined `embs, embeddings [][]float32` |

### `gocyclo` (4 issues → 0)

| Function | Baseline | Resolution |
|---|---:|---|
| `(*ONNXArtifactManifest).ValidateArtifact` | 21 | split into `validateManifestMetadata` and `validateArtifactFile`; preserved sentinel errors and tested diagnostic text |
| `(*EmbeddingsHandler).CreateEmbedding` | 30 | split into request extraction, inline validation, defaults, chunk cache/model processing, response assembly |
| `ValidateEmbeddingsRequest` | 16 | split into content length, bind error handling, payload validation, input lengths, dimensions, encoding_format helpers |
| `main` | 16 | extracted ONNX preflight logging helper; also fixed lifecycle ordering around Redis close |

One test-only complexity suppression remains:

```go
//nolint:gocyclo // table-driven production integration coverage intentionally exercises many request/error/cache paths in one matrix.
func TestCreateEmbedding_ProductionHandler(t *testing.T) { ... }
```

Rationale: this is a dense table-driven integration matrix, not production control flow. Splitting it would reduce signal locality and make regression coverage harder to read.

### `unparam` (1 issue → 0)

`runMiddleware(t, method, body, contentLen)` always received `http.MethodPost`; removed the method parameter and fixed call sites.

## Behavior-preservation notes

- `ValidateArtifact` refactor originally changed one diagnostic string (`dimensions=512`); tests caught it. Restored `expected_dimensions=512` to preserve the existing diagnostic contract.
- `resolveONNXArtifactPath` remains used by `ValidateArtifact`; the helper is important because manifests can be loaded from subdirectories while artifact paths are stored relative to repo/runtime roots.
- `maxBatchSize` remains the single validation source (128) in middleware; post-refactor validation uses the const instead of hard-coded literals.
- `CreateEmbedding` behavior is unchanged: validated request from middleware preferred, inline fallback preserved for tests/alternate mounts, cache hits short-circuit TEI calls, TEI misses are chunked at 32, response shape remains OpenAI-compatible.

## Final verification

```bash
cd /root/fd/api
PATH=$PATH:/root/go/bin go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 \
  run --config ../.golangci.yml ./...
# 0 issues, exit 0

PATH=$PATH:/root/go/bin go test ./...
# ok fd-api, fd-api/cache, fd-api/embed, fd-api/handlers, fd-api/middleware
```

Evidence files:

- `benchmark-results/m043-s02-godoc-baseline.txt`
- `benchmark-results/m043-s02-tier2-baseline.txt`
- `benchmark-results/m043-s02-tier2-after-refactor.txt`
- `benchmark-results/m043-s02-final-lint.txt`
- `benchmark-results/m043-s02-go-test.txt`

## Files changed in S02

| File | Change |
|---|---|
| `.golangci.yml` | +6 Tier 2 linters, `gocyclo.min-complexity=15`, `gocritic` settings, `revive:exported`, all linters fail mode |
| `api/cache/local.go` | godoc comments |
| `api/cache/redis.go` | godoc comments, named returns |
| `api/embed/onnx_disabled.go` | godoc comments |
| `api/embed/onnx_manifest.go` | godoc comments, complexity split, restored `resolveONNXArtifactPath` use |
| `api/embed/onnx_types.go` | godoc comments |
| `api/embed/tei.go` | godoc comments |
| `api/embed/types.go` | godoc comments for DTOs |
| `api/handlers/batch.go` | godoc comments, removed duplicate package comment |
| `api/handlers/embeddings.go` | godoc comments, `CreateEmbedding` complexity split |
| `api/handlers/health.go` | godoc comments |
| `api/middleware/validation.go` | `ValidateEmbeddingsRequest` complexity split, `maxBatchSize` restored as source of truth |
| `api/main.go` | close lifecycle fix for `os.Exit` paths, extracted ONNX preflight logging helper |
| tests | gocritic/unparam fixes (`http.NoBody`, param combine, runMiddleware method removal), one justified test-only `//nolint:gocyclo` |

## Carry-over to S03

- Add `govulncheck` as a separate CI step (not part of golangci-lint).
- Finalize static analysis docs/recommendation to reflect S01+S02 outcomes.
- Keep Tier 3 style/perf linters (`gofumpt`, `structslop`, `dupl`, `wsl`, etc.) out of scope unless explicitly requested.
