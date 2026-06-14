# S01 Replan

**Milestone:** M042-fjf2en
**Slice:** S01
**Blocker Task:** T01
**Created:** 2026-06-14T10:13:06.613Z

## Blocker Description

During T02-style profiling, TEI restart/recreate became unhealthy for >10 minutes and remained stuck around backend warmup/startup. User also explicitly rescoped the project away from ONNX runtime work, preferring TEI-first stabilization and future-only ONNX research. The existing S01 T02/T03 plan overemphasized ONNX and destructive restart profiling.

## What Changed

Preserve completed T01 snapshot. Replace T02 with a non-destructive TEI startup/concurrency failure analysis artifact using already captured failed restart evidence and any safe read-only diagnostics. Replace T03 with a TEI-first RCA/verdict document that defers ONNX and recommends TEI operational stabilization/mitigation only.
