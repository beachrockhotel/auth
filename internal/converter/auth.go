package converter

import (
	"github.com/beachrockhotel/auth/internal/model"
	desc "github.com/beachrockhotel/auth/pkg/auth_v1"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToAuthFromService(auth *model.Auth) *desc.GetResponse {
	var updatedAt *timestamppb.Timestamp
	if auth.UpdatedAt.Valid {
		updatedAt = timestamppb.New(auth.UpdatedAt.Time)
	}

	return &desc.GetResponse{
		Id: auth.ID,
		Info: &desc.UserInfo{
			Name:  auth.Info.Name,
			Email: auth.Info.Email,
			Role:  desc.Role(desc.Role_value[auth.Info.Role]),
		},
		CreatedAt: timestamppb.New(auth.CreatedAt),
		UpdatedAt: updatedAt,
		RoleEnum:  auth.Info.Role,
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

func ToAuthInfoFromDesc(info desc.AuthInfo) *model.AuthInfo {
	return &desc.AuthInfo{
		Name:     info.Name,
		Email:    info.Email,
		Role:     info.Role,
		Password: info.Password,
	}
}
