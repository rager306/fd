---
id: M016-pdcjat
title: "ONNX legal divergence and model alternatives"
status: complete
completed_at: 2026-05-20T07:21:04.356Z
key_decisions:
  - D014: Keep TEI as production/default and make the next ONNX gate a 512-token quality remediation plus long-text chunking/longer-sequence policy before packaging, tuning, or promotion.
  - Alternative models remain research candidates only; no model switch without Russian/legal corpus benchmark evidence.
key_files:
  - tools/profile_legal_divergence.py
  - tools/diagnose_onnx_sequence_length.py
  - benchmark-results/fd-legal-divergence-profile-m016-s01.txt
  - benchmark-results/fd-onnx-sequence-diagnostics-m016-s02.txt
  - benchmark-results/fd-onnx-remediation-plan-m016-s03.txt
  - benchmark-results/fd-model-alternatives-m016-s04.txt
  - .gsd/DECISIONS.md
lessons_learned:
  - Tokenizer parity is necessary but not sufficient; sequence-length policy is part of embedding correctness for legal text.
  - Worst-case legal divergence should be diagnosed with stable IDs, hashes, and token counts before runtime changes.
  - A 512-token ONNX path likely fixes most observed worst-case divergence but still needs chunking or longer sequence handling for texts beyond 512 tokens.
---

# M016-pdcjat: ONNX legal divergence and model alternatives

**M016 proved ONNX legal divergence is primarily 128-token truncation and selected a safe quality-first remediation path while keeping TEI default.**

## What Happened

M016 investigated why tagged ONNX failed the Russian/legal quality gate despite tokenizer parity and fixed-probe success. S01 extracted and profiled the M015 worst cases without leaking raw legal text, showing all 17 severe outliers are truncated at 128 tokens and only 2 are truncated at 512. S02 compared TEI API vectors to local ONNX outputs at 128 and 512 tokens, confirming the severe divergence is primarily caused by the 128-token ONNX path: mean cosine improved from 0.9204953 at 128 to 0.99885631 at 512. S03 turned that evidence into a remediation plan and D014: keep TEI default, keep ONNX experimental, validate a 512-token quality gate first, and add chunking or longer-sequence policy for >512-token legal fragments before packaging or performance claims. S04 provided alternative-model research but did not change the default model strategy.

## Success Criteria Results

- Worst-case profile: PASS.
- Sequence-length root-cause diagnostics: PASS.
- Remediation path decision: PASS.
- Alternative model planning: PASS.
- Production/default safety: PASS — TEI remains default.
- Verification: PASS — script compile, hygiene, Go tests, lint, tagged tests, and validation passed.

## Definition of Done Results

- Evidence artifacts produced for S01/S02/S03/S04: met.
- No raw legal corpus text in artifacts: met by hygiene checks.
- TEI remains production/default: met by D014 and S03.
- ONNX remains opt-in/experimental: met by D014 and S03.
- Verification run before closure: met; Go tests, lint, tagged tests, script compile, and hygiene checks passed.
- GitNexus scope check before closure: met.

## Requirement Outcomes

- ONNX legal divergence root cause: advanced from unknown to confirmed/narrowed.
- ONNX production readiness: remains blocked until 512-token/chunking remediation and full legal gate pass.
- Alternative model strategy: documented as future benchmark candidates only.
- Raw text artifact hygiene: validated for M016 artifacts.

## Deviations

S04 alternative-model research was completed before S01/S02/S03 because it was independent. This did not affect the root-cause investigation and was integrated into S03 as a non-switch recommendation.

## Follow-ups

Future implementation milestone: build and validate a 512-token ONNX quality gate for `deepvk/USER-bge-m3`, add deterministic chunking or longer-sequence handling for legal fragments beyond 512 tokens, rerun the full legal retrieval gate with isolated Redis namespace, then benchmark performance and packaging only after quality passes.
