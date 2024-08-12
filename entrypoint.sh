#!/bin/sh

# 環境変数MYSQL_HOSTからMySQLホスト名を取得（デフォルト値: "mysql"）
MYSQL_HOST=${MYSQL_HOST:-"mysql"}

echo "Waiting for MySQL to start on host $MYSQL_HOST..."

# MySQLが起動するまで待機
while ! mysqladmin ping -h"${MYSQL_HOST}" --silent; do
    sleep 1
done

echo "MySQL is up - executing command"
exec "$@"