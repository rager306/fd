# S03 Replan

**Milestone:** M014-vjfs9f
**Slice:** S03
**Blocker Task:** T02
**Created:** 2026-05-20T04:24:45.761Z

## Blocker Description

The existing benchmark Redis L2 restart check always restarted Docker Compose `api`, which is valid for TEI on port 8000 but invalid for a manually started tagged ONNX server on port 18000. Without a configurable restart command, S03 would not actually test ONNX Redis L2 after benchmarked-server restart.

## What Changed

S03 now explicitly accounts for the benchmark restart-harness fix and expects snapshot v3 for ONNX benchmark artifacts.
