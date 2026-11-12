# Upstream Sync

UPSTREAM_COMMIT=7d6a06dcdce1f01dfb24f60e1030d5efba9f3b88
UPSTREAM_DATE=2026-04-21

## Sync Log
- 2026-04-21: Initial fork. Removed send tools, added auth middleware, added upstream check.
- 2026-04-21: Updated whatsmeow to v0.0.0-20260421083005-5b8886176ff7 (WhatsApp rejected old client version 405). Fixed context API changes. Added -phone flag for phone-number pairing.
- 2026-11-10: Updated whatsmeow to v0.0.0-20261110152934-c88388f983b4 (and transitive deps). No API changes needed; build, vet and tests pass. Added weekly GitHub Action that opens a PR when whatsmeow has a new version.
- 2026-11-12: Media downloads broke after the 2026-11-10 whatsmeow update (every download returned 403). whatsmeow's Download() now uses only the direct path and appends "&hash=...", so the path must keep its signed query string; extractDirectPathFromURL was stripping it. Fixed by keeping the query string.
