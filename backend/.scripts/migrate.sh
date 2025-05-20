#!/bin/bash

# Use the passed CWD or default to the current directory
CWD=${CWD:-$(pwd)}
cd "$CWD" || exit 1

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
    goose up
    ;;
  up-by-one)
    goose up-by-one
    ;;
  down)
    goose down-to 0
    ;;
  down-by-one)
    goose down
    ;;
  status)
    goose status
    ;;
  create)
    if [ "$#" -lt 2 ]; then
      echo "Error: Missing migration name"
      echo "Usage: $0 create <migration_name>"
      exit 1
    fi
    goose create "$2" sql
    ;;
  *)
    usage
    ;;
esac
