---
id: T03
parent: S02
milestone: M010-84qfzu
key_files:
  - benchmark-results/fd-dense-comparator-m010-s02.txt
  - tools/compare_dense_embeddings.py
key_decisions:
  - Keep `benchmark-results/fd-dense-comparator-m010-s02.txt` as the TEI/API dense baseline artifact for S03.
  - Artifact verification checks required sections and confirms raw probe texts from script constants do not appear in the saved markdown.
duration: 
verification_result: passed
completed_at: 2026-05-19T18:37:27.997Z
blocker_discovered: false
---

# T03: Captured the sanitized TEI/API dense comparator baseline artifact for S03 ONNX comparison.

**Captured the sanitized TEI/API dense comparator baseline artifact for S03 ONNX comparison.**

## What Happened

Captured and validated the TEI dense comparator baseline artifact at `benchmark-results/fd-dense-comparator-m010-s02.txt`. The artifact records sanitized configuration, five probe labels with character counts, 1024-dimensional vector checks, finite-value and L2-normalization status, float32 vector hashes, pairwise cosine similarities, usage, and a PASS verdict. A parser verified required sections and confirmed the raw probe texts embedded in the script constants are not present in the artifact.

## Verification

Artifact validation passed: required sections/tokens exist, expected dimensions are recorded as 1024, raw probe text logging is false, verdict is PASS, and none of the five raw probe texts appear in the artifact.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `python3 - <<'PY'
from pathlib import Path
artifact = Path('benchmark-results/fd-dense-comparator-m010-s02.txt')
script = Path('tools/compare_dense_embeddings.py')
text = artifact.read_text(encoding='utf-8')
script_text = script.read_text(encoding='utf-8')
required = ['# Dense Embedding Comparator Baseline', '## Effective Configuration', '## Probe Summary', '## Pairwise Cosine Similarities', '## Verdict', 'PASS', '"expected_dimensions": 1024', '"raw_probe_texts_logged": false']
missing = [item for item in required if item not in text]
if missing:
    raise SystemExit(f'missing required artifact tokens: {missing}')
raw_texts = []
for line in script_text.splitlines():
    line = line.strip()
    if line.startswith('"text":'):
        raw = line.split(':', 1)[1].strip().strip(',').strip('"')
        raw_texts.append(raw)
leaked = [raw for raw in raw_texts if raw and raw in text]
if leaked:
    raise SystemExit(f'raw probe text leaked into artifact: {leaked[:2]}')
print('artifact_ok=true')
print(f'raw_probe_count_checked={len(raw_texts)}')
PY` | 0 | ✅ pass | 0ms |

## Deviations

None. T03 reused the artifact produced by the T02 implementation run and added explicit artifact validation.

## Known Issues

The artifact is a point-in-time TEI/API baseline. If model revision, tokenizer, dimensions, or runtime config change, S03 should regenerate it before comparison.

## Files Created/Modified

- `benchmark-results/fd-dense-comparator-m010-s02.txt`
- `tools/compare_dense_embeddings.py`
