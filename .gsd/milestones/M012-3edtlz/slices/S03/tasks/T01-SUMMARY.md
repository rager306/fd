---
id: T01
parent: S03
milestone: M012-3edtlz
key_files: []
key_decisions:
  - `github.com/daulet/tokenizers` is feasible locally with the prebuilt `libtokenizers.linux-amd64.tar.gz` static library and `CGO_LDFLAGS=-L<libdir>`.
  - The isolated probe loaded the exact local `tokenizer.json` and produced IDs/mask matching the known Hugging Face labor-law baseline probe.
duration: 
verification_result: passed
completed_at: 2026-05-20T02:15:08.467Z
blocker_discovered: false
---

# T01: Proved `daulet/tokenizers` can run locally with prebuilt HF Rust tokenizers and matches one known HF baseline probe.

**Proved `daulet/tokenizers` can run locally with prebuilt HF Rust tokenizers and matches one known HF baseline probe.**

## What Happened

Created a temporary module under `/tmp/fd-daulet-tokenizers-probe`, added `github.com/daulet/tokenizers@latest`, downloaded the prebuilt linux-amd64 `libtokenizers.a` release asset, linked with `CGO_LDFLAGS`, loaded `/root/fd/tei-models/deepvk--USER-bge-m3/tokenizer.json`, and encoded the labor-law probe. The IDs and attention mask match the Hugging Face baseline sequence for that probe, showing the binding is a viable candidate for full S03 comparison.

## Verification

Temp-module probe exited 0 and printed token IDs/mask matching the S01 HF baseline for the labor-law probe.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `temp module: go get github.com/daulet/tokenizers@latest; download libtokenizers.linux-amd64.tar.gz; CGO_LDFLAGS=-L<libdir> go run .` | 0 | ✅ pass — loaded tokenizer.json and emitted HF-matching IDs/mask for one probe | 8200ms |

## Deviations

None.

## Known Issues

This proves local feasibility only. Production integration still needs a committed/downloaded native library strategy for Docker/CI and must not assume `/tmp` assets exist.

## Files Created/Modified

None.
