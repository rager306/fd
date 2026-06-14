# S02: TEI active-path cleanup and safe mitigation

**Goal:** Make fd TEI-first by removing or disabling active ONNX runtime/build/config/docs surfaces, then decide whether TEI async chunking remains worth implementing from fresh evidence. ONNX stays future research only.
**Demo:** After this, active build/config/docs no longer present ONNX as a current runtime path; TEI remains the only current backend, and any TEI request-shaping mitigation has fresh safe evidence or is explicitly deferred.

## Must-Haves

- Active production/runtime config exposes TEI as the only current backend; ONNX runtime selector/build path is removed or fails closed as research-only.
- Default build/test/Docker path no longer depends on ONNX runtime/tokenizer artifacts.
- README/docs/compose comments no longer tell operators to pursue ONNX as current optimization.
- R027 is validated with code/docs evidence.
- R021 is either validated by a safe TEI mitigation proof or explicitly deferred with rationale.
- `go test ./...`, golangci-lint v2.12.2, and govulncheck pass after final changes.

## Proof Level

- This slice proves: Code, docs, runtime/config verification, and mandatory Go gates.

## Integration Closure

S02 consumes S01 RCA. It leaves fd in a coherent TEI-only current posture and prepares milestone closure without requiring ONNX S03.

## Verification

- Remove ONNX-related ambiguity from runtime health/docs. Preserve TEI runtime metadata and diagnostic evidence paths.

## Tasks

- [x] **T01: Mapped active ONNX source/config/docs/tooling surfaces and defined the TEI-only removal boundary.** `est:45m`
  Map active ONNX references across source, tests, Dockerfile, compose, docs, tools, CI, and requirements. Distinguish historical research artifacts from active runtime/build surfaces. Produce `documents/onnx-deactivation-inventory-m042.md` with a remove/keep table and risk notes. Do not edit code in this task.
  - Files: `documents/onnx-deactivation-inventory-m042.md`
  - Verify: Inventory artifact exists and lists active source/config/docs surfaces plus explicit keep/remove decisions.

- [x] **T02: Removed ONNX as an accepted active runtime backend from fd startup config; startup now supports TEI only.** `est:1h`
  Edit active Go runtime startup/config so fd no longer accepts ONNX as a current backend selector. Remove or neutralize ONNX env parsing/config branches, update tests accordingly, and ensure invalid ONNX env usage fails closed with a clear TEI-only error or is ignored only if documented. Keep TEI behavior unchanged.
  - Files: `api/main.go`, `api/main_test.go`, `api/embed/`
  - Verify: Targeted Go tests for runtime config pass; TEI startup config still passes; ONNX selector is absent or fails closed as TEI-only.

- [x] **T03: Removed active ONNX build-tagged embedder/runtime files, ONNX Dockerfile, and unused ONNX/tokenizer dependencies from the default module.** `est:1.5h`
  Remove ONNX-only build artifacts, Dockerfile paths, module dependencies, and tests that are no longer part of active product scope, or quarantine them in documentation-only research artifacts if deletion is unsafe. Run `go mod tidy` if dependencies are removed. Preserve historical benchmark/docs files.
  - Files: `api/go.mod`, `api/go.sum`, `Dockerfile.onnx`, `api/embed/onnx*.go`, `api/embed/*onnx*_test.go`, `tools/*onnx*`
  - Verify: Default `go test ./...` works without ONNX runtime/toolchain dependencies; `go list -deps ./...` no longer includes ONNX runtime packages unless justified.

- [x] **T04: Updated operator docs and compose/CI surfaces to TEI-only current posture, with ONNX retained only as historical/future research context.** `est:1h`
  Update README, same-host contract, fd-v2 docs, and relevant operations docs to state TEI is the only current runtime. Mark ONNX as historical/future research, not an operator option. Remove outdated compose comments suggesting ONNX export as current optimization. Update R021/R027/R022 statuses if evidence supports it.
  - Files: `README.md`, `docs/same-host-embedding-service-contract.md`, `docs/fd-v2.md`, `docker-compose.yaml`, `.gsd/REQUIREMENTS.md`
  - Verify: Docs contain TEI-only current posture and no active ONNX operator instructions outside historical/future research notes.

- [ ] **T05: Final TEI-only gates and milestone scope decision** `est:1h`
  Run mandatory M043 gates and final TEI-only checks: `go test ./...`, golangci-lint v2.12.2, govulncheck, and a small runtime/config smoke if Docker service is healthy. Record whether R021 async chunking is deferred or implemented separately. Validate R027. Write final S02 evidence artifacts.
  - Files: `benchmark-results/m042-s02-go-test.txt`, `benchmark-results/m042-s02-lint.txt`, `benchmark-results/m042-s02-govulncheck.txt`, `benchmark-results/m042-s02-tei-only-check.txt`
  - Verify: All mandatory gates pass; R027 validated; R021 either validated with evidence or deferred with explicit rationale.

## Files Likely Touched

- documents/onnx-deactivation-inventory-m042.md
- api/main.go
- api/main_test.go
- api/embed/
- api/go.mod
- api/go.sum
- Dockerfile.onnx
- api/embed/onnx*.go
- api/embed/*onnx*_test.go
- tools/*onnx*
- README.md
- docs/same-host-embedding-service-contract.md
- docs/fd-v2.md
- docker-compose.yaml
- .gsd/REQUIREMENTS.md
- benchmark-results/m042-s02-go-test.txt
- benchmark-results/m042-s02-lint.txt
- benchmark-results/m042-s02-govulncheck.txt
- benchmark-results/m042-s02-tei-only-check.txt
