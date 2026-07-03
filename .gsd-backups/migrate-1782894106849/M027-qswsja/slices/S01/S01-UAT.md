# S01: Startup artifact and provider preflight — UAT

**Milestone:** M027-qswsja
**Written:** 2026-05-20T12:40:28.852Z

# S01 UAT — Startup artifact and provider preflight

## Checks

- [x] Manifest exposes tokenizer JSON source metadata.
- [x] Startup validates tokenizer JSON size/sha from manifest.
- [x] Startup validates ONNX Runtime sha when `ONNX_RUNTIME_SHA256` is configured.
- [x] Startup rejects unsupported provider values.
- [x] Health metadata includes safe provider/tokenizer/runtime verification flags.
- [x] Default Go tests, lint, tagged tests, Docker, scripts, hygiene, and cleanup pass.

## Result

Pass. ONNX startup has stronger preflight diagnostics while TEI/default remains unchanged.

