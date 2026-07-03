---
id: T02
parent: S01
milestone: M042-fjf2en
key_files:
  - benchmark-results/te-concurrency-profile-m042-s01.md
  - benchmark-results/te-concurrency-profile-m042-s01-run.txt
  - tools/profile_tei_concurrency.sh
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-06-14T10:15:16.373Z
blocker_discovered: false
---

# T02: Documented the attempted TEI concurrency profile and the stronger finding: TEI restart/recreate can spend ~48 minutes in backend startup before becoming ready.

**Documented the attempted TEI concurrency profile and the stronger finding: TEI restart/recreate can spend ~48 minutes in backend startup before becoming ready.**

## What Happened

Created `benchmark-results/te-concurrency-profile-m042-s01.md` from the failed profiler run, current Docker health/log evidence, and the T01 direct TEI snapshot. The original profiler attempted sequential batch=32, parallel batch=32, parallel batch=1, idle batch=32, and restart-then-batch32 scenarios, but failed at the restart health gate and later log metrics were lost after container recreation. The revised artifact does not claim missing metrics as success. It records that TEI eventually recovered to healthy, but the recreate path spent about 48 minutes from `Starting model backend` to `Ready`, including a delayed missing-ONNX ORT backend failure before Candle/safetensors warmup. This supports the user's TEI-first rescope and ONNX deferral.

## Verification

Verified the artifact exists, is 5987 bytes, includes all five attempted scenarios, records TEI recovered to healthy, records `Could not start ORT backend`, and explicitly states ONNX runtime implementation is deferred.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `python3 - <<'PY'
from pathlib import Path
p=Path('benchmark-results/te-concurrency-profile-m042-s01.md')
text=p.read_text()
print('bytes', p.stat().st_size)
for s in ['sequential_batch32_control','parallel_4_batch32','parallel_16_batch1','idle_30s_batch32','restart_then_batch32']:
    print(s, s in text)
for s in ['TEI eventually recovered to `healthy`','Could not start ORT backend','ONNX runtime implementation is deferred']:
    print(s[:40], s in text)
PY` | 0 | ✅ pass: revised T02 artifact contains required scenarios and TEI startup/ONNX deferral evidence | 120000ms |

## Deviations

The original T02 intended to collect complete concurrency metrics. Because the restart scenario became destructive and logs were lost after recreate, T02 was replanned to preserve evidence without repeating destructive restarts. Missing metrics are documented as limitations, not pass evidence.

## Known Issues

TEI cold restart/backend startup is operationally fragile and should be addressed in the RCA. Parallel concurrency metrics should be recollected only after TEI startup is stable and with non-destructive log capture.

## Files Created/Modified

- `benchmark-results/te-concurrency-profile-m042-s01.md`
- `benchmark-results/te-concurrency-profile-m042-s01-run.txt`
- `tools/profile_tei_concurrency.sh`
