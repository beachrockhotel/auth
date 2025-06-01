package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/beachrockhotel/auth/internal/model"
	desc "github.com/beachrockhotel/auth/pkg/auth_v1"
)

func ToAuthFromService(auth *model.Auth) *desc.Auth {
	var updatedAt *timestamppb.Timestamp
	if auth.UpdatedAt.Valid {
		updatedAt = timestamppb.New(auth.UpdatedAt.Time)
	}

	return &desc.Auth{
		Id:        auth.ID,
		Info:      ToAuthInfoFromService(auth.Info),
		CreatedAt: timestamppb.New(auth.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

func ToAuthInfoFromService(info model.AuthInfo) *desc.AuthInfo {
	return &desc.AuthInfo{
		Name:     info.Name,
		Email:    info.Email,
		Role:     info.Role,
		Password: info.Password,
	}
}

func ToAuthInfoFromDesc(info *desc.AuthInfo) *model.AuthInfo {
	return &model.AuthInfo{
		Name:     info.Name,
		Email:    info.Email,
		Role:     info.Role,
		Password: info.Password,
	}
}
