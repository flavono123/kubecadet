#!/usr/bin/env bash
set -euo pipefail

COUNT="${1:-20}"
TARGET="${TARGET_URL:-http://localhost:8080/inventory}"

for i in $(seq 1 "$COUNT"); do
  curl -fsS "$TARGET" >/dev/null && echo "[$(date +%H:%M:%S)] hit $i" || echo "[$(date +%H:%M:%S)] error on request $i" >&2
  sleep 0.2
done
