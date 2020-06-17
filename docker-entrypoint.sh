#!/bin/sh
set -e

if [ "$1" = 'project' ]; then
    ./pgmigrate -database postgres://admin:moneyforward@123@db:5432/postgres?sslmode=disable -path db/migrations up
    ./project
fi

#exec "$@"