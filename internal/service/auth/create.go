package auth

import (
	"context"

	"github.com/beachrockhotel/auth/internal/model"
)

func (s *serv) Create(ctx context.Context, info *model.AuthInfo) (int64, error) {
	var id int64
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		id, errTx = s.authRepository.Create(ctx, info)
		if errTx != nil {
			return errTx
		}

		_, errTx = s.authRepository.Get(ctx, id)
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return id, nil
}
