---
id: T01
parent: S04
milestone: M040-pbp9z1
key_files:
  - tools/verify_m040_s04_recommendation.py
key_decisions:
  - Verifier source files are not leak-scanned for prohibited fixture literals, but they are checked for expected verifier markers; evidence artifacts and the final recommendation artifact remain leak-scanned fail-closed.
  - Hosted-CI readiness-gate detection is line-oriented and only permits explicit negation of the readiness-gate phrase itself, avoiding false negatives from unrelated 'not required' language later on the same line.
duration: 
verification_result: passed
completed_at: 2026-05-22T08:28:46.058Z
blocker_discovered: false
---

# T01: Added a semantic verifier for the M040 S04 final runtime recommendation artifact, with inline positive and negative self-tests.

**Added a semantic verifier for the M040 S04 final runtime recommendation artifact, with inline positive and negative self-tests.**

## What Happened

Created `tools/verify_m040_s04_recommendation.py`, following the existing S02/S03 verifier style. The verifier checks required recommendation sections, validates source evidence paths and markers, enforces the explicit stance that packaged ONNX for `deepvk/USER-bge-m3` is only the preferred same-host performance runtime under the operating contract and explicit operator opt-in, keeps TEI as the current/default posture, requires Redis cache namespace isolation plus smoke readiness and ONNX artifact/tokenizer/runtime preflight language, defers alternative candidate model replacement fail-closed, rejects hosted-CI-as-readiness-gate language, and scans for secret/raw-text leak patterns. Self-test fixtures cover a passing recommendation plus negative cases for missing required sections, missing cache isolation, accepted candidate replacement, hosted CI gate language, secret-like content, and raw-text markers. I also exercised the real source-evidence validation path with the inline passing fixture to ensure the verifier fails closed for missing/malformed evidence inputs without requiring the final artifact to exist yet.

## Verification

Verified the new verifier compiles, its self-test passes, and its real evidence-file validation path passes using an inline valid recommendation fixture. During debugging, a false positive against verifier source files was corrected by excluding verifier scripts from leak scanning while still requiring their semantic markers and scanning evidence artifacts.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `python3 -m py_compile tools/verify_m040_s04_recommendation.py` | 0 | ✅ pass | 89ms |
| 2 | `python3 tools/verify_m040_s04_recommendation.py --self-test` | 0 | ✅ pass | 69ms |
| 3 | `python3 - <<'PY'
from pathlib import Path
import tempfile
from tools.verify_m040_s04_recommendation import DEFAULT_SOURCE_PATHS, sample_recommendation, validate_path
with tempfile.TemporaryDirectory() as tmpdir:
    artifact = Path(tmpdir) / 'M040-S04-RECOMMENDATION.md'
    artifact.write_text(sample_recommendation(), encoding='utf-8')
    validate_path(artifact, DEFAULT_SOURCE_PATHS)
print('real evidence validation path: PASS')
PY` | 0 | ✅ pass | 83ms |

## Deviations

Added an extra real-evidence validation exercise beyond the two required commands, to verify source artifact marker checks before the final recommendation artifact exists.

## Known Issues

GitNexus MCP tools were listed in project instructions but not exposed in the available function namespace for this session, so no GitNexus impact/detect_changes tool call could be made. No existing symbols were edited; the task created a new standalone verifier file.

## Files Created/Modified

- `tools/verify_m040_s04_recommendation.py`
