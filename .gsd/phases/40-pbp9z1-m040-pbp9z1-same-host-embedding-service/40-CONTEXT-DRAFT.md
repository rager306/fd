# M040-pbp9z1 Context Draft

## Vision
Prepare `fd` as a local same-host embedding service for neighboring services, with excellent Russian/legal-domain quality and optimal speed on this host. ONNX is a candidate runtime, not the goal itself. Experiments outside that boundary are out of scope except a bounded quick check of alternative embedding models on the legal domain.

## Confirmed Scope
- Primary result: practical TEI-vs-ONNX runtime recommendation for same-host operation.
- Consumers: local HTTP clients on the same host using `/v1/embeddings`, batch endpoint, and `/health`.
- Model: `deepvk/USER-bge-m3` remains baseline; include only a bounded quick gate for 1-2 legal-compatible alternative candidates.
- Out of scope: hosted GitHub Actions proof, push/dispatch/upload, ONNX for its own sake, Rust rewrite, provider/INT8/NUMA experiments without same-host readiness need, embedded/library surface.

## Confirmed Architecture Direction
- Runtime recommendation basis: evidence envelope, not raw speed alone. Legal-domain quality, same-host performance, Redis/cache restart behavior, health/preflight clarity, and operational simplicity all count.
- Restart/cache proof: scripted Docker restart via `BENCHMARK_API_RESTART_COMMAND` or equivalent small local command, preserving benchmark comparability and proving packaged lifecycle.
- Alternative model quick gate: offline shortlist/profile of 1-2 candidates from known/legal-compatible sources. Stop unless a candidate clearly beats or challenges the current USER-bge-m3 envelope.

## Error Handling Strategy
- Fail fast on bad runtime/artifact/tokenizer/runtime-library configuration.
- `/health` is the primary safe readiness surface and must expose backend, model, dimensions, artifact/tokenizer/runtime verification, provider, cache namespace, sequence length, and production/default flag without secrets.
- No silent fallback between TEI and ONNX or between tokenizers/models within one service run.
- Redis namespace contamination is a correctness risk; comparisons must isolate namespaces or explicitly record cache flush side effects.
- Local clients should use bounded timeouts/retries only for transient service errors; configuration/model mismatch should not be blindly retried.
- Fallback/rollback is operational restart/reconfiguration, not per-request magic.

## Quality Bar
- Performance bar: choose the relative winner by measured evidence envelope on this host, not an arbitrary SLA. Compare cold/warm latency, throughput, batch behavior, and sanitized config.
- Legal quality bar: no regression versus TEI/USER-bge-m3 on the current legal parity gate. Alternative models cannot replace baseline without legal-domain proof.
- Final deliverable: recommendation artifact with runtime recommendation, evidence table, operating contract, and remaining caveats.

## Known Prior Evidence
- M038: local Go ONNX target-runtime smoke/legal/performance proof passed.
- M039: packaged Docker Go ONNX smoke/legal/performance proof passed for `fd-api:onnx1024-m039`.
- M039 operational finding: packaged ONNX container needs `ONNX_RUNTIME_SHA256` for `/health` to report `runtime_library_verified=true`.
- Hosted GitHub Actions proof is not required as an ONNX acceptance/production gate.

## Open Work
- Depth verification.
- Requirements and roadmap approval.
