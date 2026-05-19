# S01: Benchmark README update — UAT

**Milestone:** M005-pyn4x7
**Written:** 2026-05-19T10:41:43.320Z

# UAT: S01 Benchmark README update

## Verification performed

- `docker compose config >/tmp/fd-m005-s01-compose-config.out` — passed.
- `uv run --python 3.13 python -m py_compile benchmark.py` — passed.
- README snippet check confirmed:
  - `uv run --python 3.13 --with requests --with redis python benchmark.py`
  - `127.0.0.1:6379`
  - `FLUSHALL`
  - `docker compose restart api`
  - shared/production warning
- `gitnexus_detect_changes(repo=fd, scope=all)` — low risk, no affected processes.

## Result

README benchmark documentation now matches the validated local workflow and warns about Redis/API side effects.

