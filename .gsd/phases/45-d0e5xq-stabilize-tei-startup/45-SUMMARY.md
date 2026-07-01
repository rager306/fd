---
id: M045-d0e5xq
title: "Stabilize TEI startup"
status: complete
completed_at: 2026-06-14T12:55:10.285Z
key_decisions:
  - Use cached local USER-bge-m3 snapshot path as TEI `--model-id` for startup stabilization.
  - Reject `HF_HUB_OFFLINE=1` as the startup mitigation because it did not stop TEI from entering ONNX/model probing and becoming unhealthy.
  - Keep fd public runtime identity as `deepvk/USER-bge-m3` while TEI internally uses a local path.
  - GSD UAT closeout content must follow the exact UAT template section shape for mode extraction.
key_files:
  - docker-compose.yaml
  - docker-compose.override.yaml
  - documents/tei-startup-recon-m045.md
  - documents/tei-startup-mitigation-m045.md
  - docs/same-host-embedding-service-contract.md
  - benchmark-results/m045-tei-offline-startup-proof.md
  - benchmark-results/m045-tei-local-path-startup-proof.md
  - .gsd/milestones/M045-d0e5xq/M045-d0e5xq-VALIDATION.md
  - .gsd/milestones/M045-d0e5xq/slices/S03/S03-UAT.md
lessons_learned:
  - A mounted Hugging Face cache prevents repeated weight downloads, but a Hub ID still keeps TEI on the Hub/api_repo path; use a local directory path when the operator wants local-only startup behavior.
  - Compose override `command` replaces the base command wholesale, so startup mitigations must update both base and override when both define `command`.
  - TEI may still log a local ORT/ONNX miss when ONNX files are absent; the operational criterion is that it falls through quickly to Candle backend and becomes healthy.
  - GSD UAT parser defaults to artifact-driven unless the mode appears under `## UAT Type` as `- UAT mode: ...`.
---

# M045-d0e5xq: Stabilize TEI startup

**M045 made TEI restarts operationally bounded by switching the deployment from Hub ID startup to a validated cached local USER-bge-m3 snapshot path.**

## What Happened

M045 started from an operational failure mode: TEI-only fd deployments could spend ~45+ minutes in startup due to TEI's internal ONNX/ORT probing and Hub-style artifact resolution, even though fd's own ONNX runtime had been removed. S01 captured the current TEI image, command, env, mounted cache/model layout, logs, and smoke behavior without disrupting the running service. S02 inventoried the cache and staged the `HF_HUB_OFFLINE=1` candidate with rollback planning. S03 then executed the controlled proof: `HF_HUB_OFFLINE=1` was rejected after TEI still became unhealthy and attempted `Downloading onnx/model.onnx`; the selected mitigation became passing the cached local snapshot directory to TEI as `--model-id`. With `/data/models--deepvk--USER-bge-m3/snapshots/0cc6cfe48e260fb0474c753087a69369e88709ae`, TEI reached healthy in roughly three minutes, fd `/health` preserved backend/model/dimension contract, and fd/direct TEI embedding smokes returned 1024-dimensional vectors. Operator docs, compose base, and compose override now align on the validated local snapshot command. R028 is validated. During closeout, a GSD parser gotcha was found and documented: UAT mode extraction requires the exact `## UAT Type` plus `- UAT mode:` template shape.

## Success Criteria Results

- MET: Effective current TEI image, command, env, mounted model/cache layout, and startup log timeline were documented in S01/S03 artifacts without disrupting the service.
- MET: Candidate mitigations were evaluated: `HF_HUB_OFFLINE=1` rejected empirically; local snapshot path accepted via TEI docs, source behavior, and runtime proof.
- MET: Controlled startup proof measured TEI start to healthy and verified fd readiness and embedding smoke.
- MET: Compose and operator docs were updated only after local snapshot proof validated the safer configuration.
- MET: Fresh mandatory gates passed: go test 270 passed, golangci-lint 0 issues, govulncheck 0 reachable vulnerabilities.

## Definition of Done Results

- All planned slices complete: S01, S02, S03.
- Fresh runtime verification passed: `fd_tei` is healthy with local snapshot command; fd `/health` reports TEI/deepvk/1024; fd embedding smoke returns 1024 dimensions.
- Fresh Go gates passed: `go test ./...`, golangci-lint, govulncheck.
- R028 validated with proof artifact and UAT evidence.
- Milestone validation saved with verdict `pass`.

## Requirement Outcomes

- R028: active to validated. Proof: S03 local snapshot startup proof, mixed UAT, browser localhost `/health` assertion, and closeout runtime smoke.
- No M045-specific requirement remains unaddressed.

## Deviations

The planned `HF_HUB_OFFLINE=1` mitigation was rejected after controlled proof failure. S03 pivoted to the local snapshot path after TEI source/docs showed that local directories bypass Hub api_repo behavior. A GSD completion-format blocker was resolved by reauthoring UAT content to match the parser template.

## Follow-ups

Optional future hardening: replace the pinned snapshot path with an operator-managed stable symlink and rerun the same proof. No follow-up is required to satisfy M045/R028.
