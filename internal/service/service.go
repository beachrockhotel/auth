package service

import (
	"context"
	"github.com/beachrockhotel/auth/internal/model"
)

type AuthService interface {
	Create(ctx context.Context, info *model.AuthInfo) (int64, error)
	Get(ctx context.Context, id int64) (*model.Auth, error)
}
