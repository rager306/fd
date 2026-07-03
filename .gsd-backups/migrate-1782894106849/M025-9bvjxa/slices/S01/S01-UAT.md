# S01: Artifact provisioning contract — UAT

**Milestone:** M025-9bvjxa
**Written:** 2026-05-20T11:40:46.082Z

# S01 UAT — Artifact provisioning contract

## Checks

- [x] Provisioning contract exists.
- [x] README links provisioning contract.
- [x] Provisioning helper compiles.
- [x] Dry-run reports missing source blockers as JSON.
- [x] Non-dry-run without sources fails fast.
- [x] Strict local verifier still passes.
- [x] Binary hygiene passes.

## Result

Pass. Artifact provisioning is explicit and checksum-based, with missing external sources represented as blockers.

