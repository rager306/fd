#!/usr/bin/env python3
"""Validate sanitized legal model quick-gate artifacts."""

from __future__ import annotations

import argparse
import json
from pathlib import Path
import re
import tempfile
from typing import Any

ALLOWED_OUTCOMES = {"keep_current", "reject_candidate", "defer_candidate"}
REQUIRED_SECTIONS = {
    "Effective Configuration",
    "Baseline Status",
    "Candidate Results",
    "Cross-Model Cosine/Parity",
    "Verdict",
    "Redaction",
}
SECRET_PATTERNS = [
    re.compile(r"(?i)bearer\s+[a-z0-9._~+/=-]{12,}"),
    re.compile(r"(?i)(api[_-]?key|token|password|secret)\s*[:=]\s*[a-z0-9._~+/=-]{8,}"),
    re.compile(r"sk-[A-Za-z0-9]{16,}"),
    re.compile(r"hf_[A-Za-z0-9]{16,}"),
]
RAW_TEXT_PATTERNS = [
    re.compile(r"(?im)^\s*##\s+Raw Text\b"),
    re.compile(r"(?im)^\s*raw_text\s*[:=]"),
]


class ArtifactError(ValueError):
    pass


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Validate legal model quick-gate markdown artifact.")
    parser.add_argument("--artifact", type=Path, help="Artifact markdown path to validate.")
    parser.add_argument("--max-candidates", type=int, default=2)
    parser.add_argument("--self-test", action="store_true", help="Run inline pass/fail fixture tests.")
    return parser.parse_args()


def section_names(markdown: str) -> set[str]:
    names: set[str] = set()
    for line in markdown.splitlines():
        if line.startswith("## "):
            names.add(line[3:].strip())
    return names


def section_text(markdown: str, name: str) -> str:
    marker = f"## {name}"
    start = markdown.find(marker)
    if start < 0:
        raise ArtifactError(f"missing section: {name}")
    start = markdown.find("\n", start)
    if start < 0:
        return ""
    next_section = markdown.find("\n## ", start + 1)
    if next_section < 0:
        next_section = len(markdown)
    return markdown[start:next_section].strip()


def json_section(markdown: str, name: str) -> Any:
    body = section_text(markdown, name)
    match = re.search(r"```json\s*(.*?)\s*```", body, flags=re.S)
    if not match:
        raise ArtifactError(f"section {name} missing json code fence")
    try:
        return json.loads(match.group(1))
    except json.JSONDecodeError as exc:
        raise ArtifactError(f"section {name} contains invalid json: {exc}") from exc


def verdict(markdown: str) -> str:
    body = section_text(markdown, "Verdict")
    lines = [line.strip() for line in body.splitlines() if line.strip()]
    if not lines:
        raise ArtifactError("missing verdict value")
    value = lines[0]
    if value not in ALLOWED_OUTCOMES:
        raise ArtifactError(f"invalid verdict: {value}")
    return value


def contains_secret(markdown: str) -> str | None:
    for pattern in SECRET_PATTERNS:
        match = pattern.search(markdown)
        if match:
            return match.group(0)[:40]
    return None


def validate_artifact_text(markdown: str, max_candidates: int = 2) -> None:
    missing = REQUIRED_SECTIONS - section_names(markdown)
    if missing:
        raise ArtifactError(f"missing required sections: {', '.join(sorted(missing))}")
    secret = contains_secret(markdown)
    if secret:
        raise ArtifactError(f"artifact contains secret-like pattern: {secret}")
    for pattern in RAW_TEXT_PATTERNS:
        if pattern.search(markdown):
            raise ArtifactError("artifact contains a raw-text section or field")

    config = json_section(markdown, "Effective Configuration")
    baseline = json_section(markdown, "Baseline Status")
    candidates = json_section(markdown, "Candidate Results")
    cosine = json_section(markdown, "Cross-Model Cosine/Parity")
    redaction = json_section(markdown, "Redaction")
    final = verdict(markdown)

    if not isinstance(config, dict):
        raise ArtifactError("configuration section must be an object")
    if config.get("raw_text_logged") is not False:
        raise ArtifactError("config.raw_text_logged must be false")
    if not config.get("redaction_status"):
        raise ArtifactError("configuration missing redaction_status")
    if int(config.get("candidate_count", -1)) > max_candidates:
        raise ArtifactError(f"candidate_count exceeds max {max_candidates}")
    if not isinstance(baseline, dict):
        raise ArtifactError("baseline status must be an object")
    for key in ("api_url_label", "model", "runtime_label", "expected_dimensions", "cache_namespace"):
        if key not in baseline:
            raise ArtifactError(f"baseline missing {key}")
    if not isinstance(candidates, list):
        raise ArtifactError("candidate results must be an array")
    if len(candidates) > max_candidates:
        raise ArtifactError(f"artifact has {len(candidates)} candidates; max {max_candidates}")
    if int(config.get("candidate_count", len(candidates))) != len(candidates):
        raise ArtifactError("config candidate_count does not match candidate results")
    if not candidates:
        raise ArtifactError("at least one candidate result is required")

    for idx, candidate in enumerate(candidates, 1):
        if not isinstance(candidate, dict):
            raise ArtifactError(f"candidate {idx} must be an object")
        for key in ("api_url_label", "model", "runtime_label", "expected_dimensions", "cache_namespace", "outcome", "stop_reason"):
            if key not in candidate:
                raise ArtifactError(f"candidate {idx} missing {key}")
        if candidate["outcome"] not in ALLOWED_OUTCOMES:
            raise ArtifactError(f"candidate {idx} invalid outcome {candidate['outcome']}")
        metrics = candidate.get("metrics")
        if candidate["outcome"] == "reject_candidate" and not metrics:
            raise ArtifactError(f"candidate {idx} rejected without metrics")
        if candidate["outcome"] == "defer_candidate" and not candidate.get("stop_reason"):
            raise ArtifactError(f"candidate {idx} deferred without stop_reason")
        if metrics and "against_baseline" in metrics:
            ratios = metrics["against_baseline"]
            if not {"recall_ratio", "mrr_ratio"}.issubset(ratios):
                raise ArtifactError(f"candidate {idx} metrics missing baseline ratios")

    if not isinstance(cosine, dict) or cosine.get("applicable") is not False:
        raise ArtifactError("cross-model cosine/parity must be explicitly not applicable")
    cosine_reason = str(cosine.get("reason", "")).lower()
    if "acceptance metric" not in cosine_reason or "not" not in cosine_reason:
        raise ArtifactError("cross-model cosine/parity reason must reject acceptance use")
    acceptance_claim_patterns = [
        r"(?i)cross[-_ ]model[^\n]{0,120}(accepted as|approved as|passes as|pass acceptance)",
        r"(?i)cross[-_ ]model[^\n]{0,120}acceptance[_ -]?metric[\"']?\s*[:=]\s*true",
    ]
    if any(re.search(pattern, markdown) for pattern in acceptance_claim_patterns):
        raise ArtifactError("artifact appears to treat cross-model cosine as acceptance")

    if not isinstance(redaction, dict) or redaction.get("raw_text_logged") is not False:
        raise ArtifactError("redaction.raw_text_logged must be false")
    if "raw legal corpus text" not in str(redaction.get("statement", "")).lower():
        raise ArtifactError("redaction statement must mention raw legal corpus text")
    candidate_outcomes = {candidate["outcome"] for candidate in candidates}
    if final == "reject_candidate" and not candidate_outcomes.issubset({"reject_candidate"}):
        raise ArtifactError("reject_candidate verdict inconsistent with candidate outcomes")
    if final == "defer_candidate" and "defer_candidate" not in candidate_outcomes:
        raise ArtifactError("defer_candidate verdict requires a deferred candidate")


def validate_path(path: Path, max_candidates: int = 2) -> None:
    validate_artifact_text(path.read_text(encoding="utf-8"), max_candidates=max_candidates)


def sample_artifact(candidate_count: int = 1, verdict_value: str = "defer_candidate", raw_text_logged: bool = False, extra: str = "") -> str:
    candidates = [
        {
            "api_url_label": "not-called-dry-run",
            "model": "BAAI/bge-m3",
            "runtime_label": "candidate-bge-m3",
            "expected_dimensions": 1024,
            "cache_namespace": f"candidate-{idx}",
            "health": {"ok": False, "phase": "dry_run"},
            "smoke_embedding": {"ok": False, "phase": "dry_run"},
            "metrics": None,
            "outcome": "defer_candidate",
            "stop_reason": "dry_run_availability_only_no_endpoint_calls",
        }
        for idx in range(candidate_count)
    ]
    return "\n".join(
        [
            "# M040 S03 Legal Model Quick Gate",
            "",
            "## Effective Configuration",
            "```json",
            json.dumps({"raw_text_logged": raw_text_logged, "redaction_status": "sanitized", "candidate_count": candidate_count}),
            "```",
            "",
            "## Baseline Status",
            "```json",
            json.dumps({"api_url_label": "not-called-dry-run", "model": "deepvk/USER-bge-m3", "runtime_label": "tei-default", "expected_dimensions": 1024, "cache_namespace": "default"}),
            "```",
            "",
            "## Candidate Results",
            "```json",
            json.dumps(candidates),
            "```",
            "",
            "## Cross-Model Cosine/Parity",
            "```json",
            json.dumps({"applicable": False, "reason": "not an acceptance metric for different models"}),
            "```",
            "",
            "## Verdict",
            verdict_value,
            "",
            "## Redaction",
            "```json",
            json.dumps({"raw_text_logged": raw_text_logged, "statement": "Raw legal corpus text excluded."}),
            "```",
            extra,
            "",
        ]
    )


def expect_pass(name: str, text: str) -> None:
    validate_artifact_text(text)
    print(f"PASS expected-pass {name}")


def expect_fail(name: str, text: str) -> None:
    try:
        validate_artifact_text(text)
    except ArtifactError as exc:
        print(f"PASS expected-fail {name}: {exc}")
        return
    raise AssertionError(f"expected fixture to fail: {name}")


def run_self_test() -> None:
    with tempfile.TemporaryDirectory() as tmpdir:
        path = Path(tmpdir) / "artifact.md"
        text = sample_artifact()
        path.write_text(text, encoding="utf-8")
        validate_path(path)
        expect_pass("valid", text)
        expect_fail("too-many-candidates", sample_artifact(candidate_count=3))
        expect_fail("missing-required-metadata", sample_artifact().replace('"runtime_label": "candidate-bge-m3",', ""))
        expect_fail("raw-text-logged", sample_artifact(raw_text_logged=True))
        expect_fail("secret-pattern", sample_artifact(extra="token=supersecrettoken12345"))
        expect_fail("missing-verdict", sample_artifact().replace("## Verdict\ndefer_candidate", "## Verdict\n"))
        expect_fail("cross-model-cosine-acceptance", sample_artifact().replace("not an acceptance metric", "accepted as pass acceptance metric"))


def main() -> int:
    args = parse_args()
    if args.self_test:
        run_self_test()
        return 0
    if not args.artifact:
        raise SystemExit("--artifact is required unless --self-test is used")
    validate_path(args.artifact, max_candidates=args.max_candidates)
    print(f"artifact valid: {args.artifact}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
