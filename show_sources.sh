#!/bin/bash

# Выходной файл
OUTPUT_FILE="sources_dump.txt"
> "$OUTPUT_FILE"  # Очистить файл перед началом

# Каталоги и шаблоны, которые надо исключить
EXCLUDE_DIRS=("bin" ".git" ".idea")
EXCLUDE_FILES=("go.mod" "go.sum")
EXCLUDE_PATTERNS=("*.pb.go" "*.pb" "*.exe" "*.so" "*.out" "*.a")

# Проверка: исключить ли файл по имени
should_exclude_file() {
    local filename="$1"

    for pattern in "${EXCLUDE_PATTERNS[@]}"; do
        if [[ "$filename" == $pattern ]]; then
            return 0
        fi
    done

    for f in "${EXCLUDE_FILES[@]}"; do
        if [[ "$(basename "$filename")" == "$f" ]]; then
            return 0
        fi
    done

    return 1
}

# Основной проход по всем файлам
find . -type f | while read -r file; do
    # Пропуск исключённых директорий
    skip_file=0
    for dir in "${EXCLUDE_DIRS[@]}"; do
        if [[ "$file" == ./$dir/* ]]; then
            skip_file=1
            break
        fi
    done

    if [ "$skip_file" -eq 1 ]; then
        continue
    fi

    # Пропуск по расширению или имени
    if should_exclude_file "$file"; then
        continue
    fi

    # Убедиться, что это текстовый файл
    if file "$file" | grep -qvE 'text|ASCII'; then
        continue
    fi

    # Запись пути и содержимого в файл
    {
        echo "==== $file ===="
        cat "$file"
        echo
    } >> "$OUTPUT_FILE"
done

echo "Содержимое сохранено в $OUTPUT_FILE"
