---
id: T01
parent: S03
milestone: M014-vjfs9f
key_files:
  - docs/onnx-artifacts/hf-tokenizers-linux-amd64.json
  - docs/onnx-artifacts/user-bge-m3-dense-fp32.json
  - benchmark-results/fd-go-onnx-hf-tokenizer-m013-s03.txt
key_decisions:
  - Use isolated Redis namespace `m014-onnx-hf-tokenizer` for tagged ONNX benchmark.
  - Use tagged build env `CGO_LDFLAGS=-L../.gsd/runtime/tokenizers/linux-amd64` and `go run -tags hf_tokenizers .` for S03 server.
  - Use ONNX Runtime library sha256 `50775d390eb55e7abd9f6d734da103a04f0e5342ef0a76b1c6ec795544439295`.
duration: 
verification_result: passed
completed_at: 2026-05-20T04:22:34.380Z
blocker_discovered: false
---

# T01: Preflight confirmed tagged ONNX artifacts, checksums, and tagged tests are ready for the benchmark.

**Preflight confirmed tagged ONNX artifacts, checksums, and tagged tests are ready for the benchmark.**

## What Happened

Validated required local artifacts for the tagged ONNX benchmark: ONNX Runtime shared library, ONNX model, native `libtokenizers.a`, tokenizer JSON, tracked manifests, and M013 cosine correctness artifact. The tagged `hf_tokenizers` Go test still passes. The benchmark will use isolated namespace `m014-onnx-hf-tokenizer` and snapshot v2 runtime metadata.

## Verification

Artifact checksum checks and tagged Go tests passed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `sha256sum ONNX Runtime, ONNX model, libtokenizers.a, tokenizer.json, manifests, M013 cosine artifact` | 0 | ✅ pass — all required files exist and checksums recorded | 0ms |
| 2 | `cd api && CGO_LDFLAGS='-L../.gsd/runtime/tokenizers/linux-amd64' go test -tags hf_tokenizers ./embed -count=1` | 0 | ✅ pass — ok fd-api/embed 0.124s | 5700ms |

## Deviations

None.

## Known Issues

Native tokenizer manifest still references upstream `latest`; this remains a known production packaging blocker but does not block local benchmark evidence.

## Files Created/Modified

- `docs/onnx-artifacts/hf-tokenizers-linux-amd64.json`
- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
- `benchmark-results/fd-go-onnx-hf-tokenizer-m013-s03.txt`
