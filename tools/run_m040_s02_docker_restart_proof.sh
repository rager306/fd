#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

IMAGE_TAG="${IMAGE_TAG:-fd-api:onnx-1024}"
API_CONTAINER="${API_CONTAINER:-fd-m040-s02-onnx-api}"
REDIS_CONTAINER="${REDIS_CONTAINER:-fd-m040-s02-redis}"
DOCKER_NETWORK="${DOCKER_NETWORK:-fd-m040-s02-proof}"
API_PORT="${API_PORT:-18000}"
REDIS_PORT="${REDIS_PORT:-16379}"
CACHE_VERSION="m040-s02-onnx-restart"
RESULTS_DIR="benchmark-results"
PREFLIGHT_ARTIFACT="$RESULTS_DIR/fd-m040-s02-onnx-docker-preflight.txt"
BENCHMARK_ARTIFACT="$RESULTS_DIR/fd-benchmark-m040-s02-onnx-docker-restart.txt"
BENCHMARK_API_URL="http://127.0.0.1:${API_PORT}"
RESTART_COMMAND="docker restart ${API_CONTAINER}"

mkdir -p "$RESULTS_DIR"
: >"$PREFLIGHT_ARTIFACT"

log() {
  printf '[%s] %s\n' "$(date -u '+%Y-%m-%dT%H:%M:%SZ')" "$*" | tee -a "$PREFLIGHT_ARTIFACT"
}

record_blocker() {
  log "BLOCKER: $*"
}

require_command() {
  if ! command -v "$1" >/dev/null 2>&1; then
    record_blocker "missing required command: $1"
    exit 1
  fi
}

port_available() {
  local port="$1"
  python3 - "$port" <<'PY'
import socket
import sys
port = int(sys.argv[1])
with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as sock:
    sock.settimeout(0.2)
    sys.exit(1 if sock.connect_ex(("127.0.0.1", port)) == 0 else 0)
PY
}

wait_for_http_health() {
  local url="$1"
  local deadline=$((SECONDS + 90))
  local last_status=""
  while (( SECONDS < deadline )); do
    if last_status="$(curl -fsS --max-time 2 "$url/health" 2>&1)"; then
      return 0
    fi
    sleep 1
  done
  printf '%s\n' "$last_status"
  return 1
}

json_get() {
  python3 - "$1" <<'PY'
import json
import sys
expr = sys.argv[1].split('.')
data = json.load(sys.stdin)
for part in expr:
    data = data[part]
print(data)
PY
}

cleanup_api() {
  docker rm -f "$API_CONTAINER" >/dev/null 2>&1 || true
}

require_command docker
require_command curl
require_command python3

log "phase=preflight image=$IMAGE_TAG api_container=$API_CONTAINER redis_container=$REDIS_CONTAINER api_url=$BENCHMARK_API_URL redis_host_port=127.0.0.1:$REDIS_PORT cache_version=$CACHE_VERSION"

if ! docker info >/dev/null 2>&1; then
  record_blocker "Docker daemon is unavailable or current user lacks Docker permission"
  exit 1
fi

if docker ps --format '{{.Names}}' | grep -Fxq "$API_CONTAINER"; then
  record_blocker "API proof container already running: $API_CONTAINER"
  exit 1
fi
# Remove stale stopped proof API containers so the fixed proof name remains reusable.
docker rm -f "$API_CONTAINER" >/dev/null 2>&1 || true

if ! port_available "$API_PORT"; then
  record_blocker "API proof port 127.0.0.1:$API_PORT is already occupied"
  exit 1
fi

if docker ps --format '{{.Names}}' | grep -Fxq "$REDIS_CONTAINER"; then
  log "phase=redis status=reusing_existing_container name=$REDIS_CONTAINER"
else
  # A stopped proof Redis container is not useful evidence; recreate it with the
  # localhost-only port binding and keep it alive across API restarts.
  docker rm -f "$REDIS_CONTAINER" >/dev/null 2>&1 || true
  if ! port_available "$REDIS_PORT"; then
    record_blocker "Redis proof port 127.0.0.1:$REDIS_PORT is already occupied and $REDIS_CONTAINER is not running"
    exit 1
  fi
  log "phase=redis status=starting name=$REDIS_CONTAINER"
  docker network create "$DOCKER_NETWORK" >/dev/null 2>&1 || true
  docker run -d \
    --name "$REDIS_CONTAINER" \
    --network "$DOCKER_NETWORK" \
    -p "127.0.0.1:${REDIS_PORT}:6379" \
    redis:7-alpine \
    redis-server --save '300 1' --appendonly no >/dev/null
fi

docker network create "$DOCKER_NETWORK" >/dev/null 2>&1 || true
if ! docker network inspect "$DOCKER_NETWORK" --format '{{range .Containers}}{{.Name}}{{"\n"}}{{end}}' | grep -Fxq "$REDIS_CONTAINER"; then
  docker network connect "$DOCKER_NETWORK" "$REDIS_CONTAINER" >/dev/null 2>&1 || true
fi

log "phase=image status=build_or_reuse tag=$IMAGE_TAG"
IMAGE_TAG="$IMAGE_TAG" tools/build_onnx_image.sh 2>&1 | tee -a "$PREFLIGHT_ARTIFACT"

trap cleanup_api EXIT
cleanup_api

log "phase=api status=starting name=$API_CONTAINER bind=127.0.0.1:$API_PORT"
docker run -d \
  --name "$API_CONTAINER" \
  --network "$DOCKER_NETWORK" \
  -p "127.0.0.1:${API_PORT}:${API_PORT}" \
  -e "PORT=${API_PORT}" \
  -e "BIND_HOST=0.0.0.0" \
  -e "MODEL_ID=deepvk/USER-bge-m3" \
  -e "LOG_LEVEL=info" \
  -e "REDIS_HOST=${REDIS_CONTAINER}:6379" \
  -e "REDIS_POOL_SIZE=50" \
  -e "REDIS_CACHE_NO_EXPIRE=true" \
  -e "EMBEDDING_CACHE_VERSION=${CACHE_VERSION}" \
  -e "EMBEDDING_MODEL_ID=deepvk/USER-bge-m3" \
  -e "EMBEDDING_MODEL_REVISION=0cc6cfe48e260fb0474c753087a69369e88709ae" \
  -e "EMBEDDING_TOKENIZER_VERSION=068d9f7ed9dd190a00a567e5f7750fdc591b93bc623072ac8050a303c25f5937" \
  -e "EMBEDDING_CHUNKING_VERSION=m040-s02" \
  "$IMAGE_TAG" >/dev/null

if ! health_json="$(wait_for_http_health "$BENCHMARK_API_URL")"; then
  record_blocker "API /health did not become ready: $health_json"
  docker logs --tail 80 "$API_CONTAINER" >>"$PREFLIGHT_ARTIFACT" 2>&1 || true
  exit 1
fi

log "phase=health status=ready"
{
  echo "HEALTH_JSON_BEGIN"
  printf '%s\n' "$health_json" | python3 -m json.tool
  echo "HEALTH_JSON_END"
} | tee -a "$PREFLIGHT_ARTIFACT"

backend="$(printf '%s\n' "$health_json" | json_get runtime.backend)"
dimensions="$(printf '%s\n' "$health_json" | json_get runtime.dimensions)"
cache_namespace="$(printf '%s\n' "$health_json" | json_get runtime.cache_namespace)"
if [[ "$backend" != "onnx" || "$dimensions" != "1024" || "$cache_namespace" != *"$CACHE_VERSION"* ]]; then
  record_blocker "unexpected /health runtime metadata backend=$backend dimensions=$dimensions cache_namespace=$cache_namespace"
  exit 1
fi

log "phase=smoke status=requesting_embedding"
smoke_json="$(curl -fsS --max-time 120 \
  -H 'Content-Type: application/json' \
  -d '{"model":"deepvk/USER-bge-m3","input":"S02 packaged ONNX smoke probe","encoding_format":"base64"}' \
  "$BENCHMARK_API_URL/v1/embeddings")"
{
  echo "SMOKE_JSON_BEGIN"
  printf '%s\n' "$smoke_json" | python3 -m json.tool
  echo "SMOKE_JSON_END"
} | tee -a "$PREFLIGHT_ARTIFACT"

printf '%s\n' "$smoke_json" | python3 -c '
import json
import sys
body = json.load(sys.stdin)
assert body.get("object") == "list", body
assert body.get("data") and body["data"][0].get("object") == "embedding", body
embedding = body["data"][0].get("embedding")
assert isinstance(embedding, str) and embedding, "expected base64 embedding string"
assert body.get("model") == "deepvk/USER-bge-m3", body.get("model")
'

log "phase=container status=recording_metadata"
{
  echo "CONTAINER_JSON_BEGIN"
  docker inspect "$API_CONTAINER" | python3 -m json.tool
  echo "CONTAINER_JSON_END"
} | tee -a "$PREFLIGHT_ARTIFACT" >/dev/null

log "phase=benchmark status=running artifact=$BENCHMARK_ARTIFACT restart_command=$RESTART_COMMAND"
BENCHMARK_API_URL="$BENCHMARK_API_URL" \
BENCHMARK_REDIS_HOST="127.0.0.1" \
BENCHMARK_REDIS_PORT="$REDIS_PORT" \
BENCHMARK_RUNTIME_LABEL="onnx-docker-m040-s02" \
BENCHMARK_BUILD_TAGS="onnx,hf_tokenizers,docker" \
BENCHMARK_ONNX_ARTIFACT_MANIFEST="docs/onnx-artifacts/user-bge-m3-dense-fp32.json" \
BENCHMARK_NATIVE_TOKENIZER_MANIFEST="docs/onnx-artifacts/hf-tokenizers-linux-amd64.json" \
BENCHMARK_ONNX_RUNTIME_LIBRARY="/opt/onnxruntime/libonnxruntime.so.1.26.0" \
BENCHMARK_API_RESTART_COMMAND="$RESTART_COMMAND" \
EMBEDDING_BACKEND="onnx" \
EMBEDDING_CACHE_VERSION="$CACHE_VERSION" \
EMBEDDING_MODEL_ID="deepvk/USER-bge-m3" \
EMBEDDING_MODEL_REVISION="0cc6cfe48e260fb0474c753087a69369e88709ae" \
EMBEDDING_TOKENIZER_VERSION="068d9f7ed9dd190a00a567e5f7750fdc591b93bc623072ac8050a303c25f5937" \
EMBEDDING_CHUNKING_VERSION="m040-s02" \
python3 benchmark.py 2>&1 | tee "$BENCHMARK_ARTIFACT"

log "phase=verify status=running"
python3 tools/verify_m040_s02_artifacts.py \
  --preflight "$PREFLIGHT_ARTIFACT" \
  --benchmark "$BENCHMARK_ARTIFACT"
log "phase=complete status=pass benchmark_artifact=$BENCHMARK_ARTIFACT preflight_artifact=$PREFLIGHT_ARTIFACT"
