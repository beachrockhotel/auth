package repository

import (
	"context"
	"github.com/biggiecheesebiggie/auth/grpc/internal/repository/auth/model"
)

type AuthRepository interface {
	Create(ctx context.Context, info *model.AuthInfo) (int64, error)
	Get(ctx context.Context, id int64) (*model.Auth, error)
}
