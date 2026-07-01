#!/usr/bin/env bash
#
# verify_warmup_recovery.sh — runtime reproducer for the fd-api warmup
# auto-recovery contract (M051-h1xr44 S01/T05).
#
# Reproduces the docker-compose startup race: fd-api starts before TEI
# finishes BERT load (~15-20s) and would normally leave /health in
# `degraded` forever without manual `docker restart fd_api`. With the
# recovery contract enabled, the api self-recovers once TEI becomes
# reachable, even after the original 5x5s startup window expires.
#
# Effects on the running environment:
#   - Stops the `tei` compose service for up to 4 minutes (TEI BERT load
#     on a loaded CPU can exceed 2 minutes before /ready)
#   - Restarts the `api` compose service at the beginning
#   - Restarts `tei` after the simulated race
#   - Leaves both services in their final recovered state
#
# Prereqs: docker, docker compose v2, jq, curl. Must be run from the
# project root (where docker-compose.yaml lives) or invoked with an
# explicit COMPOSE_FILE path.
#
# Exit codes:
#   0 — recovery confirmed (status=ok, model_loaded=true)
#   1 — recovery did NOT happen; /health remained degraded
#
set -euo pipefail

# ----- Configuration --------------------------------------------------------
COMPOSE_FILE="${COMPOSE_FILE:-docker-compose.yaml}"
RECOVERY_INTERVAL_SEC="${RECOVERY_INTERVAL_SEC:-30}"
MAX_RECOVERY_WAIT_SEC=240        # cap the wait for /health=ok (TEI BERT load on loaded CPU can exceed 2 min)
TEI_STOP_WAIT_SEC=5              # grace for port release
API_RESTART_WAIT_SEC=15           # give api time to come back up and fail warmup

# ----- Output helpers -------------------------------------------------------
log() {
    printf '%s [%s] %s\n' "$(date -u '+%Y-%m-%dT%H:%M:%SZ')" "${1:-info}" "${2}"
}

die() {
    log "error" "$1"
    exit "${2:-1}"
}

# ----- Preconditions --------------------------------------------------------
command -v docker >/dev/null || die "docker not found in PATH"
command -v jq     >/dev/null || die "jq not found in PATH"
command -v curl   >/dev/null || die "curl not found in PATH"

[ -f "$COMPOSE_FILE" ] || die "compose file not found: $COMPOSE_FILE"

log "info" "composing file: $COMPOSE_FILE"

# ----- Health probe helpers -------------------------------------------------
health_payload() {
    curl -fsS --max-time 5 http://localhost:8000/health
}

status_field() {
    health_payload | jq -r '.status' 2>/dev/null || echo "unknown"
}

warmup_done_field() {
    health_payload | jq -r '.warmup_done' 2>/dev/null || echo "unknown"
}

model_loaded_field() {
    health_payload | jq -r '.model_loaded' 2>/dev/null || echo "unknown"
}

wait_for_status() {
    local target="$1" max_wait="$2"
    local elapsed=0
    while [ "$elapsed" -lt "$max_wait" ]; do
        local current
        current="$(status_field)"
        if [ "$current" = "$target" ]; then
            return 0
        fi
        sleep 2
        elapsed=$((elapsed + 2))
    done
    return 1
}

# ----- Pre-test sanity ------------------------------------------------------
log "info" "checking current /health"
initial_status="$(status_field)"
log "info" "current status: $initial_status"
if [ "$initial_status" != "ok" ]; then
    die "pre-flight /health=$initial_status (want ok). Start the stack first: docker compose up -d" 2
fi

# ----- Step 1: stop tei -----------------------------------------------------
log "info" "step 1/5: stopping tei"
docker compose -f "$COMPOSE_FILE" stop tei
sleep "$TEI_STOP_WAIT_SEC"

# ----- Step 2: restart api (forces a fresh cold start against absent tei) --
log "info" "step 2/5: restarting api"
docker compose -f "$COMPOSE_FILE" up -d --no-deps api
sleep "$API_RESTART_WAIT_SEC"
log "info" "api container status: $(docker inspect fd_api --format '{{.State.Status}} health={{.State.Health.Status}}' 2>/dev/null || echo unknown)"

post_restart_status="$(status_field)"
log "info" "post-restart status: $post_restart_status (expect degraded, warmup_done=false)"
if [ "$post_restart_status" = "ok" ]; then
    log "warn" "api reached ok unexpectedly fast; the race may not have reproduced"
fi

# ----- Step 3: capture degraded evidence ------------------------------------
log "info" "step 3/5: capturing degraded state"
degraded_evidence="$(health_payload || echo 'curl_failed')"
echo "$degraded_evidence" | jq . >/dev/null 2>&1 && {
    echo "DEGRADED_RESPONSE = $degraded_evidence" | jq .
} || {
    log "warn" "could not parse degraded response as JSON: $degraded_evidence"
}

# ----- Step 4: re-start tei and wait for recovery --------------------------
log "info" "step 4/5: starting tei (race window ~${RECOVERY_INTERVAL_SEC}s + tei load)"
docker compose -f "$COMPOSE_FILE" up -d --no-deps tei
tei_start_ts=$(date +%s)

# Use the longer of: recovery interval + TEI load buffer, OR a lower bound
# of 120s so we never give up before any retry could fire even if tei is
# still loading BERT on a loaded CPU (observed ~2-3 min under contention).
recovery_budget=$(( RECOVERY_INTERVAL_SEC + 90 ))
if [ "$recovery_budget" -lt 120 ]; then
    recovery_budget=120
fi
if [ "$recovery_budget" -gt "$MAX_RECOVERY_WAIT_SEC" ]; then
    recovery_budget=$MAX_RECOVERY_WAIT_SEC
fi
log "info" "step 5/5: waiting up to ${recovery_budget}s for /health=ok"

if wait_for_status "ok" "$recovery_budget"; then
    recovery_ts=$(date +%s)
    elapsed=$(( recovery_ts - tei_start_ts ))
    log "info" "RECOVERY CONFIRMED after ${elapsed}s since tei start"
    health_payload | jq .
    exit 0
fi

log "error" "recovery did NOT happen within ${recovery_budget}s"
log "error" "final /health payload:"
health_payload | jq . || health_payload
exit 1
