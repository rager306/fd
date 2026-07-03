---
id: T01
parent: S02
milestone: M045-d0e5xq
key_files:
  - documents/tei-startup-mitigation-m045.md
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-06-14T11:37:43.747Z
blocker_discovered: false
---

# T01: Inventoried TEI `/data` cache and confirmed USER-bge-m3 Candle files are present.

**Inventoried TEI `/data` cache and confirmed USER-bge-m3 Candle files are present.**

## What Happened

Performed symlink-aware read-only inventory of the TEI `/data` volume. The USER-bge-m3 snapshot contains config, pooling config, modules, tokenizer_config, tokenizer.json, sentence-transformers config, and model.safetensors. No ONNX files were present, which matches the TEI-only scope.

## Verification

`documents/tei-startup-mitigation-m045.md` records cached USER-bge-m3 files and required-file counts. The inventory did not read file contents or print secrets.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `docker exec fd_tei find/readlink/stat inventory under /data, no file contents` | 0 | ✅ pass: USER-bge-m3 cache appears complete for Candle/safetensors startup | 180000ms |

## Deviations

The first inventory pass only saw blob files because HF snapshots are symlink-based; reran a symlink-aware shell inventory.

## Known Issues

Cache completeness makes offline mode plausible but does not prove startup until S03 controlled restart.

## Files Created/Modified

- `documents/tei-startup-mitigation-m045.md`
