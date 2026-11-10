#!/usr/bin/env bash
# Start the WhatsApp bridge with API_KEY loaded from .env.
# Usage: ./start-bridge.sh [-phone 15551234567]
set -euo pipefail

cd "$(dirname "$0")"

if [[ ! -f .env ]]; then
  echo "Missing .env — run: cp .env.example .env and set API_KEY" >&2
  exit 1
fi

set -a
source .env
set +a

if [[ -z "${API_KEY:-}" || "$API_KEY" == "replace-with-a-random-secret" ]]; then
  echo "API_KEY is not set in .env" >&2
  exit 1
fi

cd whatsapp-bridge
exec go run . "$@"
