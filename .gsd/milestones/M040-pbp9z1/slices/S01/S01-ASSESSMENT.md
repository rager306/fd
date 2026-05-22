---
sliceId: S01
uatType: artifact-driven
verdict: PASS
date: 2026-05-22T07:30:00Z
---

# UAT Result — S01

## Checks

| Check | Mode | Result | Notes |
|-------|------|--------|-------|
| Open `README.md` and confirm it links to `docs/same-host-embedding-service-contract.md` as the canonical same-host embedding service contract. | artifact | PASS | `gsd_exec` `8702d1d2-9c57-4444-b6b2-8f646efb8863` found `README.md:59` linking to `docs/same-host-embedding-service-contract.md` and describing it as the full same-host local client service contract. |
| Open `docs/same-host-embedding-service-contract.md` and verify it documents `/health`, `/v1/embeddings`, and `/embeddings/batch` request/response expectations. | artifact | PASS | `gsd_exec` `8702d1d2-9c57-4444-b6b2-8f646efb8863` found `/health`, `POST /v1/embeddings`, and `POST /embeddings/batch` in the canonical contract. |
| Confirm the contract explains dimensions, encoding caveats, status/error behavior, timeout/retry guidance, cache namespace guidance, runtime/env requirements, no-silent-fallback rules, readiness limitations, and non-goals. | artifact | PASS | `gsd_exec` `8702d1d2-9c57-4444-b6b2-8f646efb8863` found coverage for dimensions, encoding caveats, status/error codes, timeout/retry guidance, cache namespace guidance, runtime/env requirements, no-silent-fallback rules, readiness limitations, and non-goals. |
| Confirm `/health` runtime metadata is described as safe operational metadata, not as proof of live inference readiness. | artifact | PASS | The stale TEI/default no-runtime language was removed from `docs/same-host-embedding-service-contract.md`. The contract now states TEI/default `/health` includes a safe `runtime` block with operational identity fields only (`backend`, `model`, `dimensions`, `production_default`, and `cache_namespace`), and still states `/health` is not a live inference probe. `gsd_exec` `8702d1d2-9c57-4444-b6b2-8f646efb8863` verified no stale `TEI/default /health does NOT include a runtime block`, `absence of runtime block indicates TEI`, or `empty means TEI` language remains. |
| Confirm `/v1/embeddings` request `model` is documented as compatibility metadata, while response `model` and `/health.runtime.model` are authoritative. | artifact | PASS | `gsd_exec` `8702d1d2-9c57-4444-b6b2-8f646efb8863` found the contract language that request `model` is OpenAI-compatibility metadata only and response `model` plus `/health.runtime.model` are authoritative. |
| From the repository root, run `cd api && go test ./... -short`. | runtime | PASS | `gsd_exec` `8702d1d2-9c57-4444-b6b2-8f646efb8863` ran the command successfully: `ok fd-api`, `ok fd-api/cache`, `ok fd-api/embed`, and `ok fd-api/handlers`. |
| Run the leak audit command: `rg -n "signed|token=|X-Amz|BEGIN|PRIVATE|юридическая справка" docs README.md benchmark-results .gsd/milestones/M040-pbp9z1/slices/S01` and review any matches as policy/history references rather than current secret or raw corpus leaks. | artifact | PASS | `gsd_exec` `8702d1d2-9c57-4444-b6b2-8f646efb8863` ran the exact broad audit command. Matches were policy/history references in docs, benchmark result statements, and GSD planning/summary/UAT/assessment artifacts; no live secret value or signed URL was observed in the audit output. |
| Run a focused high-risk check for actual signed URL token parameters, private key blocks, and the prohibited raw sample in current public docs/artifacts. | artifact | PASS | `gsd_exec` `8702d1d2-9c57-4444-b6b2-8f646efb8863` checked `docs`, `README.md`, and `benchmark-results` for signed URL token parameters, private-key blocks, and the prohibited raw sample. Result: `PASS no focused high-risk leaks across public docs/benchmark-results (70 files)`. |

## Overall Verdict

PASS — all automatable checks passed after updating the canonical contract to match the implementation-backed TEI/default safe `/health.runtime` metadata behavior.

## Notes

The previous FAIL was caused by stale documentation in `docs/same-host-embedding-service-contract.md`, not by a runtime/test failure. The contract now aligns with `api/handlers/health_test.go` and `api/main.go`: TEI/default `/health` exposes safe operational metadata while remaining explicitly not a live inference readiness probe.

Evidence run:

- `gsd_exec` `8702d1d2-9c57-4444-b6b2-8f646efb8863`: full rerun of artifact-driven UAT after remediation; exit code 0.