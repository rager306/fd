---
sliceId: S01
uatType: runtime-executable
verdict: PASS
attempt: 1
runId: uat:M045-d0e5xq:S01:attempt-1
worktreeRoot: /root/fd
date: 2026-06-14T11:31:42.952Z
---

# UAT Result - S01

## Checks

| Check | Mode | Result | Evidence | Notes |
|-------|------|--------|----------|-------|
| Recon artifact contains safety boundary, current runtime evidence, and candidate TEI startup findings. | artifact | PASS | gsd_uat_exec:93a29b67-73ac-4644-9fdc-c5054ea171c2 | Artifact path: documents/tei-startup-recon-m045.md. |
| fd health still reports TEI-only runtime identity. | runtime | PASS | gsd_uat_exec:307f2271-c9d4-4315-b4ea-89b6c52d5e18 | Preserves R027 TEI-only posture. |
| Direct TEI embedding smoke works without restart. | runtime | PASS | gsd_uat_exec:558c59fc-491f-4117-91b4-98dbdcc3f537 | Read-only runtime check; no container restart performed. |

## Overall Verdict

PASS - PASS: S01 completed read-only TEI startup recon and identified `HF_HUB_OFFLINE=1` plus local model path as candidate mitigations for S02.

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
      "name": "docker restart",
      "reason": "Out of S01 safety boundary."
    },
    {
      "name": "docker compose restart",
      "reason": "Out of S01 safety boundary."
    },
    {
      "name": "docker run",
      "reason": "Previously blocked as destructive and not needed for S01."
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
  "notes": "S01 UAT verifies non-destructive recon and runtime health only.",
  "toolPresentationPlanId": "run-uat/default-v1"
}
```

## Gate

Aggregate UAT gate saved as pass.
