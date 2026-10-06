#!/usr/bin/env bash
#
# restore.sh — MongoDB restore for the local e-commerce stack.
#
# Usage:
#   ./scripts/restore.sh <backup_file>
#
# The backup file must be a .archive.gz archive produced by backup.sh.
# Restoring drops the target database first, so this is destructive — a
# confirmation prompt is shown unless --yes is passed.
#
#   ./scripts/restore.sh --yes backups/ecommerce_20260101_000000.archive.gz
#
# Procedure:
#   1. Validate the backup file exists.
#   2. Confirm intent (destructive).
#   3. mongorestore --drop from the archive.
#   4. Report success/failure.

set -euo pipefail

COMPOSE_FILE="$(cd "$(dirname "$0")/.." && pwd)/docker-compose.yml"
DB_USER="${MONGO_USERNAME:-DRAGON}"
DB_PASSWORD="${MONGO_PASSWORD:-DRAGON}"
DB_NAME="${MONGO_DATABASE:-ECOMMERCE}"

AUTO_YES=0
BACKUP_FILE=""

for arg in "$@"; do
  case "$arg" in
    --yes) AUTO_YES=1 ;;
    *)     BACKUP_FILE="$arg" ;;
  esac
done

if [ -z "$BACKUP_FILE" ]; then
  echo "Usage: $0 [--yes] <backup_file>" >&2
  echo "  backup_file: path to a .archive.gz dump produced by backup.sh" >&2
  exit 1
fi

if [ ! -f "$BACKUP_FILE" ]; then
  echo "ERROR: backup file not found: ${BACKUP_FILE}" >&2
  exit 1
fi

# Validate the dump BEFORE the destructive restore: a truncated/empty/corrupt
# file would wipe the database and restore nothing.
if [ ! -s "$BACKUP_FILE" ] || [ "$(stat -c%s "$BACKUP_FILE")" -lt 1024 ]; then
  echo "ERROR: backup file is empty or too small (< 1 KB); refusing to restore." >&2
  exit 1
fi

if ! gzip -t "$BACKUP_FILE" 2>/dev/null; then
  echo "ERROR: backup file is not a valid gzip archive; refusing to restore." >&2
  exit 1
fi

if ! docker compose -f "${COMPOSE_FILE}" ps mongo >/dev/null 2>&1; then
  echo "ERROR: mongo container is not present in compose." >&2
  exit 1
fi

if [ "$AUTO_YES" -ne 1 ]; then
  read -r -p "Restore will DROP and recreate database '${DB_NAME}'. Continue? [y/N] " answer
  case "$answer" in
    y|Y) ;;
    *) echo "Aborted."; exit 1 ;;
  esac
fi

echo "Restoring ${BACKUP_FILE} -> ${DB_NAME}"

docker compose -f "${COMPOSE_FILE}" exec -T mongo \
  mongorestore -u "${DB_USER}" -p "${DB_PASSWORD}" --authenticationDatabase admin \
  --drop --gzip --archive < "$BACKUP_FILE"

echo "Restore complete: ${DB_NAME} is back in service."
