# S01 — Research

**Date:** 2026-05-19

## Summary

S01 defines the artifact contract for the M011 opt-in ONNX prototype. The project now has a tracked manifest at `docs/onnx-artifacts/user-bge-m3-dense-fp32.json` that describes the local M010 FP32 dense ONNX candidate without committing the 1.43GB ONNX file. The actual binary remains under ignored runtime storage: `.gsd/runtime/onnx/m010-s03/user-bge-m3-dense.onnx`.

The manifest records model identity, source model revision and hashes, ONNX artifact path/size/SHA256, export dependency pins, runtime input/output metadata, validation artifacts, and a failure contract. It explicitly marks the artifact as `prototype_only`, `production_default=false`, and `git_tracked=false`. This preserves the D006 boundary: ONNX is locally feasible but TEI remains production/default until separate adapter, performance, artifact distribution, and quality gates pass.

S02/S03 should treat this manifest as the source of truth for local prototype validation. Runtime code must validate existence, checksum, and expected metadata before ONNX Runtime load. Missing files and checksum mismatches are actionable configuration/artifact failures, not silent fallback conditions for benchmark evidence.

## Recommendation

Proceed to S02 with a small manifest validation/config seam before attempting any ONNX runtime wiring. The safe validation order is:

1. Read manifest JSON.
2. Resolve configured/local artifact path.
3. Check file exists.
4. Check file size and SHA256 match manifest.
5. Later, after ONNX Runtime is introduced, check input/output metadata matches manifest.
6. Only then attempt inference.

Do not commit ONNX binaries. Do not make ONNX default. Do not silently fall back to TEI when ONNX is explicitly requested for benchmark evidence; fail fast with an actionable message instead.

## Implementation Landscape

### Key Files

- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json` — tracked artifact manifest for the M010 local ONNX candidate.
- `.gsd/runtime/onnx/m010-s03/user-bge-m3-dense.onnx` — ignored local ONNX artifact; referenced by manifest, not tracked.
- `.gsd/runtime/onnx/m010-s03/export-metadata.json` — ignored local export/load metadata used to populate the manifest.
- `benchmark-results/fd-onnx-fp32-m010-s03.txt` — tracked evidence that the artifact matched TEI baseline on fixed probes.
- `README.md` — currently states ONNX is future measured optimization; S04 may update it with opt-in artifact notes if runtime prototype lands.

### Build Order

1. In S02, add manifest parsing/checksum validation without enabling ONNX by default.
2. Add tests for valid manifest, missing file, checksum mismatch, invalid dimensions/output metadata, and TEI default unaffected behavior.
3. In S03, use the validated manifest to load ONNX only when explicit opt-in config is set.
4. In S04, benchmark and document the result.

### Verification Approach

- Manifest parses as JSON.
- If the local ONNX file exists, its size and SHA256 match manifest.
- `git status --ignored` shows `.gsd/runtime/` ignored and does not stage the ONNX binary.
- Runtime validation tests should assert error categories/messages for missing artifact and checksum mismatch.
- Existing Go tests/lint must remain green.

## Don't Hand-Roll

| Problem | Existing Solution | Why Use It |
|---------|------------------|------------|
| Artifact identity | SHA256 and file size in manifest | Simple, deterministic, language-agnostic validation. |
| Runtime metadata | ONNX Runtime session input/output metadata | Confirms loaded artifact matches expected dense contract. |
| Dense comparator | Existing M010 comparator artifacts | Prevents subjective equivalence claims and raw text leakage. |

## Constraints

- ONNX artifact is ~1.43GB and must stay outside git.
- Manifest references local development storage for now; production needs external artifact storage/download workflow.
- TEI remains default and must continue working without any ONNX artifact present.
- ONNX opt-in should fail fast when artifact validation fails; silent fallback would corrupt benchmark evidence.
- Raw embedding inputs must never be logged during artifact validation or load.

## Common Pitfalls

- **Loading before checksum validation** — validate file identity before ONNX Runtime opens it.
- **Silent fallback on explicit ONNX request** — fallback hides artifact/config mistakes and invalidates performance comparisons.
- **Committing large artifacts** — only the JSON manifest is tracked; `.gsd/runtime/` remains ignored.
- **Treating manifest as production distribution** — it is a local prototype contract, not a download mechanism.

## Open Risks

- Future production artifact storage is unresolved.
- Go ONNX Runtime binding may impose cgo/shared-library requirements that change deployment complexity.
- The manifest schema may need versioning if sparse/ColBERT/provider variants are ever introduced, but those are out of scope for M011.

## Sources

- M010 final recommendation and ONNX artifact evidence (source: `.gsd/milestones/M010-84qfzu/slices/S04/S04-RESEARCH.md`).
- Local export metadata and ONNX hash (source: `.gsd/runtime/onnx/m010-s03/export-metadata.json`).
- Tracked manifest created by this slice (source: `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`).
- ONNX vs TEI comparison evidence (source: `benchmark-results/fd-onnx-fp32-m010-s03.txt`).
