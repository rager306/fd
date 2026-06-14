---
id: T01
parent: S01
milestone: M042-fjf2en
key_files:
  - documents/te-perf-snapshot-m042-s01.md
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-06-14T08:54:27.086Z
blocker_discovered: false
---

# T01: Captured direct TEI telemetry snapshot showing batch-size-sensitive queue_time growth for batch 1, 8, and 32.

**Captured direct TEI telemetry snapshot showing batch-size-sensitive queue_time growth for batch 1, 8, and 32.**

## What Happened

Collected live TEI runtime limits from `http://localhost:30080/info`, then ran 10 direct TEI `/v1/embeddings` requests for each batch size 1, 8, and 32. Parsed `fd_tei` server logs since the run timestamp for `total_time`, `tokenization_time`, `queue_time`, and `inference_time`. Wrote the intermediate snapshot artifact to `documents/te-perf-snapshot-m042-s01.md`. The snapshot confirms that fd cache-hot M041 baseline timings do not explain cold TEI latency; direct TEI batch=32 has total p50 about 5.34s with queue p50 about 2.43s even for sequential single-client requests.

## Verification

Verified `documents/te-perf-snapshot-m042-s01.md` exists, is 5634 bytes, references the raw `.gsd/exec` evidence, includes runtime limits (`max_concurrent_requests`, `max_batch_requests`, `max_client_batch_size`), and contains at least 10 rows for each batch group 1, 8, and 32.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `python3 - <<'PY'
from pathlib import Path
p=Path('documents/te-perf-snapshot-m042-s01.md')
text=p.read_text()
rows=[line for line in text.splitlines() if line.startswith('| ') and line.count('|')>=8]
for batch in ['| 1 |','| 8 |','| 32 |']:
    print(batch.strip(), sum(1 for line in rows if line.startswith(batch)))
print('bytes', p.stat().st_size)
print('has_runtime_limits', all(s in text for s in ['max_concurrent_requests', 'max_batch_requests', 'max_client_batch_size']))
print('has_evidence_ref', '.gsd/exec/177ccc56-1ecc-434d-b12b-a0044ab13c20.stdout' in text)
PY` | 0 | ✅ pass: snapshot has required size, runtime limits, evidence link, and >=10 rows per batch group | 120000ms |

## Deviations

The plan path says `documents/` although the repository also has `docs/`; used the planned `documents/` path to satisfy the GSD artifact contract. The baseline file was fd cache-hot rather than TEI cold telemetry, so fresh direct TEI measurements were used for the required server-side timing table.

## Known Issues

This task proves batch-size-sensitive queue_time but does not yet distinguish TEI backend serialization from scheduler/metric semantics; that is T02/T03 work.

## Files Created/Modified

- `documents/te-perf-snapshot-m042-s01.md`
