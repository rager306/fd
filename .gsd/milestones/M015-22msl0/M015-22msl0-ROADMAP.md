# M015-22msl0: Russian legal retrieval quality gate

**Vision:** Use the provided 44-FZ legal JSONL corpus to test whether the tagged ONNX path preserves TEI retrieval behavior on Russian legal text before investing in production packaging or tuning.

## Success Criteria

- The new JSONL corpus is profiled and hashed.
- A repeatable legal retrieval parity evaluation is implemented or scripted.
- TEI and tagged ONNX are compared on the corpus without Redis contamination.
- Artifacts avoid raw legal text and secrets.
- The milestone ends with a clear quality gate verdict and no production switch.

## Slices

- [x] **S01: S01** `risk:medium` `depends:[]`
  > After this: After this, the legal corpus and evaluation method are explicit and auditable.

- [x] **S02: S02** `risk:medium` `depends:[]`
  > After this: After this, there is a reusable local evaluator for the provided legal JSONL file.

- [x] **S03: S03** `risk:high` `depends:[]`
  > After this: After this, ONNX has a real first-pass Russian/legal retrieval parity result.

- [ ] **S04: Quality verdict and closure** `risk:medium` `depends:[S03]`
  > After this: After this, the project knows whether ONNX can proceed to packaging/tuning or needs quality work.

## Boundary Map

| Area | In scope | Out of scope |
|---|---|---|
| Corpus | Use `tests/44-FZ-2026-articles.jsonl`; derive document/query IDs without raw text in artifacts | Curating external legal datasets |
| Quality gate | Compare TEI vs tagged ONNX rankings, top-k overlap, rank correlation, and synthetic known-item retrieval | Claiming absolute legal relevance without labeled qrels |
| Runtime | Use TEI default and tagged ONNX opt-in with isolated Redis namespaces | Production default switch |
| Artifacts | Save sanitized results with config, corpus hash, metrics, and caveats | Printing raw legal text in artifacts |
| Tooling | Add minimal evaluator script if needed | Full benchmark framework rewrite |
