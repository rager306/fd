#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
IMAGE_TAG="${IMAGE_TAG:-fd-api:onnx1024-local}"
CONTEXT_DIR="${CONTEXT_DIR:-$ROOT_DIR/.gsd/runtime/docker/onnx1024-context}"
ONNX_MANIFEST="${ONNX_MANIFEST:-docs/onnx-artifacts/user-bge-m3-dense-fp32.json}"
NATIVE_TOKENIZER_MANIFEST="${NATIVE_TOKENIZER_MANIFEST:-docs/onnx-artifacts/hf-tokenizers-linux-amd64.json}"
TOKENIZER_JSON="${TOKENIZER_JSON:-tei-models/deepvk--USER-bge-m3/tokenizer.json}"
ONNX_RUNTIME_LIBRARY="${ONNX_RUNTIME_LIBRARY:-/root/.cache/uv/archive-v0/tj7L7fjW7RE4nnYdVfYZ1/lib/python3.13/site-packages/onnxruntime/capi/libonnxruntime.so.1.26.0}"

cd "$ROOT_DIR"

python3 tools/verify_onnx_artifacts.py \
  --onnx-manifest "$ONNX_MANIFEST" \
  --native-tokenizer-manifest "$NATIVE_TOKENIZER_MANIFEST" >/tmp/fd-onnx-artifacts-verify.json
python3 -m json.tool /tmp/fd-onnx-artifacts-verify.json >/dev/null

if [[ ! -f "$TOKENIZER_JSON" ]]; then
  echo "missing tokenizer JSON: $TOKENIZER_JSON" >&2
  exit 1
fi
if [[ ! -f "$ONNX_RUNTIME_LIBRARY" ]]; then
  echo "missing ONNX Runtime shared library: $ONNX_RUNTIME_LIBRARY" >&2
  exit 1
fi

rm -rf "$CONTEXT_DIR"
mkdir -p "$CONTEXT_DIR"

cp -a api "$CONTEXT_DIR/api"
mkdir -p "$CONTEXT_DIR/docs/onnx-artifacts"
cp docs/onnx-artifacts/user-bge-m3-dense-fp32.json "$CONTEXT_DIR/docs/onnx-artifacts/"
cp docs/onnx-artifacts/hf-tokenizers-linux-amd64.json "$CONTEXT_DIR/docs/onnx-artifacts/"
cp docs/onnx-artifacts/README.md "$CONTEXT_DIR/docs/onnx-artifacts/"

mkdir -p "$CONTEXT_DIR/.gsd/runtime/onnx/m010-s03"
cp .gsd/runtime/onnx/m010-s03/user-bge-m3-dense.onnx "$CONTEXT_DIR/.gsd/runtime/onnx/m010-s03/"

mkdir -p "$CONTEXT_DIR/.gsd/runtime/tokenizers/linux-amd64"
cp .gsd/runtime/tokenizers/linux-amd64/libtokenizers.a "$CONTEXT_DIR/.gsd/runtime/tokenizers/linux-amd64/"

mkdir -p "$CONTEXT_DIR/tokenizer"
cp "$TOKENIZER_JSON" "$CONTEXT_DIR/tokenizer/tokenizer.json"

mkdir -p "$CONTEXT_DIR/onnxruntime"
cp "$ONNX_RUNTIME_LIBRARY" "$CONTEXT_DIR/onnxruntime/libonnxruntime.so.1.26.0"

cp Dockerfile.onnx "$CONTEXT_DIR/Dockerfile.onnx"

docker build -f "$CONTEXT_DIR/Dockerfile.onnx" -t "$IMAGE_TAG" "$CONTEXT_DIR"

echo "built_image=$IMAGE_TAG"
echo "context_dir=$CONTEXT_DIR"
