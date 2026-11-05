#!/usr/bin/env bash
set -euo pipefail
umask 077

# Only deploy this project; never prune images or remove volumes.
release_tag=${1:?Usage: deploy.sh sha-COMMIT}
[[ "$release_tag" =~ ^sha-[0-9a-f]{40}$ ]] || { echo 'Invalid release tag' >&2; exit 1; }
project_dir="$HOME/zzhan-blog"
release_dir=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
cd "$project_dir"
test -s .env || { echo 'Missing server .env' >&2; exit 1; }
test -s config.yaml || { echo 'Missing server config.yaml' >&2; exit 1; }
test -s "$release_dir/compose.yaml"
exec 9> .deploy.lock
flock -n 9 || { echo 'Another deployment is running' >&2; exit 1; }
export APP_TAG="$release_tag"

compose() {
  docker compose --project-name zzhan-blog --project-directory "$project_dir" \
    --env-file "$project_dir/.env" -f "$release_dir/compose.yaml" "$@"
}

on_failure() {
  echo 'Deployment failed. Inspect the zzhan-blog containers; no automatic database rollback was attempted.' >&2
  compose ps --all || true
}
trap on_failure ERR
compose config --quiet

echo "Downloading application images: $release_tag"
pulled=false
for attempt in 1 2 3; do
  if compose pull backend frontend; then
    pulled=true
    break
  fi
  echo "Download attempt $attempt failed" >&2
done
[[ "$pulled" == true ]] || { echo 'Images unavailable; running services were not updated' >&2; exit 1; }

echo 'Updating blog services and waiting for health checks'
compose up -d --wait --wait-timeout 300
for port in 3000 4000; do
  curl --fail --silent --show-error --retry 5 --retry-connrefused \
    --connect-timeout 5 --max-time 15 "http://127.0.0.1:$port/" -o /dev/null
  curl --fail --silent --show-error --retry 5 --retry-connrefused \
    --connect-timeout 5 --max-time 15 "http://127.0.0.1:$port/api/v1/health" -o /dev/null
done

# Preserve all secrets; only persist the successfully deployed application tag.
env_tmp=$(mktemp "$project_dir/.env.next.XXXXXX")
awk -v tag="$release_tag" '
  /^APP_TAG=/ { print "APP_TAG=" tag; found=1; next }
  { print }
  END { if (!found) print "APP_TAG=" tag }
' .env > "$env_tmp"
chmod 600 "$env_tmp"
mv -- "$env_tmp" .env
install -m 644 "$release_dir/compose.yaml" "$project_dir/compose.yaml"
printf '%s\n' "$release_tag" > .last-successful-release
compose ps
echo "Deployment succeeded: $release_tag"
