---
id: S02
parent: M031-gn517a
milestone: M031-gn517a
provides:
  - Durable source contract for future hosted ONNX workflow proof.
requires:
  []
affects:
  []
key_files:
  - docs/onnx-artifacts/user-bge-m3-dense-fp32.json
  - docs/onnx-artifacts/hf-tokenizers-linux-amd64.json
  - docs/onnx-artifacts/PROVISIONING.md
  - benchmark-results/fd-onnx-source-contract-m031-s02.txt
key_decisions:
  - D029: pinned supporting artifact candidates are selected, but exported ONNX model binary remains blocked until immutable exact hosting or separate reproducible-export validation exists.
patterns_established:
  - Use `source_contract` metadata to separate pinned candidates from rollout-proof evidence.
  - Keep exact exported ONNX model binary distinct from upstream model source files.
observability_surfaces:
  - Tracked source_contract fields in artifact manifests.
  - Provisioning documentation source-selection table.
  - Benchmark-results outcome artifact.
drill_down_paths:
  - .gsd/milestones/M031-gn517a/slices/S02/tasks/T01-SUMMARY.md
  - .gsd/milestones/M031-gn517a/slices/S02/tasks/T02-SUMMARY.md
  - .gsd/milestones/M031-gn517a/slices/S02/tasks/T03-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-21T06:41:33.507Z
blocker_discovered: false
---

# S02: Source contract documentation and closure

**Persisted and verified the ONNX artifact source contract without production promotion.**

## What Happened

S02 persisted the M031 source contract into tracked manifests and provisioning docs, recorded an outcome artifact and decision, and ran final verification. The work converts vague artifact-source blockers into specific source statuses: native tokenizer, tokenizer JSON, and ONNX Runtime have pinned checksum-matched candidates; the exact exported ONNX model binary remains blocked. This prepares the next gate without changing runtime defaults or taking external actions.

## Verification

S02 verification passed: JSON valid; provisioning dry-run/verifier passed; Go tests 87 passed; GolangCI-Lint 0 issues; actionlint passed; tagged tokenizer tests 20 passed; tagged ONNX smoke tests 2 passed; Docker default build succeeded; docs/outcome leak checks passed; tracked binary hygiene passed; GitNexus detected only low-risk docs changes and no affected processes; no background processes; port 18000 clean.

## Requirements Advanced

None.

## Requirements Validated

None.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

Initial verification-command mistakes were corrected and rerun: provisioning dry-run does not support `--allow-missing`, and binary hygiene must check tracked files rather than ignored runtime cache artifacts.

## Known Limitations

No hosted workflow proof was run. ONNX model artifact remains local-only. ONNX remains opt-in experimental and TEI remains production/default.

## Follow-ups

Next milestone can either mirror/upload the exact ONNX model binary to an immutable non-secret source, or design a reproducible export workflow that regenerates and revalidates the artifact. Hosted workflow proof still requires explicit push approval and real source inputs.

## Files Created/Modified

- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json` — Added source_contract for ONNX model blocker, tokenizer JSON candidate, and ONNX Runtime candidate.
- `docs/onnx-artifacts/hf-tokenizers-linux-amd64.json` — Replaced mutable latest source with pinned v1.27.0 release source candidate and archive checksum metadata.
- `docs/onnx-artifacts/PROVISIONING.md` — Documented M031 source selection status and remaining ONNX model blocker.
- `benchmark-results/fd-onnx-source-contract-m031-s02.txt` — Outcome artifact for source contract gate.
- `.gsd/DECISIONS.md` — Decision D029 added by GSD.
