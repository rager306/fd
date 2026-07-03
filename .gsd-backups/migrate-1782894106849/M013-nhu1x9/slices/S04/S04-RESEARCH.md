# S04 — Research

**Date:** 2026-05-20

## Summary

M013 closes with a stronger result than expected: the tagged HF tokenizer ONNX path is now fixed-probe benchmark-ready. S01 created a native artifact contract for `libtokenizers.a`; S02 proved default build isolation and tagged tokenizer parity; S03 integrated the parity-correct tokenizer into the opt-in ONNX backend under `hf_tokenizers` and passed TEI-vs-ONNX cosine on all fixed probes with isolated Redis cache namespace.

This is still not production readiness. The native artifact remains local under `.gsd/runtime/tokenizers/linux-amd64/libtokenizers.a`, Docker/CI tagged build support is not yet implemented, the source URL should be pinned beyond `latest`, and validation is still only fixed-probe cosine rather than a larger Russian/legal retrieval corpus. TEI remains the production/default runtime.

The next milestone can now be a tagged ONNX performance benchmark, provided it uses the `hf_tokenizers` build tag, validated native artifact, isolated Redis namespace, sanitized benchmark config snapshots, and explicit comparison against the current TEI+Redis baseline. Any production switch remains out of scope until native packaging, CI, larger quality evaluation, and operational startup/memory checks pass.

## Recommendation

Proceed next to **tagged ONNX performance benchmarking**, not production deployment.

Required benchmark setup:

- Build/run API with `-tags hf_tokenizers`.
- Set `CGO_LDFLAGS=-L../.gsd/runtime/tokenizers/linux-amd64` from `api` or equivalent in Docker.
- Use `EMBEDDING_BACKEND=onnx` with the validated M010 ONNX artifact and M013 native tokenizer artifact.
- Use isolated Redis namespace, for example `EMBEDDING_CACHE_VERSION=m013-hf-tokenizer-benchmark`.
- Capture sanitized effective configuration, including build tags, native artifact checksum, ONNX artifact checksum, ORT library path/hash, and environment baseline.
- Compare TEI default vs tagged ONNX for cold, warm, batch, cache-hit, startup, and memory behavior.

Do not switch defaults. Do not benchmark the untagged ONNX path as equivalent. Do not publish or push remote changes without explicit confirmation.

## Implementation Landscape

### Key Files

- `docs/onnx-artifacts/hf-tokenizers-linux-amd64.json` — native tokenizer artifact manifest and checksum contract.
- `benchmark-results/fd-tokenizers-native-artifact-m013-s01.txt` — native artifact validation evidence.
- `api/embed/hf_tokenizer_native.go` — tagged HF native tokenizer wrapper.
- `api/embed/onnx_tokenizer_default.go` — default untagged `sugarme` tokenizer adapter.
- `api/embed/onnx_tokenizer_hf.go` — tagged tokenizer factory for ONNX embedder.
- `api/embed/onnx.go` — ONNX embedder tokenizer interface and batch encoding.
- `benchmark-results/fd-tokenizer-tagged-native-m013-s02.txt` — tagged tokenizer parity proof.
- `benchmark-results/fd-go-onnx-hf-tokenizer-m013-s03.txt` — tagged ONNX cosine pass.
- `benchmark.py` — likely next milestone entrypoint for performance benchmark extension.

### Build Order

1. Preserve M013 state as the benchmark-ready correctness baseline.
2. Extend benchmark config snapshot to record build tags/native artifact/ONNX runtime details if not already present.
3. Start TEI default stack and tagged ONNX stack separately with isolated Redis namespaces.
4. Run comparable latency/throughput/cache/startup/memory measurements.
5. Only after benchmark evidence, decide whether deeper ONNX packaging or provider optimization is worthwhile.

### Verification Approach

Final M013 closure requires:

- Default `cd api && go test ./... -short`.
- Default pinned GolangCI-Lint.
- Tagged `cd api && CGO_LDFLAGS='-L../.gsd/runtime/tokenizers/linux-amd64' go test -tags hf_tokenizers ./embed -count=1`.
- Default `/health` check.
- Artifact parser/leak/native-binary checks.
- GitNexus detect changes.
- No lingering tagged ONNX server.

## Constraints

- TEI remains default.
- ONNX remains opt-in.
- Native `libtokenizers.a` is not committed.
- No raw probe text in benchmark artifacts.
- No production switch before Docker/CI packaging and larger Russian/legal quality gates.

## Common Pitfalls

- **Benchmarking the wrong build** — performance work must use `-tags hf_tokenizers`, not the default `sugarme` ONNX tokenizer path.
- **Shared Redis masking** — keep ONNX benchmark namespaces isolated.
- **Treating fixed probes as full quality proof** — cosine pass unlocks benchmarking, not production rollout.
- **Forgetting native artifact provenance** — benchmark artifacts must record native checksum and build flags.

## Open Risks

- Tagged Docker/CI packaging is still absent.
- Native artifact source URL is currently `latest`; production reproducibility needs pinning.
- ONNX CPU performance may still be worse than TEI after correctness is fixed.
- Larger Russian/legal retrieval quality may reveal gaps not covered by fixed probes.

## Sources

- Native artifact validation: `benchmark-results/fd-tokenizers-native-artifact-m013-s01.txt`
- Tagged tokenizer parity: `benchmark-results/fd-tokenizer-tagged-native-m013-s02.txt`
- Tagged ONNX cosine: `benchmark-results/fd-go-onnx-hf-tokenizer-m013-s03.txt`
