# M013-nhu1x9: HF tokenizer native packaging gate

**Vision:** Turn M012's parity-correct HF Rust tokenizer finding into a safe, reproducible opt-in packaging/build path that preserves default TEI builds and only unlocks ONNX benchmarking after semantic equivalence is proven.

## Success Criteria

- Native HF tokenizer artifact provenance/checksum contract exists.
- Default TEI builds remain unaffected by native tokenizer dependency.
- An opt-in build-tag/native path is proven or blocked with evidence.
- Tokenizer parity remains the gate before ONNX cosine, and cosine remains the gate before performance benchmarking.
- No native binary, raw probe text, or production runtime switch is committed.

## Slices

- [x] **S01: S01** `risk:medium` `depends:[]`
  > After this: After this, the native HF tokenizer artifact has a checksum/provenance contract and reproducible local setup instructions.

- [ ] **S02: Opt in build tag boundary** `risk:high` `depends:[S01]`
  > After this: After this, default builds stay clean and the project has an explicit opt-in tokenizer build path or a packaging blocker.

- [ ] **S03: Parity correct ONNX tokenizer integration** `risk:high` `depends:[S02]`
  > After this: After this, the ONNX tokenizer path either uses the parity-correct binding under opt-in conditions or remains blocked with exact evidence.

- [ ] **S04: Final native packaging gate decision** `risk:medium` `depends:[S03]`
  > After this: After this, the project has a final gate decision: proceed to ONNX benchmark, continue packaging work, or defer ONNX runtime.

## Boundary Map

| Area | In scope | Out of scope |
|---|---|---|
| Native tokenizer artifact | Define checksum/version/architecture manifest for `libtokenizers.a` and local artifact handling | Commit large/native binary blobs without policy |
| Build integration | Add opt-in build tag or isolated package path so default TEI builds do not require native tokenizers | Make native tokenizer mandatory for default builds |
| ONNX tokenizer path | Wire parity-correct tokenizer only behind explicit opt-in build/runtime path if packaging works | Switch production/default runtime to ONNX |
| Verification | Tagged parity tests, default tests/lint, artifact checks, isolated-cache cosine if runtime integration succeeds | ONNX throughput benchmark before cosine equivalence |
| Deployment docs | Document local/Docker/CI requirements for native tokenizer artifact | Full production rollout or remote publishing |
