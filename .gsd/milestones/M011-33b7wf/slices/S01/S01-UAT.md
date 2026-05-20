# S01: Artifact manifest and checksum contract — UAT

**Milestone:** M011-33b7wf
**Written:** 2026-05-19T18:59:05.595Z

# S01 UAT — Artifact manifest and checksum contract

## Checks

- [x] M010 export metadata inspected.
- [x] Tracked manifest exists at `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`.
- [x] Manifest JSON parses.
- [x] Manifest size/SHA256 match local artifact when present.
- [x] ONNX binary remains under ignored `.gsd/runtime/`.
- [x] S01 research documents missing artifact and checksum mismatch behavior.
- [x] Production default remains TEI; no runtime code changed.

## UAT Result

Pass. S02 can now implement manifest/config validation from a tracked contract instead of hardcoding local artifact assumptions.

