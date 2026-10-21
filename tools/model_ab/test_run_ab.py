#!/usr/bin/env python3
"""Unit tests for the model A/B rig. Stdlib only; the rig itself needs requests."""

from __future__ import annotations

import contextlib
import hashlib
import http.server
import json
import math
import subprocess
import sys
import tempfile
import threading
import unittest
from pathlib import Path
from typing import Any, Iterator

HERE = Path(__file__).resolve().parent
sys.path.insert(0, str(HERE))

import run_ab  # noqa: E402


def stub_vector(text: str, salt: str, dimensions: int) -> list[float]:
    digest = hashlib.sha256(f"{salt}\x00{text}".encode("utf-8")).digest()
    values = [digest[index] / 127.5 - 1.0 for index in range(dimensions)]
    norm = math.sqrt(sum(value * value for value in values))
    return [value / norm for value in values]


class StubHandler(http.server.BaseHTTPRequestHandler):
    salt = "stub"
    dimensions = 8
    status = 200
    error_payload: dict[str, Any] | None = None

    def do_POST(self) -> None:
        length = int(self.headers.get("Content-Length", "0"))
        body = json.loads(self.rfile.read(length))
        if self.status != 200:
            payload = json.dumps(self.error_payload or {}).encode("utf-8")
        else:
            data = [
                {
                    "object": "embedding",
                    "index": index,
                    "embedding": stub_vector(text, self.salt, self.dimensions),
                }
                for index, text in enumerate(body["input"])
            ]
            data.reverse()  # the client must restore order from the index field
            payload = json.dumps({"object": "list", "data": data}).encode("utf-8")
        self.send_response(self.status)
        self.send_header("Content-Type", "application/json")
        self.send_header("Content-Length", str(len(payload)))
        self.end_headers()
        self.wfile.write(payload)

    def log_message(self, *_args: Any) -> None:
        pass


def handler_class(**attributes: Any) -> type[StubHandler]:
    return type("StubHandler", (StubHandler,), attributes)


@contextlib.contextmanager
def serve(salt: str = "stub", dimensions: int = 8, **attributes: Any) -> Iterator[str]:
    server = http.server.ThreadingHTTPServer(
        ("127.0.0.1", 0), handler_class(salt=salt, dimensions=dimensions, **attributes)
    )
    thread = threading.Thread(target=server.serve_forever, daemon=True)
    thread.start()
    try:
        yield f"http://127.0.0.1:{server.server_address[1]}"
    finally:
        server.shutdown()
        server.server_close()
        thread.join(timeout=5)


def sample_pairs() -> list[dict[str, Any]]:
    return [
        {
            "id": "p1",
            "query": "SECRET_QUERY_ONE",
            "positive": "SECRET_POSITIVE_ONE",
            "negatives": ["SECRET_NEGATIVE_ONE_A", "SECRET_NEGATIVE_ONE_B"],
        },
        {
            "id": "p2",
            "query": "SECRET_QUERY_TWO",
            "positive": "SECRET_POSITIVE_TWO",
            "negatives": ["SECRET_NEGATIVE_TWO_A"],
        },
        {
            "id": "p3",
            "query": "SECRET_QUERY_THREE",
            "positive": "SECRET_POSITIVE_THREE",
            "negatives": [
                "SECRET_NEGATIVE_THREE_A",
                "SECRET_NEGATIVE_THREE_B",
                "SECRET_NEGATIVE_THREE_C",
            ],
        },
    ]


def document_pool(pairs: list[dict[str, Any]]) -> list[str]:
    documents: list[str] = []
    for pair in pairs:
        for text in [pair["positive"], *pair["negatives"]]:
            if text not in documents:
                documents.append(text)
    return documents


def expected_ranks(pairs: list[dict[str, Any]], documents: list[str], salt: str, top_k: int = 10) -> list[int]:
    document_vectors = [stub_vector(text, salt, 8) for text in documents]
    ranks = []
    for pair in pairs:
        query = stub_vector(pair["query"], salt, 8)
        scores = [run_ab.cosine(query, vector) for vector in document_vectors]
        order = sorted(range(len(scores)), key=lambda index: (-scores[index], index))
        relevant = documents.index(pair["positive"])
        ranks.append(min(order.index(relevant) + 1, top_k + 1))
    return ranks


def expected_metrics(ranks: list[int]) -> dict[str, float]:
    count = len(ranks)
    return {
        "recall_at_1": sum(rank <= 1 for rank in ranks) / count,
        "recall_at_5": sum(rank <= 5 for rank in ranks) / count,
        "recall_at_10": sum(rank <= 10 for rank in ranks) / count,
        "mrr": sum(1.0 / rank for rank in ranks) / count,
        "ndcg_at_10": sum(
            1.0 / math.log2(rank + 1) if rank <= 10 else 0.0 for rank in ranks
        )
        / count,
    }


def run_cli(arguments: list[str]) -> subprocess.CompletedProcess[str]:
    return subprocess.run(
        [sys.executable, str(HERE / "run_ab.py"), *arguments],
        capture_output=True,
        text=True,
        timeout=60,
    )


class RankingMetricTests(unittest.TestCase):
    def test_rank_cutoff_and_ties(self) -> None:
        query = [1.0, 0.0, 0.0]
        documents = [[1.0, 0.0, 0.0], [0.9, 0.1, 0.0], [0.0, 1.0, 0.0], [0.0, 0.0, 1.0]]
        self.assertEqual(run_ab.rank_of(query, documents, 0, 10), 1)
        self.assertEqual(run_ab.rank_of(query, documents, 1, 10), 2)
        self.assertEqual(run_ab.rank_of(query, documents, 3, 3), 4)

        tied = [[1.0, 0.0], [1.0, 0.0]]
        self.assertEqual(run_ab.rank_of([1.0, 0.0], tied, 1, 10), 2)

    def test_summary_formulas(self) -> None:
        metrics = run_ab.summarize_ranks([1, 2, 4])
        self.assertAlmostEqual(metrics["recall_at_1"], 1 / 3)
        self.assertAlmostEqual(metrics["recall_at_5"], 1.0)
        self.assertAlmostEqual(metrics["recall_at_10"], 1.0)
        self.assertAlmostEqual(metrics["mrr"], (1 + 0.5 + 0.25) / 3)
        self.assertAlmostEqual(
            metrics["ndcg_at_10"], (1 + 1 / math.log2(3) + 1 / math.log2(5)) / 3
        )

    def test_non_finite_scores_are_rejected(self) -> None:
        with self.assertRaises(RuntimeError):
            run_ab.rank_of([0.0, 0.0], [[0.0, 0.0]], 0, 10)

    def test_agreement_is_skipped_when_document_prefixes_differ(self) -> None:
        self.assertEqual(
            run_ab.compute_agreement([], [], comparable=False),
            {"status": "skipped", "reason": "document_prefix_mismatch"},
        )


class PairValidationTests(unittest.TestCase):
    def test_positive_listed_as_negative_is_rejected(self) -> None:
        with tempfile.TemporaryDirectory() as tmp:
            path = Path(tmp) / "pairs.jsonl"
            path.write_text(
                json.dumps(
                    {
                        "id": "p1",
                        "query": "SECRET_QUERY",
                        "positive": "SECRET_POSITIVE",
                        "negatives": ["SECRET_POSITIVE"],
                    }
                )
                + "\n",
                encoding="utf-8",
            )
            with self.assertRaises(run_ab.InputError) as caught:
                run_ab.load_pairs(path)
            self.assertIn("also listed as a negative", str(caught.exception))
            self.assertNotIn("SECRET", str(caught.exception))

    def test_missing_negatives_are_rejected_without_echoing_text(self) -> None:
        with tempfile.TemporaryDirectory() as tmp:
            path = Path(tmp) / "pairs.jsonl"
            path.write_text(
                json.dumps(
                    {
                        "id": "p1",
                        "query": "SECRET_QUERY",
                        "positive": "SECRET_POSITIVE",
                        "negatives": [],
                    }
                )
                + "\n",
                encoding="utf-8",
            )
            with self.assertRaises(run_ab.InputError) as caught:
                run_ab.load_pairs(path)
            self.assertIn("negatives", str(caught.exception))
            self.assertNotIn("SECRET", str(caught.exception))


class EndToEndTests(unittest.TestCase):
    def test_report_metrics_and_sanitization(self) -> None:
        pairs = sample_pairs()
        documents = document_pool(pairs)
        with tempfile.TemporaryDirectory() as tmp:
            pairs_path = Path(tmp) / "pairs.jsonl"
            pairs_path.write_text(
                "\n".join(json.dumps(pair, ensure_ascii=False) for pair in pairs) + "\n",
                encoding="utf-8",
            )
            out = Path(tmp) / "report.json"
            with serve("baseline") as baseline_url, serve("candidate") as candidate_url:
                result = run_cli(
                    [
                        "--baseline-url",
                        baseline_url,
                        "--candidate-url",
                        candidate_url,
                        "--pairs",
                        str(pairs_path),
                        "--out",
                        str(out),
                    ]
                )
            self.assertEqual(result.returncode, 0, result.stderr)
            raw = out.read_text(encoding="utf-8")
            report = json.loads(raw)

        for pair in pairs:
            for text in [pair["query"], pair["positive"], *pair["negatives"]]:
                self.assertNotIn(text, raw)
        self.assertNotIn("SECRET", raw)
        self.assertFalse(report["raw_texts_in_report"])
        self.assertEqual(report["pairs"]["count"], 3)
        self.assertEqual(report["pairs"]["documents"], len(documents))
        self.assertEqual(len(report["pairs"]["sha256"]), 64)
        self.assertEqual(report["baseline"]["dimensions"], 8)
        self.assertEqual(report["candidate"]["dimensions"], 8)

        for name, salt in (("baseline", "baseline"), ("candidate", "candidate")):
            ranks = expected_ranks(pairs, documents, salt)
            self.assertEqual([entry["rank"] for entry in report[name]["ranks"]], ranks)
            self.assertEqual(
                [entry["id"] for entry in report[name]["ranks"]],
                [pair["id"] for pair in pairs],
            )
            for metric, value in expected_metrics(ranks).items():
                self.assertLessEqual(
                    abs(report[name]["retrieval"][metric] - value), 1e-6, (name, metric)
                )
            self.assertEqual(report[name]["retrieval"]["queries"], 3)
            self.assertEqual(report[name]["retrieval"]["documents"], len(documents))
            self.assertEqual(report[name]["texts"]["queries"], 3)
            self.assertEqual(report[name]["latency_ms"]["batches"], 2)

        agreement = report["agreement"]
        self.assertEqual(agreement["status"], "ok")
        self.assertEqual(agreement["scope"], "documents")
        self.assertEqual(agreement["count"], len(documents))
        self.assertGreaterEqual(agreement["cosine"]["min"], -1.0)
        self.assertLessEqual(agreement["cosine"]["mean"], 1.0)

    def test_agreement_is_one_for_identical_vectors(self) -> None:
        pairs = sample_pairs()
        with tempfile.TemporaryDirectory() as tmp:
            pairs_path = Path(tmp) / "pairs.jsonl"
            pairs_path.write_text(
                "\n".join(json.dumps(pair, ensure_ascii=False) for pair in pairs) + "\n",
                encoding="utf-8",
            )
            out = Path(tmp) / "report.json"
            with serve("same") as first_url, serve("same") as second_url:
                result = run_cli(
                    [
                        "--baseline-url",
                        first_url,
                        "--candidate-url",
                        second_url,
                        "--pairs",
                        str(pairs_path),
                        "--out",
                        str(out),
                    ]
                )
            self.assertEqual(result.returncode, 0, result.stderr)
            report = json.loads(out.read_text(encoding="utf-8"))
        self.assertEqual(report["agreement"]["status"], "ok")
        self.assertAlmostEqual(report["agreement"]["cosine"]["mean"], 1.0, places=7)
        self.assertAlmostEqual(report["agreement"]["cosine"]["min"], 1.0, places=7)
        self.assertAlmostEqual(report["agreement"]["cosine"]["p1"], 1.0, places=7)
        for metric in ("recall_at_1_pp", "mrr", "ndcg_at_10"):
            self.assertEqual(report["delta_candidate_minus_baseline"][metric], 0.0)

    def test_dimension_mismatch_skips_agreement(self) -> None:
        pairs = sample_pairs()
        with tempfile.TemporaryDirectory() as tmp:
            pairs_path = Path(tmp) / "pairs.jsonl"
            pairs_path.write_text(
                "\n".join(json.dumps(pair, ensure_ascii=False) for pair in pairs) + "\n",
                encoding="utf-8",
            )
            out = Path(tmp) / "report.json"
            with serve("baseline", dimensions=8) as baseline_url, serve(
                "candidate", dimensions=4
            ) as candidate_url:
                result = run_cli(
                    [
                        "--baseline-url",
                        baseline_url,
                        "--candidate-url",
                        candidate_url,
                        "--pairs",
                        str(pairs_path),
                        "--out",
                        str(out),
                    ]
                )
            self.assertEqual(result.returncode, 0, result.stderr)
            report = json.loads(out.read_text(encoding="utf-8"))
        self.assertEqual(report["baseline"]["dimensions"], 8)
        self.assertEqual(report["candidate"]["dimensions"], 4)
        self.assertEqual(report["agreement"]["status"], "skipped")
        self.assertEqual(report["agreement"]["reason"], "dimension_mismatch")
        self.assertEqual(report["candidate"]["retrieval"]["queries"], 3)

    def test_http_error_does_not_echo_body(self) -> None:
        pairs = sample_pairs()
        with tempfile.TemporaryDirectory() as tmp:
            pairs_path = Path(tmp) / "pairs.jsonl"
            pairs_path.write_text(
                "\n".join(json.dumps(pair, ensure_ascii=False) for pair in pairs) + "\n",
                encoding="utf-8",
            )
            out = Path(tmp) / "report.json"
            error = {
                "error": {
                    "code": "invalid_request_error",
                    "message": "upstream echoed SECRET_QUERY_ONE",
                }
            }
            with serve(status=400, error_payload=error) as url:
                result = run_cli(
                    [
                        "--baseline-url",
                        url,
                        "--candidate-url",
                        url,
                        "--pairs",
                        str(pairs_path),
                        "--out",
                        str(out),
                    ]
                )
            self.assertEqual(result.returncode, 2)
            self.assertIn("HTTP 400", result.stderr)
            self.assertIn("invalid_request_error", result.stderr)
            self.assertNotIn("SECRET", result.stderr)
            self.assertFalse(out.exists())


if __name__ == "__main__":
    unittest.main()
