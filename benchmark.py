#!/usr/bin/env python3
"""
Benchmark: Russian text embedding latency and cache performance.
Runs from the host — API at localhost:8000, Redis at localhost:6379.
"""

import hashlib
import json
import os
from pathlib import Path
import statistics
import subprocess
import time

import redis
import requests

API = os.getenv("BENCHMARK_API_URL", "http://localhost:8000")
MODEL = os.getenv("BENCHMARK_MODEL", "deepvk/USER-bge-m3")
DIMENSIONS = int(os.getenv("BENCHMARK_DIMENSIONS", "1024"))
ENVIRONMENT_BASELINE = Path(os.getenv("BENCHMARK_ENVIRONMENT_BASELINE", "benchmark-results/fd-environment-inxi-m008.txt"))
SNAPSHOT_VERSION = 1
R = redis.Redis(host=os.getenv("BENCHMARK_REDIS_HOST", "localhost"), port=int(os.getenv("BENCHMARK_REDIS_PORT", "6379")), decode_responses=True)

SAFE_ENV_KEYS = [
    "BENCHMARK_API_URL",
    "BENCHMARK_MODEL",
    "BENCHMARK_DIMENSIONS",
    "BENCHMARK_REDIS_HOST",
    "BENCHMARK_REDIS_PORT",
    "BENCHMARK_ENVIRONMENT_BASELINE",
    "PORT",
    "MODEL_ID",
    "LOG_LEVEL",
    "REDIS_HOST",
    "REDIS_PORT",
    "REDIS_POOL_SIZE",
    "REDIS_CACHE_TTL",
    "REDIS_CACHE_NO_EXPIRE",
    "REDIS_MAXMEMORY",
    "REDIS_MAXMEMORY_POLICY",
    "REDIS_RDB_SAVE",
    "REDIS_AOF_ENABLED",
    "EMBEDDING_MODEL_ID",
    "EMBEDDING_MODEL_REVISION",
    "EMBEDDING_CACHE_VERSION",
    "EMBEDDING_TOKENIZER_VERSION",
    "EMBEDDING_CHUNKING_VERSION",
    "LOCAL_CACHE_TTL",
    "LOCAL_CACHE_MAX_SIZE",
    "BATCH_CACHE_MODE",
    "BATCH_CACHE_PIPELINE_SIZE",
]
SECRET_KEY_PARTS = ("SECRET", "TOKEN", "PASSWORD", "PASS", "KEY", "CREDENTIAL", "AUTH", "COOKIE")


def is_secret_like(key: str) -> bool:
    upper_key = key.upper()
    return any(part in upper_key for part in SECRET_KEY_PARTS)


def run_metadata_command(args: list[str], timeout=10):
    try:
        result = subprocess.run(
            args,
            check=True,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True,
            timeout=timeout,
        )
        return result.stdout.strip() or None
    except (FileNotFoundError, subprocess.CalledProcessError, subprocess.TimeoutExpired):
        return None


def sha256_text(value: str) -> str:
    return hashlib.sha256(value.encode("utf-8")).hexdigest()


def sha256_file(path: Path):
    if not path.exists() or not path.is_file():
        return None
    digest = hashlib.sha256()
    with path.open("rb") as handle:
        for chunk in iter(lambda: handle.read(1024 * 1024), b""):
            digest.update(chunk)
    return digest.hexdigest()


def collect_git_metadata():
    commit = run_metadata_command(["git", "rev-parse", "HEAD"])
    branch = run_metadata_command(["git", "branch", "--show-current"])
    status = run_metadata_command(["git", "status", "--short"])
    return {
        "commit": commit,
        "branch": branch,
        "dirty": bool(status),
    }


def collect_compose_metadata():
    compose_config = run_metadata_command(["docker", "compose", "config"], timeout=20)
    images_output = run_metadata_command(["docker", "compose", "images"], timeout=20)
    images = images_output.splitlines() if images_output else []
    return {
        "config_sha256": sha256_text(compose_config) if compose_config else None,
        "images_output_sha256": sha256_text(images_output) if images_output else None,
        "images_available": images_output is not None,
        "images": images,
    }


def collect_safe_environment():
    values = {}
    omitted_secret_like = []
    for key in SAFE_ENV_KEYS:
        if key not in os.environ:
            continue
        if is_secret_like(key):
            omitted_secret_like.append(key)
            continue
        values[key] = os.environ[key]
    return {
        "values": values,
        "omitted_secret_like_keys": sorted(omitted_secret_like),
        "allowlist_size": len(SAFE_ENV_KEYS),
    }


def collect_redis_metadata():
    try:
        memory = R.info("memory")
        stats = R.info("stats")
        server = R.info("server")
        keyspace = R.info("keyspace")
        redis_config = {}
        for key in ("maxmemory", "maxmemory-policy", "save", "appendonly"):
            value = R.config_get(key).get(key)
            redis_config[key] = value
        return {
            "available": True,
            "redis_version": server.get("redis_version"),
            "config": redis_config,
            "used_memory_human": memory.get("used_memory_human"),
            "used_memory_peak_human": memory.get("used_memory_peak_human"),
            "mem_fragmentation_ratio": memory.get("mem_fragmentation_ratio"),
            "keyspace_hits": stats.get("keyspace_hits"),
            "keyspace_misses": stats.get("keyspace_misses"),
            "evicted_keys": stats.get("evicted_keys"),
            "expired_keys": stats.get("expired_keys"),
            "db0_keys": keyspace.get("db0", {}).get("keys") if isinstance(keyspace.get("db0"), dict) else None,
        }
    except redis.RedisError as err:
        return {
            "available": False,
            "error": type(err).__name__,
        }


def collect_environment_baseline_metadata():
    return {
        "path": str(ENVIRONMENT_BASELINE),
        "exists": ENVIRONMENT_BASELINE.exists(),
        "sha256": sha256_file(ENVIRONMENT_BASELINE),
    }


def effective_config_snapshot():
    return {
        "snapshot_version": SNAPSHOT_VERSION,
        "benchmark": {
            "script": "benchmark.py",
            "api_url": API,
            "model": MODEL,
            "dimensions": DIMENSIONS,
            "input_texts_logged": False,
        },
        "git": collect_git_metadata(),
        "docker_compose": collect_compose_metadata(),
        "environment": collect_safe_environment(),
        "redis_before_run": collect_redis_metadata(),
        "environment_baseline": collect_environment_baseline_metadata(),
        "redaction_policy": {
            "mode": "allowlist_with_secret_like_omission",
            "secret_key_parts": list(SECRET_KEY_PARTS),
            "raw_benchmark_texts_excluded": True,
        },
    }


def print_effective_config_snapshot():
    print("\n## 0. Effective Configuration Snapshot (sanitized)\n")
    print(json.dumps(effective_config_snapshot(), ensure_ascii=False, indent=2, sort_keys=True))


RUSSIAN_TEXTS = {
    "short (17 chars)": "Привет, как дела?",
    "medium (73 chars)": "Москва — столица России, город с богатой историей и культурным наследием.",
    "long (422 chars)": (
        "Искусственный интеллект — это область компьютерных наук, которая занимается созданием "
        "интеллектуальных агентов, способных выполнять задачи, традиционно требующие человеческого "
        "интеллекта. К таким задачам относятся распознавание речи, принятие решений, понимание "
        "естественного языка и визуальное восприятие. Современные языковые модели способны "
        "генерировать связный текст, отвечать на вопросы и даже писать программный код."
    ),
    "very_long (693 chars)": (
        "В современном мире технологии машинного обучения и глубокого обучения играют ключевую роль "
        "в развитии искусственного интеллекта. Нейронные сети с миллионами параметров обучаются на "
        "огромных корпусах текстов, что позволяет им понимать контекст, семантику и даже эмоциональную "
        "окраску сообщений. Модели типа трансформеров, такие как BERT и GPT, произвели революцию в "
        "обработке естественного языка. Они используются в поисковых системах, голосовых помощниках, "
        "системах машинного перевода и многих других приложениях. Важным направлением является "
        "обучение с подкреплением на основе обратной связи от людей, что позволяет моделям лучше "
        "соответствовать человеческим предпочтениям и этическим нормам."
    ),
}


def flush_cache():
    R.flushall()


def call_api(text: str, timeout=120):
    payload = {"model": MODEL, "input": text, "encoding_format": "base64"}
    start = time.perf_counter()
    r = requests.post(f"{API}/v1/embeddings", json=payload, timeout=timeout)
    latency_ms = (time.perf_counter() - start) * 1000
    r.raise_for_status()
    return latency_ms, r.json()


def call_batch_api(inputs: list[str], timeout=120):
    payload = {
        "inputs": inputs,
        "dimensions": DIMENSIONS,
        "encoding_format": "base64",
    }
    start = time.perf_counter()
    r = requests.post(f"{API}/embeddings/batch", json=payload, timeout=timeout)
    latency_ms = (time.perf_counter() - start) * 1000
    r.raise_for_status()
    return latency_ms, r.json()


def percentile(values: list[float], pct: float) -> float:
    if not values:
        return 0.0
    idx = min(len(values) - 1, int(len(values) * pct))
    return sorted(values)[idx]


def print_latency_summary(prefix: str, lats: list[float]):
    print(f"  {prefix}mean:   {statistics.mean(lats):.2f}ms")
    print(f"  {prefix}stdev:  {statistics.stdev(lats):.2f}ms" if len(lats) > 1 else f"  {prefix}stdev:  0.00ms")
    print(f"  {prefix}min:    {min(lats):.2f}ms")
    print(f"  {prefix}p50:    {statistics.median(lats):.2f}ms")
    print(f"  {prefix}p95:    {percentile(lats, 0.95):.2f}ms")
    print(f"  {prefix}p99:    {percentile(lats, 0.99):.2f}ms")
    print(f"  {prefix}max:    {max(lats):.2f}ms")


def redis_stats_snapshot():
    try:
        stats = R.info("stats")
        keyspace = R.info("keyspace")
        return {
            "keyspace_hits": int(stats.get("keyspace_hits", 0)),
            "keyspace_misses": int(stats.get("keyspace_misses", 0)),
            "evicted_keys": int(stats.get("evicted_keys", 0)),
            "expired_keys": int(stats.get("expired_keys", 0)),
            "total_commands_processed": int(stats.get("total_commands_processed", 0)),
            "db0_keys": int(keyspace.get("db0", {}).get("keys", 0)) if isinstance(keyspace.get("db0"), dict) else 0,
        }
    except redis.RedisError:
        return None


def redis_stats_delta(before, after):
    if before is None or after is None:
        return None
    return {key: after.get(key, 0) - before.get(key, 0) for key in before}


def print_redis_delta(label: str, delta):
    if delta is None:
        print(f"  redis_delta/{label}: unavailable")
        return
    print(f"  redis_delta/{label}: hits={delta['keyspace_hits']} misses={delta['keyspace_misses']} "
          f"evicted={delta['evicted_keys']} expired={delta['expired_keys']} "
          f"commands={delta['total_commands_processed']} db0_key_delta={delta['db0_keys']}")


def wait_for_api(timeout=60):
    deadline = time.perf_counter() + timeout
    last_error = ""
    while time.perf_counter() < deadline:
        try:
            r = requests.get(f"{API}/health", timeout=2)
            r.raise_for_status()
            return True, ""
        except Exception as err:
            last_error = str(err)
            time.sleep(1)
    return False, last_error or "timed out waiting for API health"


def restart_api_for_l2_check():
    try:
        subprocess.run(
            ["docker", "compose", "restart", "api"],
            check=True,
            stdout=subprocess.PIPE,
            stderr=subprocess.STDOUT,
            text=True,
            timeout=120,
        )
    except FileNotFoundError:
        return False, "docker command not found"
    except subprocess.TimeoutExpired:
        return False, "docker compose restart api timed out"
    except subprocess.CalledProcessError as err:
        output = (err.stdout or "").strip().splitlines()
        detail = output[-1] if output else str(err)
        return False, f"docker compose restart api failed: {detail}"
    return wait_for_api(timeout=60)


def main():
    print("=" * 70)
    print("EMBEDDING BENCHMARK — Russian Text / Redis Cache")
    print(f"Model: {MODEL} ({DIMENSIONS} dimensions)")
    print("=" * 70)
    print_effective_config_snapshot()

    # Warm up the service
    print("\n[warmup]...", end=" ", flush=True)
    call_api("Разогрев")
    print("OK")

    # --- Section 1: Single-request latency, cold vs warm ---
    print("\n## 1. Single-Request Latency (cold vs warm, 5 runs each)\n")
    print(f"{'Label':<25} {'Chars':>5}  {'COLD mean':>12} {'WARM mean':>10} {'Speedup':>8}")
    print("-" * 65)

    rows = []
    for label, text in RUSSIAN_TEXTS.items():
        # Cold: flush, then measure
        flush_cache()
        cold_lats = []
        for _ in range(5):
            lat, _ = call_api(text)
            cold_lats.append(lat)

        # Warm: prime cache, then measure
        call_api(text)
        warm_lats = []
        for _ in range(5):
            lat, _ = call_api(text)
            warm_lats.append(lat)

        cold_mean = statistics.mean(cold_lats)
        warm_mean = statistics.mean(warm_lats)
        speedup = cold_mean / warm_mean if warm_mean > 0 else float("inf")

        rows.append({
            "label": label,
            "chars": len(text),
            "cold_mean": round(cold_mean, 1),
            "warm_mean": round(warm_mean, 2),
            "speedup": round(speedup, 1),
            "cold_stdev": round(statistics.stdev(cold_lats), 1) if len(cold_lats) > 1 else 0,
            "warm_stdev": round(statistics.stdev(warm_lats), 1) if len(warm_lats) > 1 else 0,
        })

        print(f"{label:<25} {len(text):>5}  {cold_mean:>10.1f}ms  {warm_mean:>8.2f}ms  {speedup:>6.1f}x")

    # --- Section 2: 100 repeated requests (cache behavior) ---
    print(f"\n## 2. Repeated Requests — Cache Hit Behavior (100 requests)\n")
    text = RUSSIAN_TEXTS["medium (73 chars)"]
    flush_cache()
    call_api(text)  # prime

    redis_before_l1 = redis_stats_snapshot()
    lats = []
    for _ in range(100):
        lat, _ = call_api(text)
        lats.append(lat)
    redis_after_l1 = redis_stats_snapshot()

    print("  text_label: medium (73 chars)")
    print(f"  chars: {len(text)}")
    print_latency_summary("", lats)
    print_redis_delta("l1_hot_repeated", redis_stats_delta(redis_before_l1, redis_after_l1))

    # --- Section 3: Throughput with concurrency ---
    print("\n## 3. Throughput — Concurrent Requests (5s window)\n")
    print(f"{'Concurrency':>12} {'Req/s':>8} {'Mean ms':>8} {'p50 ms':>8} {'p95 ms':>8}")
    print("-" * 50)

    import concurrent.futures

    def worker(done_event):
        lats = []
        while not done_event.is_set():
            try:
                lat, _ = call_api(RUSSIAN_TEXTS["medium (73 chars)"], timeout=30)
                lats.append(lat)
            except Exception:
                break
        return lats

    throughput_rows = []
    for conc in [1, 4, 8, 16]:
        done_event = __import__("threading").Event()
        all_lats = []
        with concurrent.futures.ThreadPoolExecutor(max_workers=conc) as ex:
            futures = [ex.submit(worker, done_event) for _ in range(conc)]
            time.sleep(5)
            done_event.set()
            for f in futures:
                all_lats.extend(f.result())

        rps = len(all_lats) / 5
        throughput_rows.append({"concurrency": conc, "rps": rps})
        print(f"{conc:>12}  {rps:>7.1f}  {statistics.mean(all_lats):>7.1f}ms  "
              f"{statistics.median(all_lats):>7.1f}ms  {sorted(all_lats)[int(len(all_lats)*0.95)]:>7.1f}ms")

    # --- Section 4: Response format verification ---
    print("\n## 4. Response Format Verification\n")
    lat, resp = call_api("Тест")
    emb = resp["data"][0]["embedding"]
    print(f"  dimensions:  {len(emb)}")
    print(f"  model:       {resp.get('model', 'N/A')}")
    print(f"  object:      {resp.get('object', 'N/A')}")
    print(f"  usage/prompt_tokens: {resp.get('usage', {}).get('prompt_tokens', 'N/A')}")

    # --- Section 5: Redis L2 persistence after API restart ---
    print("\n## 5. Redis L2 Persistence — After API Restart\n")
    l2_restart_result = None
    l2_text = "Redis L2 persistence diagnostic text"
    flush_cache()
    prime_lat, _ = call_api(l2_text)
    ok, reason = restart_api_for_l2_check()
    if not ok:
        print(f"  skipped: {reason}")
    else:
        redis_before_l2 = redis_stats_snapshot()
        l2_lat, _ = call_api(l2_text)
        redis_after_l2 = redis_stats_snapshot()
        l2_restart_result = {
            "prime_ms": prime_lat,
            "after_restart_ms": l2_lat,
        }
        print(f"  prime/cold request:      {prime_lat:.2f}ms")
        print(f"  after API restart:       {l2_lat:.2f}ms")
        print("  expectation:             served from Redis L2, then backfilled into L1")
        print_redis_delta("l2_after_api_restart", redis_stats_delta(redis_before_l2, redis_after_l2))

    # --- Section 6: Cached batch endpoint behavior ---
    print("\n## 6. Cached Batch Endpoint — L1 and Redis L2\n")
    batch_inputs = [f"batch-cache-item-{i}" for i in range(16)]

    flush_cache()
    prime_batch_lat, prime_batch_resp = call_batch_api(batch_inputs)
    print(f"  batch_size:               {len(batch_inputs)}")
    print(f"  prime/cold batch:         {prime_batch_lat:.2f}ms")
    print(f"  prime/count:              {prime_batch_resp.get('count', 'N/A')}")

    redis_before_batch_l1 = redis_stats_snapshot()
    batch_l1_lats = []
    for _ in range(20):
        lat, _ = call_batch_api(batch_inputs)
        batch_l1_lats.append(lat)
    redis_after_batch_l1 = redis_stats_snapshot()
    print("  l1_hot_batch:")
    print_latency_summary("  ", batch_l1_lats)
    print_redis_delta("batch_l1_hot", redis_stats_delta(redis_before_batch_l1, redis_after_batch_l1))

    ok, reason = restart_api_for_l2_check()
    batch_l2_lats = []
    if not ok:
        print(f"  redis_l2_batch_after_api_restart: skipped: {reason}")
    else:
        redis_before_batch_l2 = redis_stats_snapshot()
        for _ in range(10):
            lat, _ = call_batch_api(batch_inputs)
            batch_l2_lats.append(lat)
        redis_after_batch_l2 = redis_stats_snapshot()
        print("  redis_l2_batch_after_api_restart:")
        print_latency_summary("  ", batch_l2_lats)
        print_redis_delta("batch_l2_after_api_restart", redis_stats_delta(redis_before_batch_l2, redis_after_batch_l2))

    # --- Section 7: Repeated chunk reuse pattern ---
    print("\n## 7. Repeated Chunk Reuse Pattern\n")
    chunk_pool = [f"research-chunk-{i}" for i in range(8)]
    repeated_batches = [chunk_pool[i % len(chunk_pool)] for i in range(32)]
    flush_cache()
    redis_before_chunks = redis_stats_snapshot()
    chunk_lats = []
    for _ in range(8):
        lat, _ = call_batch_api(repeated_batches)
        chunk_lats.append(lat)
    redis_after_chunks = redis_stats_snapshot()
    print(f"  unique_chunks:            {len(chunk_pool)}")
    print(f"  batch_items_per_round:    {len(repeated_batches)}")
    print(f"  rounds:                   {len(chunk_lats)}")
    print(f"  first_cold_round:         {chunk_lats[0]:.2f}ms")
    warm_chunk_lats = chunk_lats[1:]
    print("  warm_reuse_rounds:")
    print_latency_summary("  ", warm_chunk_lats)
    print_redis_delta("repeated_chunk_reuse", redis_stats_delta(redis_before_chunks, redis_after_chunks))

    # --- Summary ---
    print("\n" + "=" * 70)
    print("SUMMARY")
    print("=" * 70)

    # Best cold latency
    best_cold = min(rows, key=lambda x: x["cold_mean"])
    print(f"\n  Best cold latency:  {best_cold['cold_mean']}ms ({best_cold['label']})")
    print(f"  Cache speedup:      {max(r['speedup'] for r in rows):.1f}x (median text)")
    print(f"  Warm latency mean:  {statistics.mean([r['warm_mean'] for r in rows]):.2f}ms")
    max_throughput = max(throughput_rows, key=lambda x: x["rps"])
    print(f"  Max throughput:     ~{max_throughput['rps']:.0f} req/s ({max_throughput['concurrency']} concurrent)")
    if l2_restart_result is None:
        print("  Redis L2 restart:   skipped")
    else:
        print(f"  Redis L2 restart:   {l2_restart_result['after_restart_ms']:.2f}ms after API restart")
    print(f"  Batch L1 p95:       {percentile(batch_l1_lats, 0.95):.2f}ms ({len(batch_inputs)} items)")
    if batch_l2_lats:
        print(f"  Batch L2 p95:       {percentile(batch_l2_lats, 0.95):.2f}ms after API restart")
    else:
        print("  Batch L2 p95:       skipped")
    print(f"  Chunk reuse warm p95: {percentile(warm_chunk_lats, 0.95):.2f}ms")

    print("\n" + "=" * 70)


if __name__ == "__main__":
    main()
