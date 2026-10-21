# Model A/B rig (issue #480, item 7)

Status: experiment tooling. Nothing in this directory changes `api/`,
`docker-compose.yaml`, the default model, or production behavior.

The rig answers one question: on our own Russian retrieval data, does a candidate
embedding model beat or match `deepvk/USER-bge-m3` closely enough to justify the
migration cost? Public leaderboard scores are context only. They are not
evidence for this service.

## What the rig measures

`run_ab.py` sends the same JSONL pair file to two OpenAI-compatible endpoints
(`POST /v1/embeddings` with `input=[...]`, batched) and writes a sanitized JSON
report:

- **Recall@1, Recall@5, Recall@10** per model. Each query has exactly one
  relevant document, so Recall@k is the share of queries whose positive is in
  the top k.
- **MRR**: mean of `1 / rank` of the positive.
- **nDCG@10**: binary gain, one relevant document, ideal DCG equals 1.
- **Agreement**: cosine between the baseline and candidate vectors for the same
  document texts (mean, min, 1st percentile). It is skipped when the native
  widths differ or when the two sides use different document prefixes.
- **Latency**: end-to-end p50 and p95 of one embedding batch, plus the batch
  count. Latency is not throughput and not per-token time.

Every query is ranked against the whole document pool built from all positives
and negatives in the file. Negatives from other pairs act as ordinary
distractors, so a pair's own negatives stay hard negatives. `--top-k` caps the
ranking window: a positive found below it is scored as rank `top_k + 1`.

## Candidates from issue #480

| Model | Native dims | Context | License | TEI 1.9.4 support | Main blocker |
|---|---:|---:|---|---|---|
| `deepvk/USER-bge-m3` (baseline) | 1024 | 8192 | Apache-2.0 | current production image | Russian-only; no change needed |
| `Qwen/Qwen3-Embedding-0.6B` | 1024 (MRL 32-1024) | 32K | Apache-2.0 | Qwen3 is on the supported-model list | CPU cost; historical TEI#667 |
| `jinaai/jina-embeddings-v5-text-small` | 1024 (MRL 32-1024) | 32K | CC-BY-NC-4.0 | not listed among supported architectures | non-commercial license |
| `Snowflake/snowflake-arctic-embed-m-v2.0` | 768 (MRL 256) | 8192 | Apache-2.0 | Alibaba GTE architecture is supported | 768 dims change the storage/API contract |
| `google/embeddinggemma-300m` | 768 (MRL 512/256/128) | 2048 | Gemma Terms | Gemma3 is supported, model is gated | short context and gated download |

The TEI support column follows the supported-model list for the 1.9.x image and
the v1.9.0 release notes. Re-check it against the pinned image before a run.

The speed tier is `snowflake-arctic-embed-m-v2.0` and `embeddinggemma-300m`:
305M and 300M parameters respectively. The quality candidates are Qwen3-0.6B
and Jina v5 small.

### Blockers in detail

- **Jina v5 small license.** The model is CC-BY-NC-4.0; commercial use requires
  a separate license. That is a product blocker regardless of retrieval
  quality, and TEI does not list this architecture as supported. Do not spend a
  decision run on it until both points are resolved.
- **TEI#667.** Qwen3-Embedding-0.6B failed to start on CPU with TEI 1.7.3: no
  ONNX file in the repo, then an MKL error in the Candle backend. The issue is
  closed, TEI 1.9.0 added bidirectional attention and projection-layer support
  for Qwen3-based models, and current TEI docs list Qwen3 as supported. Re-test
  with the pinned `cpu-1.9.4` image before repeating the blocker.
- **Speed-tier context.** `embeddinggemma-300m` accepts 2048 tokens. Long legal
  fragments will be truncated on the way in, which changes the comparison for
  that model and makes chunking policy part of its evaluation. Also, the model
  is gated: TEI needs `HF_TOKEN` to download it.
- **768 dimensions.** fd accepts 512 or 1024 only. A 768-dim candidate must be
  compared TEI-to-TEI, and adopting it requires a storage and API contract
  change outside this slice.
- **MRL does not make the encoder faster.** Requesting 256 or 512 dimensions
  truncates the output after the full forward pass. It reduces vector storage
  and index size, not encoder latency or model memory. Do not report MRL as an
  encoder speedup.

## Decision gate

Run the candidate and the baseline on the same host, with the same batch size,
the same `--top-k`, and the same pair-file hash. The report records all of these.

Proposed defaults for a promote/reject decision:

| Gate | Threshold |
|---|---|
| Recall@10 delta | `>= -1.0` percentage points |
| MRR delta | `>= -0.01` absolute |
| p95 batch latency | `<= 1.10 x` baseline at equal batch size |
| Steady-state RSS | `<= baseline + 1 GiB`, with headroom on the host |
| License | allows production use |
| TEI | loads the pinned revision on `cpu-1.9.4` |

A candidate that only ties on quality needs a concrete runtime or licensing win
to justify a migration. Change these thresholds in the issue before a decision
run, not after seeing the results.

**The example pair file is a smoke test, not a benchmark.** Sixteen synthetic
pairs cannot support a one-percentage-point claim. The real gate needs labeled
Russian retrieval data (the 44-ФЗ corpus and query set used by
`tools/evaluate_legal_retrieval.py`, or another reviewed set) and enough queries
that the delta is measurable.

## Quick start

Pin the candidate by model id and revision:

```bash
export MODEL_AB_MODEL_ID=Qwen/Qwen3-Embedding-0.6B
export MODEL_AB_REVISION=<commit-sha>
```

The candidate mounts the production Hugging Face cache as an external volume.
Check its name first; the default is `fd_tei_data`:

```bash
docker volume ls --format '{{.Name}}' | grep tei_data
# if it does not exist yet:
docker volume create fd_tei_data
# if production uses another name:
export MODEL_AB_HF_CACHE_VOLUME=<volume-name>
```

Start only the candidate. The production stack is not touched:

```bash
docker compose -f tools/model_ab/docker-compose.ab.yaml up -d
curl -fsS http://localhost:30081/health
```

The container runs the same `cpu-1.9.4` image, the same `--max-batch-tokens 8192`
limit, and the same HF cache as production. TEI 1.9.x enables `--auto-truncate`
by default; inputs above the model or batch limit are truncated on both sides.

Choose the baseline:

- `http://localhost:30080` is production TEI. Use it for both quality and
  latency comparisons.
- `http://localhost:8000` is the fd API. It is fine for quality, but its Redis
  cache distorts latency: a document cached by an earlier run returns without
  inference.
- For a 768-dim candidate, use TEI on both sides. Do not pass `--dimensions`
  unless both models support MRL truncation to that width.

Run the comparison:

```bash
python3 tools/model_ab/run_ab.py \
  --baseline-url http://localhost:30080 \
  --candidate-url http://localhost:30081 \
  --pairs tools/model_ab/pairs.example.jsonl \
  --out .gsd/runtime/model-ab/report.json \
  --batch-size 16 \
  --warmup
```

`--baseline-model` and `--candidate-model` are optional. Pass them when the
server expects a model id; TEI only validates the name when it is started with
`--served-model-name`.

`--warmup` sends one extra request per endpoint and drops it from the latency
sample. The first request after readiness otherwise includes model warm-up.

If the endpoint requires a bearer token, pass `--api-key "$FD_API_KEY"`. The key
is never written to the report.

Instruction-aware models need their query prefix. The flag applies to queries
only:

```bash
  --candidate-query-prefix 'query: '                       # E5-style
  --candidate-query-prefix 'task: search result | query: ' # EmbeddingGemma
  --candidate-document-prefix 'title: none | text: '       # EmbeddingGemma documents
```

If the two document prefixes differ, agreement is reported as
`document_prefix_mismatch`; compare retrieval metrics only. For Qwen3, pass the
model-card instruction and its trailing `Query:` marker as the query prefix.

Stop the candidate when finished:

```bash
docker compose -f tools/model_ab/docker-compose.ab.yaml down
```

The HF cache volume is external, so `down` and `down -v` leave it in place.

## Report shape

The report uses stable ordering with sorted keys and contains no raw texts:

```json
{
  "pairs": {"file": "<name>", "count": "<int>", "documents": "<int>", "sha256": "<hex>"},
  "config": {"baseline_url": "<url>", "candidate_url": "<url>", "batch_size": "<int>", "...": "..."},
  "baseline": {
    "retrieval": {"recall_at_1": "<float>", "recall_at_5": "<float>", "recall_at_10": "<float>",
                  "mrr": "<float>", "ndcg_at_10": "<float>", "queries": "<int>",
                  "documents": "<int>", "top_k": "<int>"},
    "latency_ms": {"p50": "<float>", "p95": "<float>", "batches": "<int>"},
    "dimensions": "<int>",
    "ranks": [{"id": "<pair-id>", "rank": "<int>"}],
    "embeddings_sha256": "<hex>"
  },
  "candidate": {"...": "same shape"},
  "delta_candidate_minus_baseline": {"recall_at_10_pp": "<float>", "mrr": "<float>", "...": "..."},
  "agreement": {"status": "ok", "scope": "documents", "count": "<int>",
                "cosine": {"mean": "<float>", "min": "<float>", "p1": "<float>"}},
  "raw_texts_in_report": false
}
```

`ranks` exposes one rank per pair id, so a failing query can be found without
storing its text. `embeddings_sha256` hashes the float32 vectors in a fixed
order.

## Pair file

One JSON object per line:

```json
{"id": "ru-example", "query": "...", "positive": "...", "negatives": ["...", "..."]}
```

Rules: unique non-empty `id`, non-empty `query` and `positive`, at least one
non-empty negative. Missing fields fail with the file name, line number, and
field name; the error does not echo the text.

`pairs.example.jsonl` contains 16 synthetic RU/EN pairs and 64 unique documents.
It contains no legal corpus and no personal data. Replace it for a real gate.

## Sanitization

- The report contains counts, hashes, metrics, and latency only.
- Error messages name rows and fields, never their contents.
- HTTP errors expose only the OpenAI `error.code` or `error.type`, never the
  response body.
- Hashes fingerprint the pair file and the vectors. They do not turn
  low-entropy secrets into safe data; do not put secrets in the pair file.

## INT8 S8S8-QDQ parity

The goal is to learn whether a static INT8 artifact keeps dense retrieval
quality on CPU, not to claim a speedup before measuring one.

S8S8-QDQ means signed int8 activations and signed int8 weights, with
QuantizeLinear/DequantizeLinear nodes and static activation scales from
calibration. It is the variant worth testing on x86 with VNNI. U8 activation
variants and dynamic quantization are different artifacts; do not mix them into
this gate.

### Produce the artifacts

Use the existing export tool, which writes a dense pooled graph with dynamic
sequence axes:

```bash
python3 tools/export_user_bge_m3_dense_onnx.py \
  --model-path <local-snapshot> \
  --output-dir .gsd/runtime/model-ab/onnx-fp32 \
  --sequence-length 512
```

Quantize that graph with ONNX Runtime static QDQ. Keep tokenizer, pooling, and
normalization identical to the exported graph; quantize only the exported dense
model:

```python
from onnxruntime.quantization import QuantFormat, QuantType, quantize_static

quantize_static(
    reference_path,
    int8_path,
    calibration_reader,
    quant_format=QuantFormat.QDQ,
    activation_type=QuantType.QInt8,
    weight_type=QuantType.QInt8,
    per_channel=True,
)
```

The calibration reader must stream representative Russian texts from the
retrieval domain. Record the ONNX Runtime version, calibration set hash, and
export settings next to the artifact. The parity script compares numerics; it
does not prove that the artifact is QDQ, so check the quantization log or the
graph in Netron before trusting a pass.

### Check parity

Run this in the ONNX export environment (numpy, onnxruntime, transformers), not
in the fd Go environment:

```bash
python3 tools/model_ab/check_int8_parity.py \
  --reference-onnx .gsd/runtime/model-ab/onnx-fp32/model.onnx \
  --int8-onnx .gsd/runtime/model-ab/onnx-int8/model.onnx \
  --model-path <local-snapshot> \
  --pairs .gsd/runtime/model-ab/pairs.legal.jsonl \
  --output .gsd/runtime/model-ab/int8-parity.json
```

Acceptance:

| Check | Threshold |
|---|---|
| Mean cosine, reference vs INT8, same texts | `>= 0.999` |
| Minimum cosine | `>= 0.995` |
| Recall@10 drop | `<= 1.0` percentage point |

The script exits with code 2 when a check fails. A pass only admits the artifact
to the API-level Russian retrieval gate; it is not production proof and does not
establish a latency win.

### Artifact policy

Do not commit `.onnx` files, external ONNX data, model snapshots, calibration
sets, or reports. Keep them under `.gsd/runtime/` (already ignored) or under
`tools/model_ab/artifacts/` (ignored by the local `.gitignore`). CI does not
download models for this directory.

## Tests

```bash
python3 -m unittest discover -s tools/model_ab -p 'test_*.py'
```

The tests run a local `http.server` stub that returns deterministic vectors from
the SHA-256 of each text, including out-of-order `data` items. They check:

- rank calculation, ties, the `top_k + 1` cutoff, Recall/MRR/nDCG formulas;
- end-to-end metrics and the agreement block;
- dimension mismatch and document-prefix mismatch behavior;
- that raw texts and response bodies never appear in the report or in errors.

## Limitations

- The example pairs are synthetic and small. There is one relevant document per
  query, no graded relevance, and no confidence intervals. A real decision needs
  a labeled set large enough for the delta to mean something.
- The document pool is the pair file itself, not a production corpus, and no
  chunking is applied. Chunking policy is part of the real gate, especially for
  2048-token candidates.
- Latency is batch end-to-end on one host under current load. It is not
  throughput and it is not comparable across batch sizes.
- TEI 1.9.x truncates automatically at the model or batch-token limit. Check the
  truncated-text rate for candidates with short context; the rig does not
  report it.
- The rig uses the plain `/v1/embeddings` contract. It does not exercise Redis,
  fd caches, retries, or the fd chunking path.
