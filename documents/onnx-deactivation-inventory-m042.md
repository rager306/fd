---
milestone: M042-fjf2en
slice: S02
task: T01
captured: 2026-06-14T10:25:00Z
---

# M042 S02 T01 — ONNX deactivation inventory

## Purpose

The user directed fd to stop carrying ONNX as an active current runtime path and to focus the project on the working TEI path. ONNX may remain as future research history, but it should not be presented as a current operator option or implementation target.

This inventory defines the removal boundary for S02.

## Current state summary

- Default Docker/API path is TEI.
- Normal `api/go.mod` does not directly require `onnxruntime_go`, but active source still includes ONNX build-tag files, runtime config branches, tests, docs, tools, and `Dockerfile.onnx`.
- `api/embed/onnx_disabled.go` already fails closed when built without the `onnx` tag, but active startup config still accepts ONNX environment variables and tests them.
- Historical artifacts show ONNX was locally promising, but not operationally production-ready.
- S01 RCA found TEI startup currently wastes time on missing ONNX/ORT fallback behavior inside TEI itself; fd should not add additional ONNX ambiguity on top.

## Active surfaces and decisions

| Surface | Examples | Current role | Decision | Rationale |
|---|---|---|---|---|
| Runtime startup config | `api/main.go`, `api/main_test.go`, `EMBEDDING_BACKEND=onnx`, `ONNX_*` env parsing | Active app startup path accepts ONNX config | Remove/disable active ONNX selector | Operators should see TEI as the only current backend. ONNX env should not be a supported runtime path. |
| ONNX embedder implementation | `api/embed/onnx.go`, `api/embed/onnx_disabled.go`, `api/embed/onnx_types.go`, tokenizer files | Build-tagged implementation and disabled placeholder | Remove from active module if feasible; otherwise leave only a small fail-closed tombstone | The current milestone should not maintain runnable ONNX. A tombstone is acceptable only to produce a clear error for stale callers/tests. |
| ONNX manifest validation | `api/embed/onnx_manifest.go`, `api/embed/onnx_manifest_test.go`, docs manifests | Active package code and tests for ONNX artifact contract | Remove from active product code or quarantine as research docs | Artifact contract is not current product scope. Historical manifests may remain under docs. |
| ONNX tests | `api/embed/onnx_test.go`, manifest tests, main ONNX config test | Active test surface | Remove or rewrite to assert ONNX disabled | Current gates should not exercise ONNX runtime readiness. |
| Native tokenizer build tags | `hf_tokenizer_native.go`, `onnx_tokenizer_hf.go`, `hf_tokenizer_native_test.go` | Used by ONNX/tokenizer parity work | Remove/quarantine if only used by ONNX | Avoid native tokenizer artifact/dependency complexity in current TEI-only product. |
| Docker ONNX image | `Dockerfile.onnx`, `.github/workflows/onnx-packaging.yml` | Packaging path for ONNX artifacts | Remove from active CI/build or mark historical if retained | User wants ONNX out of current project path. Active CI should not suggest ONNX is supported. |
| Tools | `tools/*onnx*`, `tools/export_user_bge_m3_dense_onnx.py`, `tools/verify_onnx_*` | Research/provisioning helpers | Keep only if clearly under historical/research docs; otherwise remove from active tools | Tools are useful history but should not be operator workflow. |
| Docs/operator contract | README, `docs/same-host-embedding-service-contract.md`, `docs/fd-v2.md`, compose comments | Tell operators ONNX is opt-in/future optimization | Update to TEI-only current posture | Operator docs must not advertise current ONNX runtime. |
| Compose | `docker-compose.yaml` comments about fp16/ONNX | Confusing current runtime hint | Remove ONNX optimization comments | Compose should describe current TEI runtime only. |
| Historical benchmark artifacts | `benchmark-results/*onnx*`, M010-M019/M040 GSD artifacts | Historical evidence | Keep | These are audit/history, not current runtime. |
| Requirements/decisions | R022, R027, D047 | GSD contract | Keep updated | R022 deferred; R027 active until TEI-only deactivation is validated. |

## Proposed implementation boundary

S02 should not try to delete every historical mention of ONNX. It should remove or neutralize active surfaces:

1. **Startup/config:** remove ONNX as accepted runtime backend from `api/main.go`; TEI is the only current backend.
2. **Tests:** remove ONNX-positive config tests or replace with fail-closed tests.
3. **Build/dependencies:** delete or quarantine build-tagged ONNX implementation if this does not destabilize current TEI build. Run `go test ./...` and `go mod tidy` afterward.
4. **Docs/CI/tools:** remove current operator instructions and active packaging workflow references; keep historical benchmark artifacts.
5. **Runtime health:** ensure health/info still accurately report TEI and do not imply ONNX support.

## Risks

- Removing too much may erase useful research history. Mitigation: keep benchmark/GSD/docs artifacts as historical evidence.
- Removing build-tagged files may require careful go.mod/go.sum cleanup. Mitigation: verify with `go test ./...`, `go list -deps ./...`, lint, and govulncheck.
- Some tokenization helpers might be shared with tests. Mitigation: inspect compile failures and remove only ONNX-specific seams.
- The normal binary may already exclude ONNX due build tags; the main benefit is reducing product complexity and dependency/tooling surface, not necessarily a dramatic binary-size change.

## T01 recommendation

Proceed with TEI-only cleanup in T02-T04. Treat ONNX as historical/future research, not an active runtime path. Do not modify historical benchmark artifacts.
