package auth_test

import (
	"context"
	"testing"

	"github.com/beachrockhotel/auth/internal/client/db" // Вот тут ты импортируешь db!
	"github.com/beachrockhotel/auth/internal/model"
	"github.com/beachrockhotel/auth/internal/service/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Твой мок репозитория
type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) Create(ctx context.Context, info *model.AuthInfo) (int64, error) {
	args := m.Called(ctx, info)
	return int64(args.Int(0)), args.Error(1)
}

func (m *mockRepo) Get(ctx context.Context, id int64) (*model.Auth, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.Auth), args.Error(1)
}

// ✅ Исправленный мок TxManager — именно здесь нужен твой метод ReadCommitted
type mockTx struct {
	mock.Mock
}

func (m *mockTx) ReadCommitted(ctx context.Context, f db.Handler) error {
	return f(ctx)
}

func TestCreate(t *testing.T) {
	repo := &mockRepo{}
	tx := &mockTx{}

	s := auth.NewService(repo, tx)

	repo.On("Create", mock.Anything, mock.Anything).Return(1, nil)
	repo.On("Get", mock.Anything, int64(1)).Return(&model.Auth{ID: 1}, nil)

	id, err := s.Create(context.Background(), &model.AuthInfo{Name: "John"})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)
}
