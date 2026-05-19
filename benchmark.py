#!/usr/bin/env python3
"""
Benchmark: Russian text embedding latency and cache performance.
Runs from the host — API at localhost:8000, Redis at localhost:6379.
"""

import time
import statistics
import requests
import subprocess
import redis

API = "http://localhost:8000"
R = redis.Redis(host="localhost", port=6379, decode_responses=True)

# Russian test texts of varying length
RUSSIAN_TEXTS = {
    "short (11 chars)": "Привет, как дела?",
    "medium (84 chars)": "Москва — столица России, город с богатой историей и культурным наследием.",
    "long (446 chars)": (
        "Искусственный интеллект — это область компьютерных наук, которая занимается созданием "
        "интеллектуальных агентов, способных выполнять задачи, традиционно требующие человеческого "
        "интеллекта. К таким задачам относятся распознавание речи, принятие решений, понимание "
        "естественного языка и визуальное восприятие. Современные языковые модели способны "
        "генерировать связный текст, отвечать на вопросы и даже писать программный код."
    ),
    "very_long (864 chars)": (
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
    payload = {"model": "deepvk/USER-bge-m3", "input": text, "encoding_format": "base64"}
    start = time.perf_counter()
    r = requests.post(f"{API}/v1/embeddings", json=payload, timeout=timeout)
    latency_ms = (time.perf_counter() - start) * 1000
    r.raise_for_status()
    return latency_ms, r.json()


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
    print("Model: deepvk/USER-bge-m3 (1024 dimensions)")
    print("=" * 70)

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
    text = RUSSIAN_TEXTS["medium (84 chars)"]
    flush_cache()
    call_api(text)  # prime

    lats = []
    for _ in range(100):
        lat, _ = call_api(text)
        lats.append(lat)

    print(f"  text: {text[:60]}...")
    print(f"  mean:   {statistics.mean(lats):.2f}ms")
    print(f"  stdev:  {statistics.stdev(lats):.2f}ms")
    print(f"  min:    {min(lats):.2f}ms")
    print(f"  p50:    {statistics.median(lats):.2f}ms")
    print(f"  p95:    {sorted(lats)[int(len(lats)*0.95)]:.2f}ms")
    print(f"  p99:    {sorted(lats)[int(len(lats)*0.99)]:.2f}ms")
    print(f"  max:    {max(lats):.2f}ms")

    # --- Section 3: Throughput with concurrency ---
    print("\n## 3. Throughput — Concurrent Requests (5s window)\n")
    print(f"{'Concurrency':>12} {'Req/s':>8} {'Mean ms':>8} {'p50 ms':>8} {'p95 ms':>8}")
    print("-" * 50)

    import concurrent.futures

    def worker(done_event):
        lats = []
        while not done_event.is_set():
            try:
                lat, _ = call_api(RUSSIAN_TEXTS["medium (84 chars)"], timeout=30)
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
        l2_lat, _ = call_api(l2_text)
        l2_restart_result = {
            "prime_ms": prime_lat,
            "after_restart_ms": l2_lat,
        }
        print(f"  prime/cold request:      {prime_lat:.2f}ms")
        print(f"  after API restart:       {l2_lat:.2f}ms")
        print("  expectation:             served from Redis L2, then backfilled into L1")

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

    print("\n" + "=" * 70)


if __name__ == "__main__":
    main()
