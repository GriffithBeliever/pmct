#!/bin/sh
set -e

# Railway injects PORT for the public-facing port; nginx listens on it.
# Go runs on 8081 and Next.js on 3000 internally.
NGINX_PORT="${PORT:-8080}"

# Render nginx config with the correct port
export NGINX_PORT
envsubst '${NGINX_PORT}' < /etc/nginx/nginx.conf.template > /etc/nginx/nginx.conf

# Start Go backend on internal port 8081
PORT=8081 ./server &
echo "Go backend started on :8081"

# Start Next.js on internal port 3000
PORT=3000 node ./frontend/server.js &
echo "Next.js started on :3000"

# Give services a moment to bind
sleep 2

# Run nginx in the foreground (keeps the container alive)
echo "nginx listening on :${NGINX_PORT}"
exec nginx -g "daemon off;"
