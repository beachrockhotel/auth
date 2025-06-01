package converter

import (
	"github.com/beachrockhotel/auth/internal/model"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToNoteFromRepo(auth *model.Auth) *desc.Auth {
	var updatedAt *timestamppb.Timestamp
	if auth.UpdatedAt.Valid {
		updatedAt = timestamppb.New(auth.UpdatedAt.Time)
	}

	return &desc.Auth{
		Id:        auth.ID,
		Info:      ToAuthInfoFromRepo(auth.Info),
		CreatedAt: timestamppb.New(auth.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

func ToAuthInfoFromRepo(info *model.AuthInfo) *model.AuthInfo {
	return &model.AuthInfo{
		Name:     auth.Name,
		Email:    auth.Email,
		Role:     auth.Role,
		Password: auth.Password,
	}
}
