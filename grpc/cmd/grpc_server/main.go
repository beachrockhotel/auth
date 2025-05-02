package main

import (
	"context"
	"flag"
	"github.com/beachrockhotel/auth/grpc/internal/config"
	"github.com/beachrockhotel/auth/grpc/internal/config/env"
	desc "github.com/beachrockhotel/auth/grpc/pkg/auth_v1"
	"github.com/brianvoe/gofakeit"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net"
	"time"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config_path", ".env", "path to config file")
}

type Server struct {
	desc.UnimplementedAuthV1Server
	pool *pgxpool.Pool
}

func (s *Server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		log.Println("failed to acquire connection from pool: ", err)
		return nil, status.Errorf(codes.Internal, "database connection error")
	}
	defer conn.Release()

	query := `SELECT id,
       user_name,
       user_email,
       role, 
       user_created,
       user_update 
FROM users_auth
WHERE id = $1`
	row := conn.QueryRow(ctx, query, req.Id)

	var id int64
	var name, email string
	var role int32
	var createdAt, updatedAt time.Time

	err = row.Scan(&id, &name, &email, &role, &createdAt, &updatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user with ID %d not found", req.Id)
		}
		log.Println("failed to scan row: ", err)
		return nil, status.Errorf(codes.Internal, "failed to retrieve user data")
	}

	return &desc.GetResponse{
		Id:        id,
		Name:      name,
		Email:     email,
		Role:      desc.GetResponseRole(role),
		CreatedAt: timestamppb.New(createdAt),
		UpdatedAt: timestamppb.New(updatedAt),
	}, nil
}

// Update modifies user data in the database.
func (s *Server) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		log.Println("failed to acquire connection from pool: ", err)
		return nil, status.Errorf(codes.Internal, "database connection error")
	}
	defer conn.Release()

	query := `
        UPDATE users_auth
        SET 
            user_name = COALESCE($1, user_name),
            user_email = COALESCE($2, user_email),
            user_update = $3
        WHERE id = $4;
    `
	_, err = conn.Exec(ctx, query,
		req.Name.GetValue(),
		req.Email.GetValue(),
		time.Now(), // Always update user_update with the current timestamp
		req.Id,
	)
	if err != nil {
		log.Println("failed to update user: ", err)
		return nil, status.Errorf(codes.Internal, "failed to update user")
	}

	return &emptypb.Empty{}, nil
}

// Delete removes a user from the database.
func (s *Server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		log.Println("failed to acquire connection from pool: ", err)
		return nil, status.Errorf(codes.Internal, "database connection error")
	}
	defer conn.Release()

	query := `DELETE FROM users_auth WHERE id = $1`
	_, err = conn.Exec(ctx, query, req.Id)
	if err != nil {
		log.Println("failed to delete user: ", err)
		return nil, status.Errorf(codes.Internal, "failed to delete user")
	}

	return &emptypb.Empty{}, nil
}

// Create inserts a new user into the database.
func (s *Server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		log.Println("failed to acquire connection from pool: ", err)
		return nil, status.Errorf(codes.Internal, "database connection error")
	}
	defer conn.Release()

	query := "INSERT INTO users_auth (id, user_name, password, role) VALUES ($1, $2, $3, $4)"
	_, err = conn.Exec(ctx, query,
		gofakeit.Number(0, 1000), // Random ID for the user
		gofakeit.Name(),          // Random name
		gofakeit.Password(true, false, false, false, false, 32), // Random password
		gofakeit.Number(1, 2), // Random role
	)
	if err != nil {
		log.Println("failed to insert user: ", err)
		return nil, status.Errorf(codes.Internal, "failed to create user")
	}

	return &desc.CreateResponse{}, nil
}

func main() {
	flag.Parse()
	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	grpcConfig, err := env.NewGRPCConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	pgConfig, err := env.NewPGConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	pool, err := pgxpool.New(context.Background(), pgConfig.DSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterAuthV1Server(s, &Server{})
	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
