# M042 Post-M043 Replan Assessment

## Verdict

S02 and S03 were replanned because M043 static-analysis hardening materially changes how async/ONNX code should be implemented and verified.

## Reason

M042 introduces concurrency, runtime selection, ONNX native artifacts, and performance harnesses. These are high-risk areas for the new M043 gates: `contextcheck`, `gocyclo`, `gocritic`, `revive:exported`, `gosec` path/URL handling, and govulncheck.

## Changes applied

- S02 task plans now require bounded concurrency helpers, context propagation, small production functions, async observability, and final lint/test/govulncheck gates.
- S03 task plans now require safe operator-controlled path handling, cache namespace isolation for TEI vs ONNX comparisons, smoke embedding readiness proof (not just /health 200), ONNX build/perf proof, and docs explaining legal quality gate deferral.

## Important gotchas encoded in the plan

- Isolate Redis cache namespace when comparing TEI and ONNX; otherwise cached TEI vectors can make ONNX falsely appear equivalent.
- `/health` runtime metadata is not a live inference probe; readiness still requires a smoke embedding request.
- New exported ONNX/config APIs need useful godoc.
- Any new dependency or native runtime change must keep govulncheck passing with 0 reachable vulnerabilities.
