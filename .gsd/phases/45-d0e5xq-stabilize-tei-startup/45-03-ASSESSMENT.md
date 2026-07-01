---
sliceId: S03
uatType: mixed
verdict: PASS
attempt: 3
runId: uat:M045-d0e5xq:S03:attempt-3
worktreeRoot: /root/fd
date: 2026-06-14T12:24:26.889Z
---

# UAT Result - S03

## Checks

| Check | Mode | Result | Evidence | Notes |
|-------|------|--------|----------|-------|
| Local path proof runtime artifact shows healthy outcome and embedding smoke from controlled restart. | runtime | PASS | gsd_uat_exec:335f897b-97a6-4a0d-bb21-44e2df9fc8cc | Runtime output from the controlled restart proof. |
| Running TEI container is healthy and uses local snapshot command. | runtime | PASS | gsd_uat_exec:a55d4583-cd51-4b8f-943a-b56c1a045b1e | Confirms the proof is applied to the current running container. |
| fd health and embedding smoke pass after local path startup. | runtime | PASS | gsd_uat_exec:708de145-0232-4cc6-9346-75b837562a25 | Preserves client-visible runtime contract. |
| Effective compose runtime config uses local path and excludes failed HF_HUB_OFFLINE candidate. | runtime | PASS | gsd_uat_exec:298cc0b1-45c5-48cd-9748-eca5e1854320 | Ensures the rejected offline env is not in effective runtime config. |
| Browser verification of localhost health JSON shows TEI runtime identity. | browser | PASS | browser:/root/fd/.artifacts/browser/2026-06-14T12-23-36-082Z-session/m045-s03-browser-health-timeline.json | Added to satisfy browser-inclusive UAT guard for the localhost runtime surface. |

## Overall Verdict

PASS - PASS: mixed live runtime/browser UAT validates the local snapshot startup mitigation; R028 is validated.

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
    "browser_navigate",
    "browser_assert",
    "browser_timeline",
    "find",
    "glob",
    "grep",
    "ls",
    "read",
    "browser_click",
    "browser_type",
    "browser_fill_form",
    "browser_click_ref",
    "browser_fill_ref",
    "browser_wait_for",
    "browser_verify",
    "browser_screenshot",
    "browser_snapshot_refs",
    "browser_find",
    "browser_get_console_logs",
    "browser_get_network_logs",
    "browser_evaluate",
    "browser_reload",
    "browser_batch",
    "browser_act"
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
  "fallbackToolsUsed": [
    "browser_navigate",
    "browser_assert",
    "browser_timeline"
  ],
  "notes": "Mixed UAT includes live runtime container/HTTP evidence plus browser verification of the localhost health JSON surface.",
  "toolPresentationPlanId": "run-uat/default-v1"
}
```

## Gate

Aggregate UAT gate saved as pass.
