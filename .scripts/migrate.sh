#!/bin/bash

# Use the passed CWD or default to the current directory
CWD=${CWD:-$(pwd)}
cd "$CWD" || exit 1

GOOSE_DRIVER=turso
GOOSE_DBSTRING=file:$CWD/.local/local.db
GOOSE_MIGRATION_DIR=$CWD/server/db/migrations

usage() {
  echo "Usage: $0 [command]"
  echo ""
  echo "Commands:"
  echo "  up       - Apply all available migrations"
  echo "  up-by-one - Apply the next available migration"
  echo "  down     - Roll back all migrations"
  echo "  down-by-one - Roll back the most recent migration"
  echo "  status   - Show migration status"
  echo "  create [name] - Create a new migration"
  echo ""
  exit 1
}

if [ "$#" -eq 0 ]; then
  usage
fi

case "$1" in
  up)
    GOOSE_DRIVER=$GOOSE_DRIVER GOOSE_DBSTRING=$GOOSE_DBSTRING GOOSE_MIGRATION_DIR=$GOOSE_MIGRATION_DIR goose up
    ;;
  up-by-one)
    GOOSE_DRIVER=$GOOSE_DRIVER GOOSE_DBSTRING=$GOOSE_DBSTRING GOOSE_MIGRATION_DIR=$GOOSE_MIGRATION_DIR goose up-by-one
    ;;
  down)
    GOOSE_DRIVER=$GOOSE_DRIVER GOOSE_DBSTRING=$GOOSE_DBSTRING GOOSE_MIGRATION_DIR=$GOOSE_MIGRATION_DIR goose down-to 0
    ;;
  down-by-one)
    GOOSE_DRIVER=$GOOSE_DRIVER GOOSE_DBSTRING=$GOOSE_DBSTRING GOOSE_MIGRATION_DIR=$GOOSE_MIGRATION_DIR goose down
    ;;
  status)
    GOOSE_DRIVER=$GOOSE_DRIVER GOOSE_DBSTRING=$GOOSE_DBSTRING GOOSE_MIGRATION_DIR=$GOOSE_MIGRATION_DIR goose status
    ;;
  create)
    if [ "$#" -lt 2 ]; then
      echo "Error: Missing migration name"
      echo "Usage: $0 create <migration_name>"
      exit 1
    fi
    GOOSE_DRIVER=$GOOSE_DRIVER GOOSE_DBSTRING=$GOOSE_DBSTRING GOOSE_MIGRATION_DIR=$GOOSE_MIGRATION_DIR goose create "$2" sql
    ;;
  *)
    usage
    ;;
esac
