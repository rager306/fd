---
milestone: M028-y63tog
slice: S01
type: security-review
status: complete
---

# Security Review — ONNX Operational Artifact and Startup Surfaces

## Summary

Reviewed the opt-in ONNX operational surfaces added through M025-M027: Go startup preflight, health metadata, provisioning helper, verifier, manual packaging workflow, Docker packaging script, and operations docs. Overall posture is acceptable for an opt-in/admin-operated experimental path: no artifact integrity bypass was found because model/native artifacts are hash-verified before use. The main risks are manual-workflow SSRF/unbounded downloads, unbounded archive extraction/copying, and local path disclosure in logs/errors.

TEI remains the default runtime. ONNX remains explicit opt-in experimental.

## Attack surface map

| Surface | Entry point | Reachability | Sinks |
|---|---|---|---|
| Go startup config | `EMBEDDING_BACKEND`, `ONNX_ARTIFACT_MANIFEST`, `ONNX_RUNTIME_LIBRARY`, `ONNX_TOKENIZER_PATH`, `ONNX_RUNTIME_SHA256`, `ONNX_PROVIDER`, `ONNX_MAX_SEQUENCE_LENGTH` in `api/main.go:149-199` | same-host process / deployment operator | filesystem stat/open/hash, startup logs, process exit |
| ONNX manifest paths | `artifact.local_path` read by `api/embed/onnx_manifest.go:133-153` and resolved by `api/embed/onnx_manifest.go:172-189` | trusted tracked manifest / local operator | filesystem stat/open/hash and error strings |
| Health metadata | `/health` response assembled in `api/handlers/health.go:35-42` | HTTP caller | safe runtime JSON metadata |
| Manual ONNX workflow | `workflow_dispatch` inputs in `.github/workflows/onnx-packaging.yml:4-25` | GitHub actor with workflow dispatch rights | shell env, provisioning helper, docker build |
| Provisioning downloads | explicit URL/file args in `tools/provision_onnx_artifacts.py:29-39`; URL download in `tools/provision_onnx_artifacts.py:112-117` | workflow operator / local operator | outbound HTTP(S), temp files, artifact destinations |
| Native tokenizer archive extraction | archive member copy in `tools/provision_onnx_artifacts.py:125-140` | workflow operator / local operator | tar parsing, destination write |
| Verifier | manifest-local paths in `tools/verify_onnx_artifacts.py:77-87`, hash verification in `tools/verify_onnx_artifacts.py:106-147` | CI/local operator | filesystem stat/hash, JSON output |
| Docker packaging | local artifact copies in `tools/build_onnx_image.sh:20-47` | local/CI operator | filesystem copy into ignored Docker context |

## Findings

### MEDIUM-1: Manual artifact URLs allow arbitrary outbound HTTP(S) fetches and unbounded downloads

**Location:** `.github/workflows/onnx-packaging.yml:4-25`, `.github/workflows/onnx-packaging.yml:74-90`, `tools/provision_onnx_artifacts.py:108-117`

**Category:** Information disclosure / Denial of service / SSRF-adjacent (STRIDE); OWASP SSRF pattern.

**Exploit scenario:** A GitHub actor who can manually dispatch the ONNX packaging workflow supplies `onnx_source_url`, `native_tokenizer_source_url`, or `onnx_runtime_source_url` pointing to an attacker-controlled or internal/probing URL. The workflow passes the values to `tools/provision_onnx_artifacts.py` (`.github/workflows/onnx-packaging.yml:74-90`), which accepts any `http` or `https` URL (`tools/provision_onnx_artifacts.py:108-117`) and streams the response to disk without host allowlisting, private-address blocking, redirect policy, `Content-Length` enforcement, or a byte cap. Artifact sha checks prevent silent artifact substitution later, but they do not prevent runner network probing, runner IP disclosure to an attacker-controlled server, or disk/time exhaustion before checksum failure.

**Reachability:** Authenticated GitHub user with permission to run `workflow_dispatch`, or local operator running the helper.

**Business impact:** CI runner abuse, service disruption, unintended network egress, noisy security/audit events. Artifact integrity is still protected by sha verification.

**Remediation:** Require immutable non-secret artifact references from approved hosts or GitHub cache/artifact IDs instead of arbitrary URLs. If URLs remain supported, add URL policy checks before `urlopen`: allowlist expected domains, block localhost/private/link-local ranges after DNS resolution and after redirects, cap maximum response bytes, require and validate `Content-Length` when available, and fail on excessive redirects.

### MEDIUM-2: Provisioning can copy unbounded archive members before size/hash verification

**Location:** `tools/provision_onnx_artifacts.py:125-140`, `tools/provision_onnx_artifacts.py:143-152`

**Category:** Denial of service (STRIDE).

**Exploit scenario:** A workflow/local operator supplies a native tokenizer archive whose member basename is `libtokenizers.a`, but the member is extremely large or compressed to consume significant disk/CPU during extraction. The helper selects the member by basename and copies it to the destination at `tools/provision_onnx_artifacts.py:125-140`; only after the full copy does `verify_destination` check size and sha at `tools/provision_onnx_artifacts.py:143-152`. The malicious artifact cannot pass verification unless it matches the expected sha, but it can still exhaust temporary storage or runner time before failing.

**Reachability:** Authenticated GitHub workflow dispatcher or local operator.

**Business impact:** CI/service disruption and wasted runner resources.

**Remediation:** Before copying archive members, inspect `TarInfo.size` and fail if it exceeds the manifest expected size or a hard maximum. During copy, stream with a byte counter and abort after the expected size plus a small tolerance. Consider rejecting non-regular tar members explicitly.

### LOW-3: Startup and verifier errors disclose local filesystem paths in shared logs

**Location:** `api/main.go:47-72`, `api/main.go:77-90`, `api/main.go:198`, `api/main.go:221-224`, `api/embed/onnx_manifest.go:137-153`, `tools/verify_onnx_artifacts.py:117-138`, `tools/build_onnx_image.sh:20-24`

**Category:** Information disclosure (STRIDE).

**Exploit scenario:** An operator misconfigures ONNX paths or checksum values in a shared CI/deployment environment. Startup and verifier errors include full configured paths or resolved artifact paths, for example path-bearing Go errors in `api/main.go:47-72` and manifest artifact errors in `api/embed/onnx_manifest.go:137-153`. `tools/build_onnx_image.sh:20-24` also prints missing tokenizer/runtime paths. These are not secret values by themselves, but they can expose host layout, cache locations, usernames, local artifact staging paths, and operational conventions in logs available to a broad CI/deployment audience.

**Reachability:** Same-host/deployment operator; any user with access to CI/container logs.

**Business impact:** Minor information disclosure that can assist follow-on attacks or leak environment topology. No raw embedding input text or signed URL leak was observed in the reviewed code path.

**Remediation:** Split diagnostics into safe and debug modes. Default logs should use env key, artifact_id, basename or repo-relative path, and error class; full absolute paths should be enabled only under debug/admin mode. Keep existing sha mismatch details, but avoid absolute host paths where possible.

### LOW-4: Manifest path resolution accepts absolute/repo-external paths for local tools

**Location:** `api/embed/onnx_manifest.go:172-189`, `tools/verify_onnx_artifacts.py:77-87`, `tools/provision_onnx_artifacts.py:55-58`, `tools/provision_onnx_artifacts.py:68-84`

**Category:** Tampering / Information disclosure (STRIDE); path traversal adjacent.

**Exploit scenario:** The tracked manifests are currently trusted, but the tooling accepts absolute manifest artifact paths and repo-external paths. `resolveONNXArtifactPath` returns absolute artifact paths unchanged (`api/embed/onnx_manifest.go:172-177`), and the verifier similarly honors absolute `artifact.local_path` (`tools/verify_onnx_artifacts.py:77-87`). A local operator or workflow run against a malicious branch could make the verifier hash/read unintended host files or make provisioning write artifacts outside the expected ignored artifact directories. The attacker needs control over the manifest or workflow branch, so this is not remotely exploitable by ordinary API users.

**Reachability:** Local operator, malicious branch selected for manual workflow, or committer modifying manifests.

**Business impact:** Local/CI filesystem disclosure or artifact placement outside expected boundaries; possible CI confusion. Production runtime remains protected by opt-in startup and artifact checks.

**Remediation:** Enforce path policy for manifests and provisioning destinations: reject absolute paths and `..` traversal for artifact destinations, or restrict them to explicit allowlisted roots such as `.gsd/runtime/onnx/`, `.gsd/runtime/tokenizers/`, `.gsd/runtime/onnxruntime/`, and `tei-models/`. Keep an explicit override only for trusted local development if needed.

## Non-findings considered

- **Shell injection through workflow inputs:** Not found. The workflow stores inputs in environment variables and passes them as quoted array entries (`.github/workflows/onnx-packaging.yml:74-90`), so shell metacharacters in URLs are not directly evaluated.
- **Health endpoint path disclosure:** Not found. `/health.runtime` is built from `RuntimeHealth` in `api/handlers/health.go:35-42` and exposes metadata fields only; M026/M027 fields do not include manifest, tokenizer, runtime library paths, signed URLs, or raw input text.
- **Artifact integrity bypass after download:** Not found for ONNX/native tokenizer artifacts. Provisioning verifies expected size/sha after materialization (`tools/provision_onnx_artifacts.py:143-152`), and verifier repeats strict hash checks (`tools/verify_onnx_artifacts.py:132-147`).
- **Tar path traversal via archive member name:** Not found in the current extraction path. The helper does not call `extractall`; it opens the selected member and writes to a fixed destination (`tools/provision_onnx_artifacts.py:125-140`). The remaining concern is unbounded size/member type, captured in MEDIUM-2.
- **Production default bypass:** Not found. Manifest validation rejects `production_default=true` in Go startup (`api/embed/onnx_manifest.go:94-96`) and verifier requires `production_default=false` (`tools/verify_onnx_artifacts.py:67-75`, `tools/verify_onnx_artifacts.py:106-110`).

## Out of scope

- Remediation patches. This review is intentionally read-only.
- Remote hosted workflow execution, because no push/workflow dispatch was authorized.
- Production rollout security, because ONNX remains opt-in experimental.
- Third-party dependency vulnerability scanning.
- Runtime memory safety of ONNX Runtime, native tokenizer library, or CGO bindings.

## Recommended follow-up order

1. Add URL/download safety to `tools/provision_onnx_artifacts.py`: approved source policy, private-address blocking, redirect handling, and byte caps.
2. Add archive member size/type checks before extraction/copying.
3. Sanitize startup/verifier/build-script path output by default while preserving debug detail for local operators.
4. Enforce manifest artifact path roots for CI/provisioning, or document an explicit trusted-local override.

## Review integrity

No code remediation was performed as part of this review. Findings are intended to seed a follow-up remediation milestone.
