---
verdict: pass
remediation_round: 0
---

# Milestone Validation: M016-pdcjat

## Success Criteria Checklist
- [x] Extract worst M015 divergence IDs without raw text leakage — S01 complete.
- [x] Profile token lengths and truncation flags — S01 complete.
- [x] Compare sequence-length behavior against TEI — S02 complete.
- [x] Choose remediation path — S03 complete.
- [x] Research alternative Russian/legal-capable models without switching defaults — S04 complete.
- [x] Preserve TEI as production/default — D014 and S03 complete.
- [x] Keep ONNX opt-in/experimental — D014 and S03 complete.
- [x] Verify scripts/tests/lint/hygiene — S03 T03 evidence complete.

## Slice Delivery Audit
| Slice | Claimed output | Delivered evidence | Verdict |
|---|---|---|---|
| S01 | Worst-case divergence profile | `benchmark-results/fd-legal-divergence-profile-m016-s01.txt`; 17/17 resolved and truncated at 128 | PASS |
| S02 | Sequence-length root-cause diagnostics | `benchmark-results/fd-onnx-sequence-diagnostics-m016-s02.txt`; 128 mean cosine 0.9204953 vs 512 mean 0.99885631 | PASS |
| S03 | Remediation decision | `benchmark-results/fd-onnx-remediation-plan-m016-s03.txt`; D014 saved | PASS |
| S04 | Alternative model research | `benchmark-results/fd-model-alternatives-m016-s04.txt` | PASS |

## Cross-Slice Integration
| Boundary | Result |
|---|---|
| S01 -> S02 | PASS: S01 produced the 17 worst IDs and token/truncation profile used by S02 sequence diagnostics. |
| S02 -> S03 | PASS: S02 confirmed 128-token truncation and quantified 512 improvement; S03 used that evidence to choose remediation. |
| S04 -> S03 | PASS: S04 alternative-model research was explicitly kept as future research only; S03 rejected model switch for this remediation. |

No cross-slice mismatch found. The milestone stayed within investigation/planning scope and did not switch production defaults.

## Requirement Coverage
- ONNX legal divergence investigation: covered by S01 and S02.
- Remediation path selection: covered by S03 and D014.
- Alternative model planning: covered by S04.
- Production safety: covered by repeated decision that TEI remains default and ONNX stays experimental.

Unaddressed by design: implementation of a 512-token ONNX Go/runtime path and full legal gate rerun. These are follow-up implementation milestones, not validation failures for M016.

## Verification Class Compliance
- Artifact hygiene: PASS (`raw_legal_text_leaks=0`).
- Python script compile: PASS.
- Default Go tests: PASS (`78 passed in 4 packages`).
- Pinned GolangCI-Lint: PASS (`0 issues`).
- Tagged HF tokenizer tests: PASS (`20 passed in 1 package`).
- GitNexus scope check: PASS (low/no changed symbols for uncommitted scope at validation time).


## Verdict Rationale
M016 achieved its investigation and remediation-planning goal. It identified ONNX 128-token truncation as the main severe legal divergence cause, quantified the 512-token improvement, selected a safe remediation path, and kept production defaults unchanged. Remaining work is future implementation, not incomplete M016 scope.
