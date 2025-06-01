package auth

import (
	"context"
	"log"

	"github.com/beachrockhotel/auth/internal/converter"
	desc "github.com/beachrockhotel/auth/pkg/auth_v1"
)

func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	id, err := i.authService.Create(ctx, converter.ToAuthInfoFromDesc(req.GetInfo()))
	if err != nil {
		return nil, err
	}

	log.Printf("inserted user with id: %d", id)

	return &desc.CreateResponse{
		Id: id,
	}, nil
}
