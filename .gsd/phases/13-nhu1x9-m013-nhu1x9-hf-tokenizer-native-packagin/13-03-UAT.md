# S03: Parity correct ONNX tokenizer integration — UAT

**Milestone:** M013-nhu1x9
**Written:** 2026-05-20T03:56:23.414Z

# S03 UAT — Parity correct ONNX tokenizer integration

## Checks

- [x] Default untagged tests pass.
- [x] Default lint passes.
- [x] Tagged embed tests pass with native `libtokenizers.a`.
- [x] Tagged ONNX API starts locally.
- [x] Redis cache namespace is isolated with `EMBEDDING_CACHE_VERSION=m013-hf-tokenizer`.
- [x] TEI-vs-tagged-ONNX cosine passes all five fixed probes.
- [x] Raw probe text leak check passes.
- [x] No native binaries are tracked.
- [x] Tagged server cleaned up.

## UAT Result

Pass. The tagged HF tokenizer ONNX path is semantically equivalent on fixed probes and ready for a separate performance benchmark milestone.

