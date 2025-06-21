package auth_test

import (
	"context"
	"os"
	"testing"

	"github.com/beachrockhotel/auth/internal/client/db/pg"
	"github.com/beachrockhotel/auth/internal/model"
	repoAuth "github.com/beachrockhotel/auth/internal/repository/auth"
	"github.com/stretchr/testify/require"
)

func TestCreateAndGet(t *testing.T) {

	dsn := os.Getenv("PG_DSN")
	ctx := context.Background()

	client, err := pg.New(ctx, dsn)
	require.NoError(t, err)
	defer client.Close()

	repo := repoAuth.NewRepository(client)

	id, err := repo.Create(ctx, &model.AuthInfo{
		Name:     "Test",
		Email:    "test@example.com",
		Role:     1,
		Password: "secret",
	})
	require.NoError(t, err)

	authObj, err := repo.Get(ctx, id)
	require.NoError(t, err)
	require.Equal(t, "Test", authObj.Info.Name)
}
