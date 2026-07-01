# S01: Provisioning URL and archive hardening — UAT

**Milestone:** M029-4nh2ca
**Written:** 2026-05-21T04:32:56.555Z

# S01 UAT — Provisioning URL and archive hardening

## Checks

- [x] Private/localhost remote URLs are blocked by default.
- [x] HTTPS-only remote URL policy exists.
- [x] Optional allowed-host policy exists.
- [x] Remote downloads enforce byte caps.
- [x] Redirects are disabled.
- [x] Native tokenizer archive member must be a regular file.
- [x] Archive member size is checked before copy.
- [x] Provisioning dry-run/missing-source/verifier checks pass.
- [x] Default Go/lint/actionlint/Docker and hygiene checks pass.

## Result

Pass. M028 MEDIUM provisioning risks are remediated in code pending S02 docs/outcome closure.

