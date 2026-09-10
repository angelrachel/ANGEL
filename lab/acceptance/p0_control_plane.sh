#!/usr/bin/env bash
set -euo pipefail

root_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
port="${ANGEL_CONTROL_TEST_PORT:-18081}"
token="p0-local-operator-key"
log_file="$(mktemp)"
binary_file="$(mktemp)"

cleanup() {
  if [[ -n "${server_pid:-}" ]]; then
    kill "$server_pid" 2>/dev/null || true
    wait "$server_pid" 2>/dev/null || true
  fi
  rm -f "$log_file"
  rm -f "$binary_file"
}
trap cleanup EXIT

(
  cd "$root_dir"
  go build -o "$binary_file" ./src/assessment
)
ANGEL_OPERATOR_KEY="$token" ANGEL_HOST=127.0.0.1 ANGEL_PORT="$port" "$binary_file" >"$log_file" 2>&1 &
server_pid=$!

for _ in $(seq 1 30); do
  if curl --fail --silent "http://127.0.0.1:${port}/healthz" >/tmp/angel-control-health.json; then
    break
  fi
  sleep 0.25
done

jq -e '.status == "ok"' /tmp/angel-control-health.json >/dev/null
headers=(-H "Authorization: Bearer ${token}" -H 'Content-Type: application/json')
curl --fail --silent "${headers[@]}" -X POST "http://127.0.0.1:${port}/api/v1/register" --data '{"id":"worker-p0","hostname":"fixture-worker","os":"linux"}' | jq -e '.status == "ok"' >/dev/null
task_response="$(curl --fail --silent "${headers[@]}" -X POST "http://127.0.0.1:${port}/api/v1/tasks" --data '{"id":"task-p0-1","agent_id":"worker-p0","type":"http-posture","target_ref":"fixture://lab/web-app-01","mode":"observe","requested_by":"operator","payload":{"path":"/security-posture"}}')"
if ! jq -e '.status == "queued" and .target_ref == "fixture://lab/web-app-01"' <<<"$task_response" >/dev/null; then
  printf 'unexpected task response: %s\n' "$task_response" >&2
  exit 1
fi
task_list="$(curl --fail --silent "${headers[@]}" "http://127.0.0.1:${port}/api/v1/tasks?agent_id=worker-p0")"
if ! jq -e 'length == 1 and .[0].id == "task-p0-1" and .[0].status == "queued"' <<<"$task_list" >/dev/null; then
  printf 'unexpected task list: %s\n' "$task_list" >&2
  exit 1
fi

if [[ "$(curl --silent --output /dev/null --write-out '%{http_code}' "${headers[@]}" -X POST "http://127.0.0.1:${port}/api/v1/tasks" --data '{"id":"task-p0-2","agent_id":"worker-p0","type":"http-posture","target_ref":"https://outside.invalid","mode":"observe","requested_by":"operator"}')" != "403" ]]; then
  echo 'expected external target to be rejected' >&2
  exit 1
fi

rm -f /tmp/angel-control-health.json
printf '%s\n' 'P0 control-plane acceptance test passed'
