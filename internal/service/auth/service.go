package auth

import (
	"github.com/beachrockhotel/auth/internal/client/db"
	"github.com/beachrockhotel/auth/internal/repository"
	"github.com/beachrockhotel/auth/internal/service"
)

type serv struct {
	authRepository repository.AuthRepository
	txManager      db.TxManager
}

func NewService(
	authRepository repository.AuthRepository,
	txManager db.TxManager,
) service.AuthService {
	return &serv{
		authRepository: authRepository,
		txManager:      txManager,
	}
}
