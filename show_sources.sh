#!/bin/bash

# Корневая директория (можно заменить на текущую, если скрипт запускается в нужной папке)
ROOT_DIR="./"

# Расширения/имена файлов, которые нас интересуют
EXTENSIONS=("*.go" "go.mod" "go.sum" "*.proto" "*.Dockerfile" "Dockerfile" "docker-compose.yaml" "*.sh" "*.sql" "*.env")

echo "Файлы и их содержимое (только релевантные):"
echo "------------------------------------------"

# Поиск и вывод содержимого
for ext in "${EXTENSIONS[@]}"; do
    find "$ROOT_DIR" -type f -name "$ext" | while read -r file; do
        echo -e "\n--- $file ---"
        cat "$file"
    done
done
