#!/bin/sh
set -e

NGINX_PORT="${PORT:-8080}"

export NGINX_PORT
envsubst '${NGINX_PORT}' < /etc/nginx/nginx.conf.template > /etc/nginx/nginx.conf

# Start Go backend on internal port 8081
PORT=8081 ./server &
GO_PID=$!
echo "Go backend starting (pid $GO_PID) on :8081"

# Start Next.js on internal port 3000
PORT=3000 node ./frontend/server.js &
NEXT_PID=$!
echo "Next.js starting (pid $NEXT_PID) on :3000"

# Wait for a port to accept connections (timeout in seconds)
wait_for_port() {
  port=$1
  name=$2
  timeout=60
  elapsed=0
  echo "Waiting for $name on :$port ..."
  while ! nc -w1 127.0.0.1 "$port" </dev/null >/dev/null 2>&1; do
    elapsed=$((elapsed + 1))
    if [ "$elapsed" -ge "$timeout" ]; then
      echo "ERROR: $name failed to start on :$port after ${timeout}s" >&2
      exit 1
    fi
    sleep 1
  done
  echo "$name ready on :$port"
}

wait_for_port 8081 "Go backend"
wait_for_port 3000 "Next.js"

echo "All services ready — nginx listening on :${NGINX_PORT}"
exec nginx -g "daemon off;"
