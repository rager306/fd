---
id: S03
parent: M045-d0e5xq
milestone: M045-d0e5xq
provides:
  - Validated R028 bounded startup posture.
  - Working TEI container using local snapshot path.
  - GSD parser gotcha for future UAT completion attempts.
requires:
  - slice: S01
    provides: Read-only recon and TEI source/docs findings.
  - slice: S02
    provides: Cache completeness inventory and failed offline candidate context.
affects:
  []
key_files:
  - docker-compose.yaml
  - docker-compose.override.yaml
  - docs/same-host-embedding-service-contract.md
  - documents/tei-startup-mitigation-m045.md
  - benchmark-results/m045-tei-local-path-startup-proof.md
key_decisions:
  - Use cached local USER-bge-m3 snapshot path as TEI `--model-id`; reject `HF_HUB_OFFLINE=1` as ineffective for this TEI startup issue.
  - GSD UAT content passed to `gsd_slice_complete` must follow the exact UAT template section shape for mode extraction.
patterns_established:
  - For TEI Docker deployments, `/data` volume prevents repeated weight downloads, but passing Hub ID still uses Hub/api_repo path. Use local snapshot path to avoid Hub resolution and remote ONNX probe.
  - Compose override command replacement is total; keep base and override commands aligned.
  - GSD UAT mode parser expects `## UAT Type` plus a `- UAT mode:` bullet; prose headers such as `UAT Type: mixed` are not recognized.
observability_surfaces:
  - benchmark-results/m045-tei-local-path-startup-proof.md records preflight, proof result, logs, health timing, and smoke results.
  - Browser timeline records localhost `/health` verification.
drill_down_paths:
  []
duration: ""
verification_result: passed
completed_at: 2026-06-14T12:50:13.274Z
blocker_discovered: false
---

# S03: Controlled startup proof

**Validated local snapshot `--model-id` as the TEI startup mitigation and rejected `HF_HUB_OFFLINE=1`.**

## What Happened

S03 first recorded that the earlier `HF_HUB_OFFLINE=1` candidate failed: TEI remained unhealthy after a 15-minute timeout and still entered `Downloading onnx/model.onnx`. The slice then switched to the official TEI local-model path approach using the cached USER-bge-m3 snapshot directory as `--model-id`. Indexed TEI source and official CLI docs showed this should take TEI's local model branch instead of the Hub repo path. The controlled restart proof applied the local snapshot command through both compose files, measured TEI startup, and verified fd/direct TEI smoke. TEI reached healthy at 2026-06-14T12:15:15 after container start at 2026-06-14T12:12:14. R028 was validated. The earlier GSD completion blocker was traced to UAT parser format: `extractUatType` only recognizes a `## UAT Type` section containing a `- UAT mode:` bullet.

## Verification

Fresh runtime check before completion: `fd_tei` is running and healthy with command `--model-id /data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae --max-batch-tokens 8192`; fd `/health` reports backend `tei`, model `deepvk/USER-bge-m3`, dimensions 1024; browser assertions against `http://localhost:8000/health` passed. Prior final gates also passed: `go test ./...` 270 tests, golangci-lint 0 issues, govulncheck 0 reachable vulnerabilities. S03 mixed UAT attempt 3 is PASS.

## Requirements Advanced

None.

## Requirements Validated

- R028 — Local snapshot startup proof reached TEI healthy and fd/direct TEI smokes passed; browser health verification passed; evidence in benchmark-results/m045-tei-local-path-startup-proof.md and S03 UAT attempt 3.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

Original S03 offline-cache proof failed and was replaced with local snapshot path proof. The first slice completion attempts also failed because the UAT body did not match the parser's exact `## UAT Type` / `- UAT mode:` format; completion was retried with the template-compliant UAT content.

## Known Limitations

TEI still logs an immediate local ORT failure because ONNX files are absent, but it no longer blocks on remote Hub ONNX probing and reaches healthy. If the cached snapshot revision changes, compose must be updated and proof rerun.

## Follow-ups

None required for R028. Future improvement could replace the pinned snapshot path with an operator-managed stable symlink if desired.

## Files Created/Modified

- `docker-compose.yaml` — TEI command now uses cached local USER-bge-m3 snapshot path and max-batch-tokens 8192.
- `docker-compose.override.yaml` — Override command aligned with base compose local snapshot path.
- `docs/same-host-embedding-service-contract.md` — Documents local snapshot mode as startup mitigation.
- `documents/tei-startup-mitigation-m045.md` — Records validated local snapshot outcome and rejected offline env candidate.
- `benchmark-results/m045-tei-local-path-startup-proof.md` — Runtime proof record.
