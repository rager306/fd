---
id: M015-22msl0
title: "Russian legal retrieval quality gate"
status: complete
completed_at: 2026-05-20T05:09:04.184Z
key_decisions:
  - D011: use the user-provided 44-ФЗ JSONL as the first Russian/legal parity corpus.
  - D012: block ONNX packaging/tuning until long-text legal quality divergence is investigated.
key_files:
  - tests/44-FZ-2026-articles.jsonl
  - tools/evaluate_legal_retrieval.py
  - benchmark-results/fd-legal-corpus-profile-m015-s01.txt
  - benchmark-results/fd-legal-retrieval-dry-run-m015-s02.txt
  - benchmark-results/fd-legal-retrieval-m015-s03.txt
  - benchmark-results/fd-legal-retrieval-m015-summary.txt
  - .gsd/DECISIONS.md
  - .gsd/milestones/M015-22msl0/M015-22msl0-VALIDATION.md
lessons_learned:
  - The tagged ONNX path can be fast but still fail legal-corpus vector equivalence on longer texts.
  - Ranking parity can look mostly acceptable while per-vector cosine has severe outliers; both signals are needed.
  - Unlabeled legal corpora are useful for TEI-vs-ONNX parity gates, but explicit qrels are still needed for absolute retrieval quality claims.
---

# M015-22msl0: Russian legal retrieval quality gate

**M015 executed the 44-ФЗ Russian/legal quality gate and blocked tagged ONNX due severe long-text vector divergence.**

## What Happened

M015 used the user-provided 44-ФЗ JSONL corpus to run the first Russian/legal quality gate for tagged ONNX. S01 profiled the corpus and defined the parity/known-item gate. S02 implemented a sanitized evaluator. S03 ran the evaluator against TEI default and tagged ONNX, fixed an ID fallback issue, reran the gate, and produced a FAIL artifact. S04 summarized the verdict, recorded decision D012, and ran final verification. The result blocks ONNX packaging/tuning as the next priority: severe cross-backend cosine outliers on longer legal text must be investigated first. TEI remains production/default.

## Success Criteria Results

- Corpus profiled and hashed: pass.
- Repeatable evaluator implemented: pass.
- TEI and tagged ONNX compared with isolated ONNX namespace: pass.
- Artifacts avoid raw legal text/secrets: pass.
- Clear quality gate verdict: pass, verdict is ONNX FAIL.
- No production switch: pass.

## Definition of Done Results

- [x] Corpus profile artifact records counts, schema assumptions, and hash without raw text dump.
- [x] Evaluator produces sanitized TEI-vs-ONNX legal retrieval parity metrics.
- [x] Tagged ONNX evaluated with isolated namespace and M013/M014 context.
- [x] Final recommendation states ONNX fails this first Russian/legal parity gate.
- [x] No production/default switch occurred.
- [x] Tests/lint/tagged tests/artifact hygiene/GitNexus checks passed.

## Requirement Outcomes

- Russian/legal quality: gate executed; tagged ONNX failed strict equivalence.
- Production readiness: ONNX remains blocked; TEI remains default.
- Evaluator reproducibility: validated with corpus hash and sanitized artifacts.
- Artifact hygiene: validated; raw legal text leaks 0.

## Deviations

S03 discovered an evaluator ID fallback bug for unnumbered subclauses; it was fixed and the gate rerun. The final gate verdict remains FAIL. M015 closes as successful evidence production, not ONNX approval.

## Follow-ups

Next milestone should investigate long-text legal divergence before packaging/tuning: compare TEI truncation/tokenization behavior, inspect ONNX max_sequence_length=128 effects, test longer sequence ONNX export or chunking policy, and rerun the legal gate.
