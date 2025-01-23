#!/usr/bin/bash
source local.env

# Обновляем переменную окружения для корректного пути
export MIGRATION_DSN="host=$PG_HOST port=$PG_PORT dbname=$PG_DATABASE_NAME user=$PG_USER password=$PG_PASSWORD sslmode=disable"

# Проверка содержимого директории миграций для отладки
echo "Checking migrations directory: ${MIGRATIONS_DIR}"
ls -la "${MIGRATIONS_DIR}"

# Запуск миграций
sleep 2 && goose -dir "${MIGRATIONS_DIR}" postgres "${MIGRATION_DSN}" up -v