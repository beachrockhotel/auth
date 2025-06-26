package access

import (
	"context"
	"google.golang.org/grpc/metadata"
	"strings"

	"github.com/beachrockhotel/auth/internal/model"
	"github.com/beachrockhotel/auth/internal/utils"
	"github.com/beachrockhotel/auth/pkg/access_v1"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	authPrefix        = "Bearer "
	accessTokenSecret = "VqvguGiffXILza1f44TWXowDT4zwf03dtXmqWW4SYyE="
)

type Implementation struct {
	access_v1.UnimplementedAccessV1Server
}

func NewImplementation() *Implementation {
	return &Implementation{}
}

var accessibleRoles = map[string]string{
	model.ExamplePath: "admin",
}

func (i *Implementation) Check(ctx context.Context, req *access_v1.CheckRequest) (*emptypb.Empty, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("metadata is not provided")
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return nil, errors.New("authorization header is not provided")
	}

	if !strings.HasPrefix(authHeader[0], authPrefix) {
		return nil, errors.New("invalid authorization header format")
	}

	token := strings.TrimPrefix(authHeader[0], authPrefix)
	claims, err := utils.VerifyToken(token, []byte(accessTokenSecret))
	if err != nil {
		return nil, errors.New("access token is invalid")
	}

	role, ok := accessibleRoles[req.GetEndpointAddress()]
	if !ok || role != claims.Role.String() {
		return nil, errors.New("access denied")
	}

	return &emptypb.Empty{}, nil
}
