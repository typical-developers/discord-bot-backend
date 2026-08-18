#!/bin/bash
set -e

missing_env=()
required_env=(
    DATABASE_HOST
    DATABASE_PORT
    DATABASE_USERNAME
    DATABASE_PASSWORD
    DATABASE_OPTIONS
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
until pg_isready -h "$DATABASE_HOST" -p "$DATABASE_PORT" -U "$DATABASE_USERNAME"; do
    sleep 2
done

echo "Running migrations . . ."
migrate -path /migrations -database "postgres://$DATABASE_USERNAME:$DATABASE_PASSWORD@$DATABASE_HOST:$DATABASE_PORT?$DATABASE_OPTIONS" up

echo "Migrations complete."
