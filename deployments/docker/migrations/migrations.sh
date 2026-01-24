#!/bin/bash
set -e

missing_env=()
required_env=(
    POSTGRES_HOST
    POSTGRES_PORT
    POSTGRES_USER
    POSTGRES_PASSWORD
    POSTGRES_DB
)

for var in "${required_env[@]}"; do
    if [ -z "${!var}" ]; then
        missing_env+=("$var")
    fi
done

if [ ${#missing_env[@]} -gt 0 ]; then
    echo "Missing required environment variables: ${missing_env[@]}"
    exit 1
fi

echo "Waiting for primary database to be available before running migrations . . ."
until pg_isready -h "$POSTGRES_HOST" -p "$POSTGRES_PORT" -U "$POSTGRES_USER"; do
    sleep 2
done

echo "Running migrations . . ."
migrate -path /migrations -database "postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_HOST:$POSTGRES_PORT/$POSTGRES_DB?sslmode=disable" up

echo "Migrations complete."