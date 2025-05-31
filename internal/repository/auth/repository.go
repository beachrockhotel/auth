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

type repo struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *repo {
	return &repo{db: db}
}

func (r *repo) Get(ctx context.Context, req *desc.GetRequest) (model.User, error) {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		log.Println("failed to acquire connection from pool: ", err)
		return model.User{}, status.Errorf(codes.Internal, "database connection error")
	}
	defer conn.Release()

	query, args, err := squirrel.
		Select("id", "user_name", "user_email", "role", "user_created", "user_update").
		From("users_auth").
		Where(squirrel.Eq{"id": req.Id}).
		ToSql()
	if err != nil {
		log.Println("failed to build query: ", err)
		return model.User{}, status.Errorf(codes.Internal, "failed to build query")
	}

	row := conn.QueryRow(ctx, query, args...)

	var user model.User
	err = row.Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return model.User{}, status.Errorf(codes.NotFound, "user with ID %d not found", req.Id)
		}
		log.Println("failed to scan row: ", err)
		return model.User{}, status.Errorf(codes.Internal, "failed to retrieve user data")
	}

	return user, nil
}

func (r *repo) Update(ctx context.Context, req *desc.UpdateRequest) error {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		log.Println("failed to acquire connection from pool: ", err)
		return status.Errorf(codes.Internal, "database connection error")
	}
	defer conn.Release()

	query, args, err := squirrel.
		Update("users_auth").
		Set("user_name", req.Name.GetValue()).
		Set("user_email", req.Email.GetValue()).
		Set("user_update", time.Now()).
		Where(squirrel.Eq{"id": req.Id}).
		ToSql()
	if err != nil {
		log.Println("failed to build update query: ", err)
		return status.Errorf(codes.Internal, "failed to build update query")
	}

	_, err = conn.Exec(ctx, query, args...)
	if err != nil {
		log.Println("failed to update user: ", err)
		return status.Errorf(codes.Internal, "failed to update user")
	}

	return nil
}

func (r *repo) Delete(ctx context.Context, req *desc.DeleteRequest) error {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		log.Println("failed to acquire connection from pool: ", err)
		return status.Errorf(codes.Internal, "database connection error")
	}
	defer conn.Release()

	query, args, err := squirrel.
		Delete("users_auth").
		Where(squirrel.Eq{"id": req.Id}).
		ToSql()
	if err != nil {
		log.Println("failed to build delete query: ", err)
		return status.Errorf(codes.Internal, "failed to build delete query")
	}

	_, err = conn.Exec(ctx, query, args...)
	if err != nil {
		log.Println("failed to delete user: ", err)
		return status.Errorf(codes.Internal, "failed to delete user")
	}

	return nil
}

func (r *repo) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		log.Println("failed to acquire connection from pool: ", err)
		return nil, status.Errorf(codes.Internal, "database connection error")
	}
	defer conn.Release()

	id := gofakeit.Number(1000, 9999)
	name := gofakeit.Name()
	password := gofakeit.Password(true, false, false, false, false, 32)
	role := gofakeit.Number(1, 2)

	query, args, err := squirrel.
		Insert("users_auth").
		Columns("id", "user_name", "password", "role").
		Values(id, name, password, role).
		ToSql()
	if err != nil {
		log.Println("failed to build insert query: ", err)
		return nil, status.Errorf(codes.Internal, "failed to build insert query")
	}

	_, err = conn.Exec(ctx, query, args...)
	if err != nil {
		log.Println("failed to insert user: ", err)
		return nil, status.Errorf(codes.Internal, "failed to insert user")
	}

	return &desc.CreateResponse{Id: int64(id)}, nil
}
