package auth

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"github.com/beachrockhotel/internal/client/db"
	"github.com/beachrockhotel/internal/model"
	"github.com/beachrockhotel/internal/repository"
	"github.com/beachrockhotel/internal/repository/auth/converter"
	modelRepo "github.com/beachrockhotel/internal/repository/auth/model"
)

const (
	tableName = "auth"

	idColumn        = "id"
	nameColumn      = "name"
	emailColumn     = "email"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.AuthRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, info *model.AuthInfo) (int64, error) {
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, emailColumn).
		Values(info.Name, info.Email).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "auth_repository.Create",
		QueryRaw: query,
	}

	var id int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repo) Get(ctx context.Context, id int64) (*model.Auth, error) {
	builder := sq.Select(idColumn, nameColumn, emailColumn, createdAtColumn, updatedAtColumn).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{idColumn: id}).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "auth_repository.Get",
		QueryRaw: query,
	}

	var auth modelRepo.Auth
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(
		&auth.ID,
		&auth.Name,
		&auth.Email,
		&auth.CreatedAt,
		&auth.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return converter.ToAuthFromRepo(&auth), nil
}
