package repository

import (
	"context"
	desc "github.com/beachrockhotel/auth/grpc/pkg/auth_v1"
)

type AuthRepository interface {
	Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error)
	Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error)
	Update(ctx context.Context, req *desc.UpdateRequest) error
	Delete(ctx context.Context, req *desc.DeleteRequest) error
}
