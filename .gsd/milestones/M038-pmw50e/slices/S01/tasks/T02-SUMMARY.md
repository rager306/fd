---
id: T02
parent: S01
milestone: M038-pmw50e
key_files: []
key_decisions: []
duration: 
verification_result: mixed
completed_at: 2026-05-21T10:34:42.588Z
blocker_discovered: false
---

# T02: Ran the live Go ONNX embedder test successfully after correcting package-relative paths.

**Ran the live Go ONNX embedder test successfully after correcting package-relative paths.**

## What Happened

Ran the live tagged Go ONNX embedder test against the current local artifact. The first attempt exposed a path issue in the harness invocation; the corrected run used package-relative manifest/tokenizer paths and passed. This proves the Go ONNX embedder can load the current artifact through the native tokenizer and ONNX Runtime path and produce a normalized 1024-dimensional embedding.

## Verification

Live tagged Go test passed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `go test -tags 'onnx hf_tokenizers' ./embed -run TestONNXEmbedderLiveLocalArtifact (initial path attempt)` | 1 | ❌ fail — tokenizer path invalid due to package cwd | 6300ms |
| 2 | `go test -tags 'onnx hf_tokenizers' ./embed -run TestONNXEmbedderLiveLocalArtifact with ../../ manifest/tokenizer paths` | 0 | ✅ pass — 1 passed in 1 package | 7800ms |

## Deviations

First run used API-root relative tokenizer/manifest paths, but Go package tests execute with cwd `api/embed`; rerun used package-relative `../../` paths and passed.

## Known Issues

The test only proves direct Go embedder load/infer, not HTTP API, legal quality, or packaged runtime behavior.

## Files Created/Modified

None.
