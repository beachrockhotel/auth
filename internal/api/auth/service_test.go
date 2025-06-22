package auth_test

import (
	"context"
	"github.com/stretchr/testify/require"
	"log"
	"net"
	"testing"

	"github.com/beachrockhotel/auth/internal/api/auth"
	"github.com/beachrockhotel/auth/internal/model"
	desc "github.com/beachrockhotel/auth/pkg/auth_v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

type mockAuthService struct {
	mock.Mock
}

func (m *mockAuthService) Create(ctx context.Context, info *model.AuthInfo) (int64, error) {
	args := m.Called(ctx, info)
	return int64(args.Int(0)), args.Error(1)
}

func (m *mockAuthService) Get(ctx context.Context, id int64) (*model.Auth, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.Auth), args.Error(1)
}

func dialer(srv desc.AuthV1Server) func(context.Context, string) (net.Conn, error) {
	lis := bufconn.Listen(1024 * 1024)
	s := grpc.NewServer()
	desc.RegisterAuthV1Server(s, srv)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()
	return func(ctx context.Context, s string) (net.Conn, error) {
		return lis.Dial()
	}
}

func TestCreate(t *testing.T) {
	mockSvc := &mockAuthService{}
	impl := auth.NewImplementation(mockSvc)

	// Ожидание мока
	mockSvc.On("Create", mock.Anything, mock.Anything).Return(123, nil)

	conn, err := grpc.DialContext(context.Background(), "", grpc.WithInsecure(), grpc.WithContextDialer(dialer(impl)))
	require.NoError(t, err)
	client := desc.NewAuthV1Client(conn)

	resp, err := client.Create(context.Background(), &desc.CreateRequest{
		Info: &desc.AuthInfo{Name: "Test", Email: "test@example.com"},
	})
	require.NoError(t, err)
	assert.Equal(t, int64(123), resp.GetId())
}
