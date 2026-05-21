---
id: M038-pmw50e
title: "Go ONNX target runtime acceptance proof"
status: complete
completed_at: 2026-05-21T10:49:27.859Z
key_decisions: []
key_files:
  - benchmark-results/fd-onnx-go-runtime-smoke-m038-s01.txt
  - benchmark-results/fd-legal-retrieval-m038-go-onnx-target-runtime.txt
  - benchmark-results/fd-benchmark-m038-go-onnx-target-runtime.txt
  - benchmark-results/fd-onnx-go-target-runtime-acceptance-m038-s02.txt
lessons_learned:
  - The live Go ONNX package test runs with package cwd `api/embed`, so manifest/tokenizer test paths must be package-relative.
  - Python legal/performance tools can produce target-runtime evidence when configured against actual Go endpoints and isolated cache namespaces.
---

# M038-pmw50e: Go ONNX target runtime acceptance proof

**Produced fresh local Go ONNX target-runtime smoke, legal, and performance evidence for the current artifact.**

## What Happened

M038 executed the first fresh Go target-runtime acceptance proof for the current local ONNX artifact. It verified local artifacts, passed the live Go ONNX embedder test, started the Go ONNX API with isolated namespaces, verified health and embedding smoke, passed the Russian/legal retrieval gate through actual Go endpoints, and ran performance benchmark against the Go ONNX endpoint. The acceptance matrix records passed/skipped/remaining gates. TEI remains production/default and ONNX remains opt-in experimental.

## Success Criteria Results

- PASS — Fresh Go ONNX target-runtime evidence exists.
- PASS — Evidence is through actual Go endpoints, not Python-only helpers.
- PASS — Redis namespace isolation explicit.
- PASS — ONNX remains opt-in experimental and TEI remains production/default.

## Definition of Done Results

- Done — Live Go ONNX embedder test passed.
- Done — Go ONNX API smoke passed with health and 1024-dimensional embedding proof.
- Done — Legal retrieval evaluator passed through actual Go endpoints.
- Done — Performance benchmark passed against actual Go ONNX endpoint.
- Done — Outcome records exact commands, namespaces, and non-actions without raw text/secrets.
- Done — Final guardrails passed.

## Requirement Outcomes

The implicit target-runtime validation requirement is advanced with local Go evidence, but packaged and hosted proof remain open.

## Deviations

Redis L2 restart subchecks were skipped because benchmark.py was run against a bg_shell-managed Go ONNX server with empty restart command. Packaged Docker ONNX legal/performance reruns remain future gates.

## Follow-ups

Next recommended gate: packaged Docker Go ONNX target-runtime rerun for the current artifact/source setup, or a reusable target-runtime harness that automates the smoke/legal/performance sequence.
