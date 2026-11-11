#!/usr/bin/env bash
# Add the WhatsApp MCP server to Claude Desktop's config, keeping any other servers.
# Usage: ./setup-claude-desktop.sh   (then quit Claude Desktop with Cmd+Q and reopen it)
set -euo pipefail

cd "$(dirname "$0")"
repo="$PWD"

if [[ ! -f .env ]]; then
  echo "Missing .env — run: cp .env.example .env and set API_KEY" >&2
  exit 1
fi
api_key=$(grep '^API_KEY=' .env | cut -d= -f2-)
if [[ -z "$api_key" || "$api_key" == "replace-with-a-random-secret" ]]; then
  echo "API_KEY is not set in .env" >&2
  exit 1
fi

uv_path=$(command -v uv) || { echo "uv not found — run: brew install uv" >&2; exit 1; }

config="$HOME/Library/Application Support/Claude/claude_desktop_config.json"

UV_PATH="$uv_path" REPO="$repo" API_KEY="$api_key" CONFIG="$config" python3 - <<'EOF'
import json, os
p = os.environ["CONFIG"]
os.makedirs(os.path.dirname(p), exist_ok=True)
cfg = {}
if os.path.exists(p) and os.path.getsize(p) > 0:
    with open(p) as f:
        cfg = json.load(f)
cfg.setdefault("mcpServers", {})["whatsapp"] = {
    "command": os.environ["UV_PATH"],
    "args": ["--directory", os.path.join(os.environ["REPO"], "whatsapp-mcp-server"), "run", "main.py"],
    "env": {"API_KEY": os.environ["API_KEY"]},
}
with open(p, "w") as f:
    json.dump(cfg, f, indent=2)
print("Updated", p)
EOF

echo "Done. Quit Claude Desktop (Cmd+Q) and reopen it."
