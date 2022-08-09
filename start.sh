#!/bin/sh
set -e

echo "Start the db"
wait-for "${DB_DRIVER}:5432" --

echo "Start the app"
exec "$@"