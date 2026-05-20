# M012-3edtlz: Tokenizer parity gate

**Vision:** Resolve the M011 tokenizer-parity blocker by proving whether Go tokenization can exactly match Hugging Face tokenizer output for `deepvk/USER-bge-m3`, so ONNX embedding comparisons become meaningful before any performance claims.

## Success Criteria

- A sanitized Hugging Face tokenizer baseline exists for fixed Russian/legal probes.
- Go tokenizer output is compared token-by-token against the baseline.
- The milestone produces either exact tokenizer parity or a concrete blocker with evidence.
- No ONNX throughput benchmark is run unless tokenizer parity and cosine equivalence pass.
- TEI remains default and ONNX remains opt-in.
- No raw probe text or large model artifacts are committed.

## Slices

- [x] **S01: S01** `risk:medium` `depends:[]`
  > After this: After this, the project has a sanitized Hugging Face tokenizer baseline for fixed Russian/legal probes.

- [x] **S02: S02** `risk:high` `depends:[]`
  > After this: After this, current and candidate Go tokenizers are compared against the HF baseline with exact mismatch evidence.

- [x] **S03: S03** `risk:high` `depends:[]`
  > After this: After this, Go tokenizer parity either passes through a tested HF Rust tokenizers binding or is blocked with concrete linking/correctness evidence.

- [ ] **S04: Final parity decision and ONNX gate** `risk:medium` `depends:[S03]`
  > After this: After this, the project either has a valid post-parity ONNX cosine comparison or a final recommendation to defer ONNX benchmarking.

## Boundary Map

| Area | In scope | Out of scope |
|---|---|---|
| Tokenizer correctness | Compare Hugging Face Python tokenization against Go candidates for fixed Russian/legal probes | Replace the embedding model |
| Evidence artifacts | Persist labels, lengths, token hashes, token IDs when safe, and mismatch summaries without raw probe text | Print raw legal probe text in artifacts |
| Go implementation | Adjust or replace Go tokenizer path only if parity can be proven | Rewrite handlers/cache or switch production runtime |
| ONNX validation | Rerun Go ONNX vs TEI cosine only after tokenizer parity passes, with isolated cache | ONNX throughput benchmarking before equivalence |
| Runtime | Keep TEI default and ONNX opt-in | Make ONNX production default |
