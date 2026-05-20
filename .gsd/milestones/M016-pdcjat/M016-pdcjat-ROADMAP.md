# M016-pdcjat: ONNX legal divergence and model alternatives

**Vision:** Explain the M015 long-text legal quality failure and decide the next evidence-backed path: fix the current ONNX route, introduce chunking/longer sequence export, and separately identify model alternatives worth legal-corpus benchmarking.

## Success Criteria

- M015 legal divergence is narrowed to a likely cause or a concrete unknown.
- Remediation path is recommended before packaging/tuning resumes.
- Alternative model candidates are researched and ranked for future legal benchmark trials.
- Artifacts are sanitized and avoid raw legal text/secrets.
- TEI remains production/default and no model switch occurs.

## Slices

- [x] **S01: S01** `risk:medium` `depends:[]`
  > After this: After this, the M015 failure cases are concrete and diagnosable by ID/hash/length/token counts.

- [x] **S02: S02** `risk:high` `depends:[]`
  > After this: After this, we know whether max sequence length 128 is the likely cause or whether another ONNX path issue remains.

- [ ] **S03: Remediation option assessment** `risk:medium` `depends:[S02]`
  > After this: After this, there is a concrete next remediation plan for ONNX quality.

- [x] **S04: S04** `risk:medium` `depends:[]`
  > After this: After this, alternative models are ranked for future legal-corpus benchmarking.

## Boundary Map

| Area | In scope | Out of scope |
|---|---|---|
| ONNX divergence | Diagnose long-text legal divergence between TEI and tagged ONNX, especially max sequence length and truncation behavior | Production/default ONNX switch |
| Tokenization/truncation | Compare TEI/HF tokenizer behavior, ONNX tokenizer behavior, and sequence-length assumptions on worst M015 IDs | Broad tokenizer library rewrite unless evidence requires it |
| Remediation options | Evaluate longer-sequence ONNX export feasibility and explicit chunking policy as alternatives | Full production packaging before quality passes |
| Alternative models | Research Russian/legal-capable embedding candidates and define benchmark criteria for this environment | Replacing `deepvk/USER-bge-m3` without corpus benchmark evidence |
| Evidence | Save sanitized artifacts with IDs/hashes/metrics only | Printing raw 44-ФЗ text in artifacts |
