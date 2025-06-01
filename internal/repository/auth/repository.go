package auth

import (
	"context"
	"log"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/beachrockhotel/auth/internal/repository/auth/model"
	desc "github.com/beachrockhotel/auth/pkg/auth_v1"
	"github.com/brianvoe/gofakeit"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	tableName = "auth"

	idColumn        = "id"
	nameColumn      = "name"
	emailColumn     = "email"
	roleColumn      = "role"
	passwordColumn  = "password"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

type repo struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *repo {
	return &repo{db: db}
}

func (r *repo) Get(ctx context.Context, req *desc.GetRequest) (model.User, error) {
	builder := (sq.Select(idColumn, nameColumn, emailColumn, roleColumn, passwordColumn, createdAtColumn, updatedAtColumn).
		From(tableName).
		Where(sq.Eq{idColumn: id}).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	var auth model.Auth
	err = r.db.QueryRow(ctx, query, args...).Scan(&auth.ID, &auth.Info.Name, &auth.Info.Email, &auth.Info.Role, &auth.Info.Password, &auth.Info.CreatedAt, &auth.Info.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return converter.ToAuthFromRepo(auth), nil
}

func (r *repo) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	builder := (sq.Insert(tableName).
		Columns(nameColumn, emailColumn, roleColumn, passwordColumn).
		Values(info.Name, info.Email, info.Role, info.Password).
		Suffix("RETURNING id")

	query, args, err := builder.ToSQL()
	if err != nil {
		return nil, err
	}

	var id int64
	err = r.db.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
