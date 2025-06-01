package auth

import (
	"context"
	"log"

	"github.com/beachrockhotel/auth/internal/converter"
	desc "github.com/beachrockhotel/auth/pkg/auth_v1"
)

func (i *Implementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	authObj, err := i.authService.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	log.Printf("id: %d, title: %s, body: %s, created_at: %v, updated_at: %v\n", authObj.ID, authObj.Info.Title, authObj.Info.Content, authObj.CreatedAt, authObj.UpdatedAt)

	return &desc.GetResponse{
		Auth: converter.ToAuthFromService(authObj),
	}, nil
}
