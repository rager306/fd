# S02: Contract validation and next gate — UAT

**Milestone:** M020-mvkq4d
**Written:** 2026-05-20T10:03:14.287Z

# S02 UAT — Contract validation and next gate

## Checks

- [x] D018 saved in the decision register.
- [x] Manifest JSON parses.
- [x] Required contract fields are present.
- [x] Evidence artifact links exist.
- [x] `production_default=false` preserved.
- [x] Tracked binary check reports zero `.onnx` and `libtokenizers.a` files.
- [x] Default Go tests pass.
- [x] Pinned GolangCI-Lint reports 0 issues.
- [x] Tagged HF tokenizer tests pass.
- [x] No background processes remain.
- [x] Port 18000 is clean.
- [x] GitNexus scope check is low.

## UAT Result

Pass. M020 can close as an artifact-contract milestone. Next work is Docker/CI packaging and artifact provisioning.

