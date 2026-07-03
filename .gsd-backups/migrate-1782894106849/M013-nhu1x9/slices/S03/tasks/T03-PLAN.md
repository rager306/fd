---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T03: Run tagged ONNX cosine comparison

Start tagged ONNX API locally using validated native tokenizer library and isolated Redis cache namespace, then rerun TEI-vs-ONNX cosine comparison. If startup fails, capture blocker evidence.

## Inputs

- `benchmark-results/fd-go-onnx-m011-s03.txt`
- `benchmark-results/fd-tokenizer-tagged-native-m013-s02.txt`

## Expected Output

- `benchmark-results/fd-go-onnx-hf-tokenizer-m013-s03.txt`

## Verification

Cosine artifact passes threshold or records startup/runtime blocker; Redis namespace isolated.

## Observability Impact

Produces the first semantically meaningful Go ONNX cosine evidence if tagged runtime succeeds.
