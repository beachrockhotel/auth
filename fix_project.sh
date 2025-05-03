#!/bin/bash

set -e

echo "🔧 Исправляем импорты model -> internal/repository/auth/model..."
find . -type f -name "*.go" -exec sed -i 's|"github.com/beachrockhotel/auth/internal/model"|"github.com/beachrockhotel/auth/internal/repository/auth/model"|g' {} +

echo "🔧 Исправляем return nil -> return err в Update и Delete..."
sed -i 's|return nil|return err|g' internal/repository/auth/repository.go

echo "🔧 Исправляем return status.Errorf(...) в Get на правильный return..."
sed -i 's|return status.Errorf(\([^)]*\))|return model.User{}, status.Errorf(\1)|g' internal/repository/auth/repository.go

echo "🔧 Исправляем Create: возвращаем ID..."
sed -i 's|return &desc.CreateResponse{}, nil|return &desc.CreateResponse{Id: int64(id)}, nil|' internal/repository/auth/repository.go

echo "🔧 Обновляем model.User: UpdatedAt -> sql.NullTime..."
sed -i 's|UpdatedAt \*time.Time|UpdatedAt sql.NullTime|' internal/repository/auth/model/auth.go
sed -i '1 a\\nimport "database/sql"' internal/repository/auth/model/auth.go

echo "🔧 Исправляем проверку .UpdatedAt.Valid..."
sed -i 's|if auth.UpdatedAt != nil {|if auth.UpdatedAt.Valid {|' internal/repository/auth/converter/auth.go
sed -i 's|updatedAt = timestamppb.New(auth.UpdatedAt)|updatedAt = timestamppb.New(auth.UpdatedAt.Time)|' internal/repository/auth/converter/auth.go

echo "🐳 Обновляем Dockerfile: COPY .. -> COPY ."
sed -i 's|COPY \.\.|COPY .|' Dockerfile

echo "✅ Все исправления внесены. Теперь можешь запустить:"
echo "   go mod tidy && go build ./cmd/grpc_server"
