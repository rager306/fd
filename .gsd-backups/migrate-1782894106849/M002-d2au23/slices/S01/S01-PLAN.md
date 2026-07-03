# S01: Compose and git hygiene

**Goal:** Clean up Compose metadata and gitignore local runtime noise safely.
**Demo:** `docker compose config` is clean of obsolete version warning, and runtime GSD files are ignored while durable milestone artifacts remain trackable.

## Must-Haves

- Remove obsolete top-level Compose `version`.
- Ignore `.bg-shell/` and local GSD runtime/audit/journal/exec/graph/state files.
- Do not ignore durable `.gsd/milestones/` or `.gsd/gsd.db`.
- `docker compose config` emits no obsolete `version` warning.

## Proof Level

- This slice proves: config validation and git status check

## Integration Closure

No application behavior changes; only config metadata and ignore rules.

## Verification

- Cleaner git status makes real project changes easier to detect.

## Tasks

- [x] **T01: Assess hygiene cleanup scope** `est:small`
  Assess blast radius for config-only changes to docker-compose.yaml and .gitignore. Verify which GSD files are durable versus runtime noise before editing ignore rules.
  - Files: `docker-compose.yaml`, `.gitignore`, `.gsd/milestones/`, ` .gsd/gsd.db`
  - Verify: No code changes; document direct file scope.

- [x] **T02: Apply compose and gitignore cleanup** `est:small`
  Remove obsolete Compose version field and add narrow .gitignore entries for runtime-only GSD/bg-shell artifacts while preserving durable GSD artifacts.
  - Files: `docker-compose.yaml`, `.gitignore`
  - Verify: docker compose config >/tmp/fd-compose-clean.txt 2>/tmp/fd-compose-clean.err && ! grep -q 'obsolete' /tmp/fd-compose-clean.err

- [x] **T03: Verify hygiene cleanup** `est:small`
  Verify git status behavior and complete S01.
  - Verify: git status --short && git status --short --ignored

## Files Likely Touched

- docker-compose.yaml
- .gitignore
- .gsd/milestones/
-  .gsd/gsd.db
