#!/bin/sh

MYSQL_HOST=${MYSQL_HOST:-"mysql"}

echo "Waiting for MySQL to start on host $MYSQL_HOST..."

while ! mysqladmin ping -h"${MYSQL_HOST}" --silent; do
    sleep 1
done

echo "MySQL is up - executing command"
exec "$@"