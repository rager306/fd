#!/usr/bin/env python3
"""Verify the M040 S04 TEI-vs-ONNX final recommendation artifact.

This verifier is intentionally semantic rather than performance-threshold based.
It fails closed when the final recommendation drifts outside the evidence
contract: TEI remains the default until explicit operator opt-in; packaged ONNX
for deepvk/USER-bge-m3 is only preferred under same-host/S01-style operating
checks; Redis cache namespaces must be isolated; /health is not live inference
readiness; S03 candidate-model replacement remains deferred/fail-closed; and
artifacts must not leak secrets or raw benchmark/legal probe text.
"""

from __future__ import annotations

import argparse
import re
import sys
import tempfile
from pathlib import Path

EXPECTED_MODEL = "deepvk/USER-bge-m3"
EXPECTED_ONNX_ARTIFACT = "user-bge-m3-dense-fp32"
EXPECTED_CACHE_VERSION = "m040-s02-onnx-restart"
EXPECTED_RUNTIME_LABEL = "onnx-docker-m040-s02"
EXPECTED_LEGAL_RUNTIME_LABEL = "docker-onnx-go-api-m040-s02"

DEFAULT_SOURCE_PATHS = {
    "s02_verifier": Path("tools/verify_m040_s02_artifacts.py"),
    "s03_verifier": Path("tools/verify_legal_model_quick_gate_artifact.py"),
    "s02_benchmark": Path("benchmark-results/fd-benchmark-m040-s02-onnx-docker-restart.txt"),
    "s02_preflight": Path("benchmark-results/fd-m040-s02-onnx-docker-preflight.txt"),
    "s02_legal": Path("benchmark-results/fd-legal-retrieval-m040-s02-onnx-docker-restart.txt"),
    "s02_audit": Path("benchmark-results/fd-m040-s02-proof-audit.txt"),
    "s03_gate": Path("benchmark-results/fd-legal-model-quick-gate-m040-s03.md"),
    "contract": Path("docs/same-host-embedding-service-contract.md"),
}

REQUIRED_SECTIONS = {
    "Recommendation",
    "Decision Inputs",
    "Same-Host Operating Contract",
    "Required Operator Checks",
    "Evidence Links",
    "Caveats",
    "Non-Goals",
    "Redaction",
}

PROHIBITED_PATTERNS = [
    re.compile(r"-----BEGIN [A-Z ]*PRIVATE KEY-----"),
    re.compile(r"\bBearer\s+[A-Za-z0-9._~+/=-]{12,}", re.IGNORECASE),
    re.compile(r"(?i)(api[_-]?key|secret|password|credential|auth[_-]?token)\s*[:=]\s*[^\s,;}]{6,}"),
    re.compile(r"sk-[A-Za-z0-9]{16,}"),
    re.compile(r"hf_[A-Za-z0-9]{16,}"),
    re.compile(r"https?://[^\s]+\?(?:[^\s]*\b(?:X-Amz-Signature|Signature|token|access_token)=)", re.IGNORECASE),
    # Raw probe/legal text that should never be needed in the recommendation.
    re.compile(r"Привет, как дела\?"),
    re.compile(r"Москва — столица России"),
    re.compile(r"Искусственный интеллект — это область компьютерных наук"),
    re.compile(r"В современном мире технологии машинного обучения"),
    re.compile(r"Redis L2 persistence diagnostic text"),
    re.compile(r"batch-cache-item-\d+"),
    re.compile(r"research-chunk-\d+"),
    re.compile(r"(?im)^\s*raw_text\s*[:=]\s*['\"]?[^\n]+"),
    re.compile(r"(?im)^\s*##\s+Raw Text\b"),
]


class ArtifactError(ValueError):
    """Raised for a verifier failure with a future-agent-readable message."""


def fail(message: str) -> None:
    raise ArtifactError(message)


def read_required(path: Path) -> str:
    if not path.exists():
        fail(f"missing evidence file: {path}")
    if not path.is_file():
        fail(f"evidence path is not a file: {path}")
    return path.read_text(encoding="utf-8", errors="replace")


def section_names(markdown: str) -> set[str]:
    return {line[3:].strip() for line in markdown.splitlines() if line.startswith("## ")}


def section_text(markdown: str, name: str) -> str:
    marker = f"## {name}"
    start = markdown.find(marker)
    if start < 0:
        fail(f"missing section: {name}")
    body_start = markdown.find("\n", start)
    if body_start < 0:
        return ""
    next_section = markdown.find("\n## ", body_start + 1)
    if next_section < 0:
        next_section = len(markdown)
    return markdown[body_start:next_section].strip()


def require_all(text: str, needles: list[str], context: str) -> None:
    missing = [needle for needle in needles if needle not in text]
    if missing:
        fail(f"{context} missing required marker(s): {', '.join(missing)}")


def require_any(text: str, needles: list[str], context: str) -> None:
    if not any(needle in text for needle in needles):
        fail(f"{context} missing one of: {', '.join(needles)}")


def assert_no_prohibited_patterns(text: str, context: str) -> None:
    for pattern in PROHIBITED_PATTERNS:
        match = pattern.search(text)
        if match:
            snippet = match.group(0)[:80].replace("\n", " ")
            fail(f"{context} contains prohibited secret/raw-text pattern {pattern.pattern!r}: {snippet!r}")


def assert_no_hosted_ci_gate(text: str) -> None:
    for line_number, line in enumerate(text.splitlines(), start=1):
        lower = line.lower()
        mentions_ci = "hosted ci" in lower or "github actions" in lower or re.search(r"\bci\b", lower)
        mentions_gate = "readiness gate" in lower or ("readiness" in lower and any(word in lower for word in ["gate", "required", "must pass", "requirement"]))
        if not mentions_ci or not mentions_gate:
            continue
        negated = any(
            phrase in lower
            for phrase in [
                "not a readiness gate",
                "not readiness gate",
                "is not a readiness gate",
                "is not readiness",
                "never a readiness gate",
                "must not be a readiness gate",
                "non-goal",
            ]
        )
        if not negated:
            fail(f"hosted CI is incorrectly treated as a readiness gate on line {line_number}: {line.strip()}")


def assert_candidate_replacement_deferred(text: str) -> None:
    candidate_models = ["BAAI/bge-m3", "intfloat/multilingual-e5-large"]
    positive_pattern = re.compile(r"(?i)\b(accepted?|approved?|recommended?|replace|replacement|switch to)\b")
    deferred_pattern = re.compile(r"(?i)\b(defer(?:red)?|fail-closed|failed_closed|not accepted|not recommended|non-goal)\b")
    for line_number, line in enumerate(text.splitlines(), start=1):
        for model in candidate_models:
            if model not in line:
                continue
            if positive_pattern.search(line) and not deferred_pattern.search(line):
                fail(f"candidate model replacement appears accepted rather than deferred on line {line_number}: {model}")


def validate_source_evidence(paths: dict[str, Path]) -> dict[str, str]:
    texts = {name: read_required(path) for name, path in paths.items()}
    for name, text in texts.items():
        if name not in {"s02_verifier", "s03_verifier"}:
            assert_no_prohibited_patterns(text, str(paths[name]))

    require_all(
        texts["s02_verifier"],
        ["EXPECTED_CACHE_VERSION", "verify_preflight", "verify_benchmark", "verify_legal", "assert_no_leaks"],
        str(paths["s02_verifier"]),
    )
    require_all(
        texts["s03_verifier"],
        ["ALLOWED_OUTCOMES", "defer_candidate", "Cross-Model Cosine/Parity", "raw_text_logged"],
        str(paths["s03_verifier"]),
    )

    require_all(
        texts["s02_preflight"],
        [
            "phase=health status=ready",
            "phase=smoke status=requesting_embedding",
            '"backend": "onnx"',
            f'"model": "{EXPECTED_MODEL}"',
            f'"artifact_id": "{EXPECTED_ONNX_ARTIFACT}"',
            '"artifact_verified": true',
            '"tokenizer_verified": true',
            EXPECTED_CACHE_VERSION,
        ],
        str(paths["s02_preflight"]),
    )
    require_all(
        texts["s02_benchmark"],
        [
            "Redis L2 Persistence — After API Restart",
            "after API restart:",
            "redis_delta/l2_after_api_restart",
            "Cached Batch Endpoint — L1 and Redis L2",
            f'"runtime_label": "{EXPECTED_RUNTIME_LABEL}"',
            f'"EMBEDDING_CACHE_VERSION": "{EXPECTED_CACHE_VERSION}"',
            '"raw_benchmark_texts_excluded": true',
        ],
        str(paths["s02_benchmark"]),
    )
    require_all(
        texts["s02_legal"],
        [
            "## Verdict",
            "PASS",
            '"raw_text_logged": false',
            f'"model": "{EXPECTED_MODEL}"',
            f'"runtime_label": "{EXPECTED_LEGAL_RUNTIME_LABEL}"',
            f'"cache_namespace": "{EXPECTED_CACHE_VERSION}"',
        ],
        str(paths["s02_legal"]),
    )
    require_all(
        texts["s02_audit"],
        [
            "# M040 S02 Proof Audit",
            "## Legal Gate",
            "verdict: PASS",
            "## Leak Audit",
            "details: no prohibited patterns found",
            "## Blockers",
            "None.",
        ],
        str(paths["s02_audit"]),
    )
    require_all(
        texts["s03_gate"],
        [
            "# M040 S03 Legal Model Quick Gate",
            "defer_candidate",
            "failed_closed",
            "BAAI/bge-m3",
            "intfloat/multilingual-e5-large",
            '"raw_text_logged": false',
            "not an acceptance metric",
        ],
        str(paths["s03_gate"]),
    )
    require_all(
        texts["contract"],
        [
            "Production default: TEI",
            "Opt-in: ONNX",
            "What `/health` does NOT prove",
            "Clients sharing a Redis instance",
            "No-Silent-Fallback Rules",
        ],
        str(paths["contract"]),
    )
    return texts


def validate_recommendation_text(markdown: str, source_paths: dict[str, Path] | None = None) -> None:
    assert_no_prohibited_patterns(markdown, "recommendation artifact")
    assert_no_hosted_ci_gate(markdown)
    assert_candidate_replacement_deferred(markdown)

    missing = REQUIRED_SECTIONS - section_names(markdown)
    if missing:
        fail(f"missing required section(s): {', '.join(sorted(missing))}")

    recommendation = section_text(markdown, "Recommendation")
    require_all(recommendation, [EXPECTED_MODEL], "Recommendation section")
    require_any(recommendation.lower(), ["packaged onnx", "onnx package", "onnx runtime"], "Recommendation section")
    require_any(recommendation.lower(), ["prefer", "preferred", "recommend"], "Recommendation section")
    require_any(recommendation.lower(), ["tei remains", "tei stays", "tei is the current", "tei is current", "default posture"], "Recommendation section")
    require_any(recommendation.lower(), ["explicit switch", "operator explicitly switches", "opt-in", "explicit operator"], "Recommendation section")
    require_any(recommendation.lower(), ["same-host", "s01 contract", "same host"], "Recommendation section")

    decision_inputs = section_text(markdown, "Decision Inputs")
    require_all(
        decision_inputs,
        ["S02", "S03", EXPECTED_ONNX_ARTIFACT, EXPECTED_CACHE_VERSION, "defer_candidate"],
        "Decision Inputs section",
    )

    contract = section_text(markdown, "Same-Host Operating Contract")
    require_all(contract, ["/health", "smoke", "Redis", "cache", "namespace"], "Same-Host Operating Contract section")
    require_any(contract.lower(), ["not live inference", "does not perform", "not an inference", "smoke embedding"], "Same-Host Operating Contract section")
    require_any(contract.lower(), ["isolate", "disjoint", "separate", "namespace"], "Same-Host Operating Contract section")

    checks = section_text(markdown, "Required Operator Checks")
    require_all(
        checks,
        ["ONNX_ARTIFACT_MANIFEST", "ONNX_TOKENIZER_PATH", "ONNX_RUNTIME_LIBRARY", "EMBEDDING_BACKEND=onnx"],
        "Required Operator Checks section",
    )
    require_any(checks.lower(), ["smoke", "embedding request"], "Required Operator Checks section")
    require_any(checks.lower(), ["redis", "cache namespace", "embedding_cache_version"], "Required Operator Checks section")
    require_any(checks.lower(), ["artifact", "tokenizer", "runtime"], "Required Operator Checks section")

    caveats = section_text(markdown, "Caveats")
    require_all(caveats, ["legal", "quality"], "Caveats section")
    require_any(caveats.lower(), ["no fallback", "fail closed", "no silent fallback"], "Caveats section")
    require_any(caveats.lower(), ["not prove", "does not prove", "not a quality benchmark", "bounded"], "Caveats section")

    non_goals = section_text(markdown, "Non-Goals")
    require_any(non_goals.lower(), ["candidate", "alternative model", "replacement"], "Non-Goals section")
    require_any(non_goals.lower(), ["hosted ci", "github actions", "ci"], "Non-Goals section")
    require_any(non_goals.lower(), ["not", "non-goal", "defer", "fail-closed"], "Non-Goals section")

    redaction = section_text(markdown, "Redaction")
    require_all(redaction.lower(), ["raw", "text", "secret"], "Redaction section")
    require_any(redaction.lower(), ["excluded", "redacted", "not logged", "no raw"], "Redaction section")

    links = section_text(markdown, "Evidence Links")
    paths = source_paths or DEFAULT_SOURCE_PATHS
    for path in paths.values():
        if str(path) not in links:
            fail(f"Evidence Links section missing source artifact reference: {path}")


def validate_path(artifact: Path, source_paths: dict[str, Path]) -> None:
    validate_source_evidence(source_paths)
    validate_recommendation_text(read_required(artifact), source_paths)


def sample_recommendation(extra: str = "") -> str:
    links = "\n".join(f"- `{path}`" for path in DEFAULT_SOURCE_PATHS.values())
    return f"""# M040 S04 Runtime Recommendation

## Recommendation
Prefer packaged ONNX runtime for {EXPECTED_MODEL} as the same-host performance runtime only under the S01 contract and explicit operator opt-in. TEI remains the current/default posture until an operator explicitly switches with EMBEDDING_BACKEND=onnx. Alternative model replacement is deferred and fail-closed.

## Decision Inputs
S02 proved the {EXPECTED_ONNX_ARTIFACT} packaged ONNX path with cache namespace {EXPECTED_CACHE_VERSION}, restart/cache proof, legal gate PASS, and sanitized audit. S03 ended with defer_candidate for BAAI/bge-m3 and intfloat/multilingual-e5-large, both failed_closed rather than accepted.

## Same-Host Operating Contract
Use the same-host contract: /health is liveness and runtime metadata, not live inference readiness; readiness requires a smoke embedding request. Redis cache namespace isolation is mandatory; use a disjoint EMBEDDING_CACHE_VERSION and verify runtime.cache_namespace before comparing TEI and ONNX.

## Required Operator Checks
Before switching set EMBEDDING_BACKEND=onnx, ONNX_ARTIFACT_MANIFEST, ONNX_TOKENIZER_PATH, and ONNX_RUNTIME_LIBRARY. Run artifact, tokenizer, and runtime preflight; verify smoke embedding response shape; verify Redis cache namespace isolation via EMBEDDING_CACHE_VERSION.

## Evidence Links
{links}

## Caveats
The legal quality evidence is bounded to the recorded gate and does not prove universal legal quality. ONNX startup must fail closed with no fallback and no silent fallback to TEI inside an ONNX process.

## Non-Goals
Candidate or alternative model replacement is not accepted; it remains deferred/fail-closed. Hosted CI is not a readiness gate and GitHub Actions are not required for same-host runtime readiness.

## Redaction
Raw benchmark/legal text and secrets are excluded; only sanitized markers, hashes, dimensions, counts, and paths to evidence artifacts are referenced.
{extra}
"""


def expect_pass(name: str, text: str) -> None:
    validate_recommendation_text(text)
    print(f"PASS expected-pass {name}")


def expect_fail(name: str, text: str) -> None:
    try:
        validate_recommendation_text(text)
    except ArtifactError as exc:
        print(f"PASS expected-fail {name}: {exc}")
        return
    raise AssertionError(f"expected fixture to fail: {name}")


def run_self_test() -> None:
    good = sample_recommendation()
    with tempfile.TemporaryDirectory() as tmpdir:
        path = Path(tmpdir) / "recommendation.md"
        path.write_text(good, encoding="utf-8")
        validate_recommendation_text(path.read_text(encoding="utf-8"))

    expect_pass("valid", good)
    expect_fail("missing-required-section", good.replace("## Caveats", "## Notes"))
    expect_fail("missing-cache-isolation", good.replace("Redis cache namespace isolation is mandatory; use a disjoint EMBEDDING_CACHE_VERSION and verify runtime.cache_namespace before comparing TEI and ONNX.", "Redis can use the usual shared cache."))
    expect_fail("accepted-candidate-replacement", good.replace("BAAI/bge-m3 and intfloat/multilingual-e5-large, both failed_closed rather than accepted", "BAAI/bge-m3 is accepted and recommended to replace deepvk/USER-bge-m3"))
    expect_fail("hosted-ci-as-gate", good.replace("Hosted CI is not a readiness gate", "Hosted CI is a readiness gate"))
    expect_fail("secret-pattern", sample_recommendation("\napi_key=supersecrettoken12345\n"))
    expect_fail("raw-text-marker", sample_recommendation("\nraw_text: privileged legal probe text\n"))


def parse_args(argv: list[str]) -> argparse.Namespace:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--artifact", type=Path, help="Final M040 S04 recommendation markdown artifact to verify.")
    parser.add_argument("--self-test", action="store_true", help="Run inline positive and negative fixture tests.")
    for name, default in DEFAULT_SOURCE_PATHS.items():
        parser.add_argument(f"--{name.replace('_', '-')}", type=Path, default=default)
    return parser.parse_args(argv)


def main(argv: list[str]) -> int:
    args = parse_args(argv)
    if args.self_test:
        run_self_test()
        return 0
    if not args.artifact:
        raise SystemExit("--artifact is required unless --self-test is used")
    source_paths = {
        "s02_verifier": args.s02_verifier,
        "s03_verifier": args.s03_verifier,
        "s02_benchmark": args.s02_benchmark,
        "s02_preflight": args.s02_preflight,
        "s02_legal": args.s02_legal,
        "s02_audit": args.s02_audit,
        "s03_gate": args.s03_gate,
        "contract": args.contract,
    }
    validate_path(args.artifact, source_paths)
    print(f"M040 S04 recommendation verification: PASS {args.artifact}")
    return 0


if __name__ == "__main__":
    try:
        raise SystemExit(main(sys.argv[1:]))
    except ArtifactError as exc:
        print(f"M040 S04 recommendation verification: FAIL: {exc}", file=sys.stderr)
        raise SystemExit(1)
