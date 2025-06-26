package auth

import (
	"context"
	"strings"
	"time"

	"github.com/beachrockhotel/auth/internal/utils"
	"github.com/beachrockhotel/auth/pkg/auth_v1"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	refreshTokenSecretKey = "W4/X+LLjehdxptt4YgGFCvMpq5ewptpZZYRHY6A72g0="
	accessTokenSecretKey  = "VqvguGiffXILza1f44TWXowDT4zwf03dtXmqWW4SYyE="

	refreshTokenExpiration = 60 * time.Minute
	accessTokenExpiration  = 5 * time.Minute
)

// Встраивай эти методы в свою реализацию Implementation

func (i *Implementation) Login(ctx context.Context, req *auth_v1.LoginRequest) (*auth_v1.LoginResponse, error) {
	// здесь должен быть реальный поиск юзера в БД и валидация пароля
	token, err := utils.GenerateToken(
		model.AuthInfo{Name: req.GetUsername(), Role: auth_v1.Role_ADMIN},
		[]byte(refreshTokenSecretKey),
		refreshTokenExpiration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate token")
	}
	return &auth_v1.LoginResponse{RefreshToken: token}, nil
}

func (i *Implementation) GetRefreshToken(ctx context.Context, req *auth_v1.GetRefreshTokenRequest) (*auth_v1.GetRefreshTokenResponse, error) {
	claims, err := utils.VerifyToken(req.GetRefreshToken(), []byte(refreshTokenSecretKey))
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "invalid refresh token")
	}

	token, err := utils.GenerateToken(
		model.AuthInfo{Name: claims.Name, Role: claims.Role},
		[]byte(refreshTokenSecretKey),
		refreshTokenExpiration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate new refresh token")
	}
	return &auth_v1.GetRefreshTokenResponse{RefreshToken: token}, nil
}

func (i *Implementation) GetAccessToken(ctx context.Context, req *auth_v1.GetAccessTokenRequest) (*auth_v1.GetAccessTokenResponse, error) {
	claims, err := utils.VerifyToken(req.GetRefreshToken(), []byte(refreshTokenSecretKey))
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "invalid refresh token")
	}

	accessToken, err := utils.GenerateToken(
		model.AuthInfo{Name: claims.Name, Role: claims.Role},
		[]byte(accessTokenSecretKey),
		accessTokenExpiration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate access token")
	}
	return &auth_v1.GetAccessTokenResponse{AccessToken: accessToken}, nil
}
