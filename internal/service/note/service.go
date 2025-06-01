package note

import (
	"github.com/beachrockhotel/auth/internal/repository"
	def "github.com/beachrockhotel/auth/internal/service"
)

var _ def.AuthService = (*serv)(nil)

type serv struct {
	authRepository repository.AuthRepository
}

func NewService(authRepository repository.AuthRepository) *serv {
	return &serv{
		authRepository: authRepository,
	}
}
