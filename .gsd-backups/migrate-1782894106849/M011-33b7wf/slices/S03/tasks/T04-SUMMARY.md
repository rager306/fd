---
id: T04
parent: S03
milestone: M011-33b7wf
key_files:
  - benchmark-results/fd-go-onnx-m011-s03.txt
  - api/embed/onnx.go
  - api/main.go
key_decisions:
  - Do not claim Go ONNX API equivalence yet; isolated-cache comparison failed cosine threshold.
  - Treat Go tokenizer mismatch as the S03 blocker. ONNX Runtime loads/runs, but `sugarme/tokenizer` token IDs differ from Hugging Face tokenizer for the same Russian legal probe.
  - Stop short of production/prototype success; S04 should recommend resolving tokenization parity before benchmarking throughput.
duration: 
verification_result: mixed
completed_at: 2026-05-19T19:26:25.426Z
blocker_discovered: true
---

# T04: Ran the opt-in Go ONNX API and found a real blocker: Go tokenizer output diverges from Hugging Face tokenizer, causing TEI-vs-ONNX cosine failures.

**Ran the opt-in Go ONNX API and found a real blocker: Go tokenizer output diverges from Hugging Face tokenizer, causing TEI-vs-ONNX cosine failures.**

## What Happened

Ran the opt-in Go ONNX API on port 18000 with the local ONNX artifact, tokenizer, and ONNX Runtime shared library. The first comparison was invalid because shared Redis cache returned TEI vectors, so the ONNX server was restarted with isolated cache namespace `EMBEDDING_CACHE_VERSION=m011-onnx`. With isolated cache, the Go ONNX backend executed and returned normalized 1024-dimensional vectors, but TEI-vs-Go-ONNX cosine failed the `0.999` threshold: observed values ranged from `0.98266755` to `0.99713198`. A quick tokenization probe found the likely root cause: Hugging Face Python tokenizer produced 21 token IDs for the labor-law probe, while `sugarme/tokenizer` produced 27 IDs with divergent tokenization. This is a real blocker for ONNX API equivalence; runtime loading itself is not the blocker.

## Verification

ONNX API server started successfully with isolated cache and explicit ONNX config. Comparison artifact was saved as FAIL due cosine below threshold. Tokenization probe confirmed Python HF tokenizer and Go sugarme tokenizer produce different IDs for the same Russian text.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `bg_shell start ONNX API on port 18000 with EMBEDDING_BACKEND=onnx and default cache namespace` | 1 | ⚠️ invalid comparison setup — shared Redis returned TEI cache hits | 0ms |
| 2 | `bg_shell start ONNX API on port 18000 with EMBEDDING_CACHE_VERSION=m011-onnx` | 0 | ✅ pass — server ready | 0ms |
| 3 | `python3 TEI default vs Go ONNX opt-in comparison script writing benchmark-results/fd-go-onnx-m011-s03.txt` | 2 | ❌ fail — cosine range 0.98266755 to 0.99713198 below 0.999 threshold | 0ms |
| 4 | `uv run --python 3.13 --with transformers --with torch --with sentencepiece Python tokenizer probe + Go sugarme tokenizer probe` | 0 | ❌ mismatch evidence — Python produced 21 token IDs; Go sugarme produced 27 token IDs for same probe | 0ms |
| 5 | `bg_shell kill ONNX API server` | 0 | ✅ cleanup | 0ms |

## Deviations

The first ONNX API comparison used the default Redis namespace and falsely matched TEI hashes because cached TEI vectors were returned. The ONNX server was restarted with `EMBEDDING_CACHE_VERSION=m011-onnx` to isolate cache and force real ONNX inference.

## Known Issues

Go ONNX backend returns normalized 1024-dimensional vectors, but they do not match TEI/Python ONNX baseline sufficiently. Token probe showed Python HF tokenizer produced 21 IDs while Go sugarme tokenizer produced 27 IDs for the same Russian legal text, with divergent token IDs after the first few tokens.

## Files Created/Modified

- `benchmark-results/fd-go-onnx-m011-s03.txt`
- `api/embed/onnx.go`
- `api/main.go`
