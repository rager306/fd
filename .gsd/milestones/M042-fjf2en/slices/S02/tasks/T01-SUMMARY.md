---
id: T01
parent: S02
milestone: M042-fjf2en
key_files:
  - documents/onnx-deactivation-inventory-m042.md
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-06-14T10:23:58.267Z
blocker_discovered: false
---

# T01: Mapped active ONNX source/config/docs/tooling surfaces and defined the TEI-only removal boundary.

**Mapped active ONNX source/config/docs/tooling surfaces and defined the TEI-only removal boundary.**

## What Happened

Created `documents/onnx-deactivation-inventory-m042.md`. The inventory distinguishes active product surfaces from historical research artifacts. It recommends removing ONNX from runtime startup config, tests, build/dependency surfaces, operator docs, compose comments, and active CI/tooling references, while preserving historical benchmark/GSD artifacts. It notes that the normal binary may already exclude build-tagged ONNX, so the main gain is reducing product complexity and operator ambiguity rather than promising a dramatic binary-size win.

## Verification

Verified the inventory exists, is 5656 bytes, includes an active surfaces/decisions table, covers runtime config, embedder implementation, Docker/CI, docs/operator contract, historical artifacts, and explicit remove/keep decisions.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `python3 - <<'PY'
from pathlib import Path
p=Path('documents/onnx-deactivation-inventory-m042.md')
text=p.read_text()
checks={
 'exists': p.exists(),
 'bytes': p.stat().st_size > 3000,
 'active_surfaces': '## Active surfaces and decisions' in text,
 'runtime_config': 'Runtime startup config' in text,
 'embedder': 'ONNX embedder implementation' in text,
 'docker_ci': 'Docker ONNX image' in text,
 'docs': 'Docs/operator contract' in text,
 'keep_history': 'Historical benchmark artifacts' in text,
 'remove_keep': 'Decision' in text and 'Remove/disable' in text and 'Keep' in text,
}
print('bytes', p.stat().st_size)
for k,v in checks.items(): print(k, v)
PY` | 0 | ✅ pass: inventory contains required active surface and keep/remove decisions | 120000ms |

## Deviations

None.

## Known Issues

Actual code/doc removal remains for T02-T04. Need verify whether deleting build-tagged ONNX files affects shared tokenizer or manifest tests.

## Files Created/Modified

- `documents/onnx-deactivation-inventory-m042.md`
