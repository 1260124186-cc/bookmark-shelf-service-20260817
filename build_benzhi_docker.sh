#!/usr/bin/env sh
set -eu

image_name="${1:-bookmark-shelf-service:local}"
platform="${2:-linux/arm64}"

docker build --platform "$platform" -f benzhi.Dockerfile -t "$image_name" .
container_id="$(docker run -d -p 8080:8080 "$image_name")"
trap 'docker rm -f "$container_id" >/dev/null 2>&1 || true' EXIT

docker exec "$container_id" go build ./...
curl --fail --retry 10 --retry-connrefused http://localhost:8080/healthz
curl --fail http://localhost:8080/report
