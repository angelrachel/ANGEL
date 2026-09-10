#!/usr/bin/env bash
set -euo pipefail

root_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/../services/fixture-http" && pwd)"
port="${FIXTURE_HTTP_TEST_PORT:-18080}"
log_file="$(mktemp)"

cleanup() {
  if [[ -n "${server_pid:-}" ]]; then
    kill "$server_pid" 2>/dev/null || true
    wait "$server_pid" 2>/dev/null || true
  fi
  rm -f "$log_file"
}
trap cleanup EXIT

(
  cd "$root_dir"
  FIXTURE_HTTP_ADDR="127.0.0.1:${port}" go run . >"$log_file" 2>&1
) &
server_pid=$!

for _ in $(seq 1 20); do
  if curl --fail --silent "http://127.0.0.1:${port}/healthz" >/tmp/angel-fixture-health.json; then
    break
  fi
  sleep 0.25
done

jq -e '.status == "ok" and .service == "fixture-http"' /tmp/angel-fixture-health.json >/dev/null
curl --fail --silent "http://127.0.0.1:${port}/version" | jq -e '.version == "1.0.0"' >/dev/null
curl --fail --silent "http://127.0.0.1:${port}/security-posture" | jq -e '.security_headers.content_type_options == true' >/dev/null
[[ "$(curl --silent --output /dev/null --write-out '%{http_code}' "http://127.0.0.1:${port}/not-registered")" == "404" ]]
rm -f /tmp/angel-fixture-health.json
printf '%s\n' 'P0 fixture HTTP smoke test passed'
