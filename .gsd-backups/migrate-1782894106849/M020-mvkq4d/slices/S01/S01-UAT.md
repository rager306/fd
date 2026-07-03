# S01: ONNX 1024 runtime contract — UAT

**Milestone:** M020-mvkq4d
**Written:** 2026-05-20T10:01:05.423Z

# S01 UAT — ONNX 1024 runtime contract

## Checks

- [x] Manifest JSON parses.
- [x] `export.sequence_length` remains `128`.
- [x] `runtime.validated_max_sequence_length` is `1024`.
- [x] M018 legal quality evidence is linked.
- [x] M019 performance evidence is linked.
- [x] `production_default` remains `false`.
- [x] Future gates include Docker/CI, packaged reruns, and TEI default preservation.
- [x] Tracked binary check reports zero `.onnx` and `libtokenizers.a` files.

## UAT Result

Pass. S02 can record the contract decision and close the metadata milestone.

