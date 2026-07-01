---
estimated_steps: 10
estimated_files: 1
skills_used: []
---

# T05: Docker compose runtime reproducer verifying auto-recovery

Создать `tools/verify_warmup_recovery.sh` — bash reproducer доказывающий end-to-end recovery. Скрипт:

1. `docker compose down` (если запущено).
2. Стартовать только fd_api и fd_redis БЕЗ tei (или с tei который сознательно задерживается): `docker compose up -d redis api` + вручную `docker compose up -d tei` через 30 секунд после старта api.
3. Подождать 60 секунд после старта tei.
4. `curl -s localhost:8000/health | jq` — assert status=="ok", warmup_done==true, model_loaded==true.
5. Без recovery скрипт бы показал status=="degraded", warmup_done==false. Скрипт логирует оба ветви для доказательства.
6. Также может воспроизвести "kill tei и restart" сценарий.
7. Артефакт логируется в stdout.

Скрипт должен быть idempotent и safe. ВАЖНО: Скрипт НЕ должен деструктивно менять данные Redis (flush). EXIT 0 если recovery подтверждён (status=ok), EXIT 1 если нет.

Скрипт должен работать с аутентификацией (FD_API_KEY) — берёт токен из docker env контейнера fd_api: `docker exec fd_api printenv FD_API_KEY`.

## Inputs

- `docker-compose.yaml`
- `api/main.go`

## Expected Output

- `tools/verify_warmup_recovery.sh`

## Verification

bash tools/verify_warmup_recovery.sh && echo "RECOVERY OK"

## Observability Impact

Скрипт логирует timeline startup, каждый healthcheck, время до recovery.
