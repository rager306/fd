---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T02: Apply compose and gitignore cleanup

Remove obsolete Compose version field and add narrow .gitignore entries for runtime-only GSD/bg-shell artifacts while preserving durable GSD artifacts.

## Inputs

- `S01 T01`

## Expected Output

- `docker-compose.yaml`
- `.gitignore`

## Verification

docker compose config >/tmp/fd-compose-clean.txt 2>/tmp/fd-compose-clean.err && ! grep -q 'obsolete' /tmp/fd-compose-clean.err

## Observability Impact

Makes git status highlight actionable changes only.
