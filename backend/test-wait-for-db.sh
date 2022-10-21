#!/bin/sh

set -e
  
until PGPASSWORD=$DATABASE_PASSWORD psql -h "$DATABASE_URL" -d "$DATABASE_NAME" -U "$DATABASE_USERNAME" -c '\q'; do
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 1
done
  
>&2 echo "Postgres is up - executing command"
exec "$@"