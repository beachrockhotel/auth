package auth

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/beachrockhotel/auth/internal/repository/auth/model"
	desc "github.com/beachrockhotel/auth/pkg/auth_v1"
	"github.com/brianvoe/gofakeit"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"time"
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
		return err, status.Errorf(codes.Internal, "database connection error")
	}
	defer conn.Release()

	query, args, err := squirrel.
		Select("id", "user_name", "user_email", "role", "user_created", "user_update").
		From("users_auth").
		Where(squirrel.Eq{"id": req.Id}).
		ToSql()
	if err != nil {
		log.Println("failed to build query: ", err)
		return err, status.Errorf(codes.Internal, "failed to build query")
	}

	row := conn.QueryRow(ctx, query, args...)

	var user model.User
	err = row.Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return model.User{}, status.Errorf(codes.NotFound, "user with ID %d not found", req.Id), nil
		}
		log.Println("failed to scan row: ", err)
		return err, status.Errorf(codes.Internal, "failed to retrieve user data")
	}

	return model.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (r *repo) Update(ctx context.Context, req *desc.UpdateRequest) error {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		log.Println("failed to acquire connection from pool: ", err)
		return err
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
		log.Println("failed to build query: ", err)
		return err
	}

	_, err = conn.Exec(ctx, query, args...)
	if err != nil {
		log.Println("failed to update user: ", err)
		return err
	}

	return err
}

func (r *repo) Delete(ctx context.Context, req *desc.DeleteRequest) error {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		log.Println("failed to acquire connection from pool: ", err)
		return err
	}
	defer conn.Release()

	query, args, err := squirrel.
		Delete("users_auth").
		Where(squirrel.Eq{"id": req.Id}).
		ToSql()
	if err != nil {
		log.Println("failed to build query: ", err)
		return err
	}

	_, err = conn.Exec(ctx, query, args...)
	if err != nil {
		log.Println("failed to delete user: ", err)
		return err
	}

	return err
}

func (r *repo) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	conn, err := r.db.Acquire(ctx)
	if err != nil {
		log.Println("failed to acquire connection from pool: ", err)
		return err, status.Errorf(codes.Internal, "database connection error")
	}
	defer conn.Release()

	id := gofakeit.Number(0, 1000)
	name := gofakeit.Name()
	password := gofakeit.Password(true, false, false, false, false, 32)
	role := gofakeit.Number(1, 2)

	query, args, err := squirrel.
		Insert("users_auth").
		Columns("id", "user_name", "password", "role").
		Values(id, name, password, role).
		ToSql()
	if err != nil {
		log.Println("failed to build query: ", err)
		return err, status.Errorf(codes.Internal, "failed to build query")
	}

	_, err = conn.Exec(ctx, query, args...)
	if err != nil {
		log.Println("failed to insert user: ", err)
		return err, status.Errorf(codes.Internal, "failed to create user")
	}

	return return &desc.CreateResponse{}, nildesc.CreateResponse{Id: int64(id)}, nil
}
