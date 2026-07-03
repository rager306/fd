# S01: Same-host service contract — UAT

**Milestone:** M040-pbp9z1
**Written:** 2026-05-22T05:35:34.077Z

# S01 UAT: Same-host embedding service contract

## UAT Type
Operator/documentation contract review with implementation-backed health metadata verification.

## Preconditions
1. Repository checkout is at the completed S01 state.
2. Go tests can run from the `api` module directory.
3. The reviewer has access to `README.md`, `docs/same-host-embedding-service-contract.md`, and the S01 task summaries under `.gsd/milestones/M040-pbp9z1/slices/S01/tasks/`.

## Steps
1. Open `README.md` and confirm it links to `docs/same-host-embedding-service-contract.md` as the canonical same-host embedding service contract.
2. Open `docs/same-host-embedding-service-contract.md` and verify it documents `/health`, `/v1/embeddings`, and `/embeddings/batch` request/response expectations.
3. Confirm the contract explains dimensions, encoding caveats, status/error behavior, timeout/retry guidance, cache namespace guidance, runtime/env requirements, no-silent-fallback rules, readiness limitations, and non-goals.
4. Confirm `/health` runtime metadata is described as safe operational metadata, not as proof of live inference readiness.
5. Confirm `/v1/embeddings` request `model` is documented as compatibility metadata, while response `model` and `/health.runtime.model` are authoritative.
6. From the repository root, run `cd api && go test ./... -short`.
7. Run the leak audit command: `rg -n "signed|token=|X-Amz|BEGIN|PRIVATE|юридическая справка" docs README.md benchmark-results .gsd/milestones/M040-pbp9z1/slices/S01` and review any matches as policy/history references rather than current secret or raw corpus leaks.
8. Run a focused high-risk check for actual signed URL token parameters, private key blocks, and the prohibited raw sample in current public docs/artifacts.

## Expected Outcomes
- README points readers to the canonical contract without duplicating the full operating contract.
- The contract gives neighboring same-host services enough information to call fd and interpret readiness/runtime metadata safely.
- `/health` exposes TEI/default runtime identity through safe fields while preserving ONNX non-leaky metadata behavior.
- The model-field compatibility behavior is explicit and tested/documented.
- Go short tests pass from the `api` module.
- No current public docs or benchmark artifacts contain signed URLs, token query parameters, private key blocks, or the prohibited raw legal sample text.

## Edge Cases
- A client sends a non-empty `/v1/embeddings.model` value that differs from configured runtime model: it must treat that field as compatibility metadata and rely on response model plus `/health.runtime.model` instead.
- `/health` is healthy but no embedding smoke request has been run: clients must not assume full TEI/ONNX inference readiness solely from `/health`.
- Redis cache namespace changes between runtime comparisons: operators must isolate namespaces or flush deliberately to avoid cache contamination.

## Not Proven By This UAT
- Docker restart/cache behavior for packaged ONNX deployments; that is S02 scope.
- Alternative legal-domain model quality; that is S03 scope.
- Final TEI-vs-ONNX runtime recommendation; that is S04 scope.
- Live production throughput or latency benchmarking.
