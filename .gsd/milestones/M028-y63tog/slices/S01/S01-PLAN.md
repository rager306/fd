# S01: ONNX operational security review

**Goal:** Map attack surfaces and produce STRIDE/OWASP security findings for ONNX operational surfaces.
**Demo:** After this, there is a file:line security review of ONNX artifact/provisioning/startup surfaces with prioritized findings.

## Must-Haves

- Attack surfaces are mapped.
- Findings have file:line, exploit, severity, reachability, impact, remediation.
- Non-findings are listed.
- No code changes beyond GSD/report artifacts.
- Report avoids raw secrets and signed URLs.

## Proof Level

- This slice proves: Code-read report with file:line citations and verification that no code remediation occurred.

## Integration Closure

Provides a concrete remediation backlog for future GSD milestones without changing runtime behavior.

## Verification

- Identifies where startup errors/logs/health/provisioning can leak or be abused.

## Tasks

- [x] **T01: Map ONNX security attack surfaces** `est:medium`
  Map ONNX operational attack surfaces across Go startup/preflight, health metadata, provisioning helper, workflows, verifier, Docker packaging, and operations docs. Capture entry points, trust boundaries, and sinks with file:line references.
  - Files: `api/main.go`, `api/embed/onnx_manifest.go`, `api/handlers/health.go`, `tools/provision_onnx_artifacts.py`, `tools/verify_onnx_artifacts.py`, `.github/workflows/onnx-packaging.yml`, `tools/build_onnx_image.sh`, `docs/onnx-artifacts/OPERATIONS.md`
  - Verify: Mapped surfaces have file:line citations.

- [x] **T02: Write security findings report** `est:medium`
  Run STRIDE/OWASP pass on mapped surfaces and write the security review artifact with prioritized findings, non-findings, out-of-scope, and follow-up recommendations.
  - Files: `.gsd/milestones/M028-y63tog/slices/S01/S01-RESEARCH.md`
  - Verify: Report has severity/reachability/exploit/remediation and no secret/raw text leaks.

- [x] **T03: Verify read-only review scope** `est:small`
  Verify S01 is read-only with no code remediation, run artifact hygiene checks, complete S01.
  - Verify: Git diff excludes code remediation except GSD/report docs; leak checks pass.

## Files Likely Touched

- api/main.go
- api/embed/onnx_manifest.go
- api/handlers/health.go
- tools/provision_onnx_artifacts.py
- tools/verify_onnx_artifacts.py
- .github/workflows/onnx-packaging.yml
- tools/build_onnx_image.sh
- docs/onnx-artifacts/OPERATIONS.md
- .gsd/milestones/M028-y63tog/slices/S01/S01-RESEARCH.md
