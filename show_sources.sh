#!/bin/bash

ROOT_DIR="./"

echo "Файлы и их содержимое (только релевантные):"
echo "------------------------------------------"

find "$ROOT_DIR" -type f \( \
    -name "*.go" -o \
    -name "go.mod" -o \
    -name "go.sum" -o \
    -name "*.proto" -o \
    -name "*.Dockerfile" -o \
    -name "Dockerfile" -o \
    -name "docker-compose.yaml" -o \
    -name "*.sh" -o \
    -name "*.sql" -o \
    -name "*.env" \
\) | while read -r file; do
    echo -e "\n--- $file ---"
    cat "$file"
done
