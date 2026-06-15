---
sliceId: S02
uatType: artifact-driven
verdict: PASS
attempt: 1
runId: uat:M048-l4sctg:S02:attempt-1
worktreeRoot: /root/fd
date: 2026-06-15T11:19:03.623Z
---

# UAT Result - S02

## Checks

| Check | Mode | Result | Evidence | Notes |
|-------|------|--------|----------|-------|
| RuntimeHealth no longer exposes inactive ONNX-only fields. | artifact | PASS | gsd_uat_exec:a242ec62-a82e-48e6-962f-f21f9b87bc28 |  |
| Handlers and lifecycle use shared embed.Embedder. | artifact | PASS | gsd_uat_exec:7fe8e369-348d-4184-8292-880b5425da2b |  |
| Lifecycle state is explicitly constructed in main. | artifact | PASS | gsd_uat_exec:346c2f5f-bc4b-4905-98ff-8f9edb6b37c4 |  |
| S02 evidence artifact covers #26/#29/#30 and R038. | artifact | PASS | gsd_uat_exec:7ffa9040-28ba-4f7f-aa9f-e50b58a036ed |  |

## Overall Verdict

PASS - PASS: S02 validates issue #7 findings #26, #29, and #30.

## Tool Presentation

```json
{
  "surface": "provider-tools",
  "presentedTools": [
    "gsd_resume",
    "gsd_milestone_status",
    "gsd_journal_query",
    "gsd_uat_exec",
    "gsd_uat_result_save",
    "find",
    "glob",
    "grep",
    "ls",
    "read"
  ],
  "blockedTools": [
    {
      "name": "gsd_exec",
      "reason": "forbidden during run-uat"
    },
    {
      "name": "gsd_summary_save",
      "reason": "forbidden during run-uat"
    },
    {
      "name": "gsd_save_gate_result",
      "reason": "forbidden during run-uat"
    },
    {
      "name": "browser",
      "reason": "S02 validates backend contract source shape, not browser UI."
    },
    {
      "name": "edit",
      "reason": "forbidden during run-uat"
    },
    {
      "name": "write",
      "reason": "forbidden during run-uat"
    },
    {
      "name": "search-the-web",
      "reason": "forbidden during run-uat"
    },
    {
      "name": "WebSearch",
      "reason": "forbidden during run-uat"
    },
    {
      "name": "Bash",
      "reason": "forbidden during run-uat"
    },
    {
      "name": "Write",
      "reason": "forbidden during run-uat"
    },
    {
      "name": "Edit",
      "reason": "forbidden during run-uat"
    }
  ],
  "fallbackToolsUsed": [],
  "notes": "Artifact-driven UAT is sufficient for S02 runtime contract cleanup.",
  "toolPresentationPlanId": "run-uat/default-v1"
}
```

## Gate

Aggregate UAT gate saved as pass.
