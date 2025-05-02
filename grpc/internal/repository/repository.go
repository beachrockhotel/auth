package repository

import (
	"context"
)

type AuthRepository interface {
	Create(ctx context.Context, name, email, password string) (int64, error)
	Get(ctx context.Context, id int64) (string, error)
}
