#!/bin/bash

LOCAL_TURSO_DB_PORT=9000
LOCAL_TURSO_DB_URL=http://localhost

echo "Starting a local db instance for unit/integration tests..."
echo "HOST=$LOCAL_TURSO_DB_URL"
echo "PORT=$LOCAL_TURSO_DB_PORT"

turso dev -p $LOCAL_TURSO_DB_PORT > /dev/null 2>&1 &

# get the bg running turso pid
LOCAL_TURSO_DB_PID=$!
echo "PID=$LOCAL_TURSO_DB_PID"
echo "Wait for local db instance to boot..."
sleep 2

echo "Running migrations..."
make up DB_URL="$LOCAL_TURSO_DB_URL:$LOCAL_TURSO_DB_PORT"
# make up DB_URL="sqlite://./test.db"

echo "Run tests..."
TURSO_DATABASE_URL="$LOCAL_TURSO_DB_URL:$LOCAL_TURSO_DB_PORT" go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
# TURSO_DATABASE_URL="file:./test.db" go test -v ./...

if ps -p $LOCAL_TURSO_DB_PID > /dev/null; then
    echo "Terminating local db instance..."
    kill -SIGINT $LOCAL_TURSO_DB_PID

    if ! ps -p $LOCAL_TURSO_DB_PID > /dev/null; then
        echo "Local db instance terminated. PID=$LOCAL_TURSO_DB_PID"
    else
        echo "Failed to terminate local db instance. PID=$LOCAL_TURSO_DB_PID"
    fi
else
    echo "Local db instance terminated. PID=$LOCAL_TURSO_DB_PID"
fi
