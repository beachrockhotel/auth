package auth

import (
	"context"
	"github.com/pkg/errors"
	"log"

	"github.com/beachrockhotel/auth/internal/converter"
	desc "github.com/beachrockhotel/auth/pkg/auth_v1"
)

func (i *Implementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	if req.GetId() == 0 {
		return nil, errors.Errorf("id is empty")
	}

	authObj, err := i.authService.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	log.Printf("id: %d, name: %s, email: %s, role: %s, created_at: %v, updated_at: %v\n", authObj.ID, authObj.Info.Name, authObj.Info.Email, authObj.Info.Role, authObj.CreatedAt, authObj.UpdatedAt)

	return &desc.GetResponse{
		Auth: converter.ToAuthFromService(authObj),
	}, nil
}
