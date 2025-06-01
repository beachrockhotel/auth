package note

import "github.com/beachrockhotel/auth/internal/model"

func (s *serv) Create(ctx context.Context, info *model.AuthInfo) (int64, error) {
	id, err := s.authRepository.Create(ctx, info)
	if err != nil {
		return 0, err
	}

	return id, nil
}
