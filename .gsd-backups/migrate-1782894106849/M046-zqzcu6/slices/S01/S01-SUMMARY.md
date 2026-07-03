---
id: S01
parent: M046-zqzcu6
milestone: M046-zqzcu6
provides:
  - Verified P0/P1 issue #3 inventory.
  - S02 batch endpoint guardrail starting checklist.
  - Root decision map for S02-S05.
requires:
  []
affects:
  []
key_files:
  - .gsd/milestones/M046-zqzcu6/slices/S01/S01-RESEARCH.md
  - benchmark-results/m046-s01-audit-validation.md
  - documents/issue-3-audit-remediation-plan-m046.md
key_decisions:
  - Treat issue #3 as a wave remediation program, not a single broad rewrite.
  - Fix P0 batch guardrails before N+1 batch performance work.
  - Separate same-host policy risks from implementation defects before changing auth defaults.
patterns_established:
  - Validate audit findings against current code before coding fixes.
  - Use static probes instead of abusive load tests for DoS-class findings.
observability_surfaces:
  - benchmark-results/m046-s01-audit-validation.md records static validation signals.
  - documents/issue-3-audit-remediation-plan-m046.md records remediation waves and root decisions.
drill_down_paths:
  []
duration: ""
verification_result: passed
completed_at: 2026-06-14T15:51:08.337Z
blocker_discovered: false
---

# S01: Audit validation map

**Validated issue #3 P0/P1 findings against current code and mapped confirmed defects to remediation waves.**

## What Happened

S01 converted issue #3 from a broad audit report into a verified remediation map. It inventoried all 10 P0/P1 findings, confirmed the main defect signals with static probes, and separated true implementation defects from policy risks and lower-priority cleanup. The slice identified three root decision clusters: batch endpoints evolved outside the hardened `/v1/embeddings` request boundary, same-host assumptions leaked into default exposure behavior, and LocalCache used simple concurrency primitives without deterministic lifecycle/accounting. No product source behavior was changed in S01; it prepared S02-S05 execution inputs.

## Verification

Static probe `cafb3f84-b852-4d21-b71f-c13c9f5afd77` confirmed current-code signals for batch route guardrails, default-open auth, public metrics, LocalCache no Close/size counter, and per-input TEI calls. Artifact check `b013c88f-400c-4627-ae21-08fcb2c83ac0` verified P0/P1 coverage and S02 inputs. S01 UAT passed with evidence `50b0abab-5c25-42e0-bce8-b98a57d65b98`, `af89d6de-6281-4f28-9ed5-9c4bf70c799f`, and `85b897af-c19f-4d0c-9acc-2e3b3841de68`.

## Requirements Advanced

- R029 — Confirmed batch endpoint guardrail defects and prepared S02 checklist.
- R030 — Confirmed exposure posture policy risk and prepared S04 policy inputs.
- R031 — Confirmed LocalCache lifecycle/accounting risk and prepared S05 inputs.

## Requirements Validated

None.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

None. S01 intentionally avoided source edits and destructive load tests.

## Known Limitations

P1 #9 still needs exact contract confirmation in S06; P2/P3 findings are intentionally deferred to closure triage after P0/P1 waves.

## Follow-ups

Proceed to S02 to fix batch endpoint guardrails first.

## Files Created/Modified

- `.gsd/milestones/M046-zqzcu6/slices/S01/S01-RESEARCH.md` — Validated issue #3 P0/P1 inventory and wave mapping.
- `benchmark-results/m046-s01-audit-validation.md` — Static validation evidence for confirmed P0/P1 signals.
- `documents/issue-3-audit-remediation-plan-m046.md` — Root-decision map and downstream remediation plan.
