package main

import (
	"context"
	"flag"
	"github.com/beachrockhotel/auth/internal/config"
	"github.com/beachrockhotel/auth/internal/config/env"
	desc "github.com/beachrockhotel/auth/pkg/auth_v1"
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
	"strings"
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

func toRoleEnum(roleStr string) desc.Role {
	switch strings.ToUpper(roleStr) {
	case "USER":
		return desc.Role_USER
	case "ADMIN":
		return desc.Role_ADMIN
	default:
		return desc.Role_ROLE_UNSPECIFIED
	}
}

func (s *Server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		log.Println("failed to acquire connection from pool: ", err)
		return nil, status.Errorf(codes.Internal, "database connection error")
	}
	defer conn.Release()

	query := `SELECT id, name, email, role, created_at, updated_at FROM auth WHERE id = $1`
	row := conn.QueryRow(ctx, query, req.Id)

	var id int64
	var name, email, roleStr string
	var createdAt, updatedAt time.Time

	err = row.Scan(&id, &name, &email, &roleStr, &createdAt, &updatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user with ID %d not found", req.Id)
		}
		log.Println("failed to scan row: ", err)
		return nil, status.Errorf(codes.Internal, "failed to retrieve user data")
	}

	roleEnum := toRoleEnum(roleStr)

	return &desc.GetResponse{
		Id: id,
		Info: &desc.UserInfo{
			Name:  name,
			Email: email,
			Role:  roleEnum,
		},
		CreatedAt: timestamppb.New(createdAt),
		UpdatedAt: timestamppb.New(updatedAt),
	}, nil
}

func (s *Server) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		log.Println("failed to acquire connection from pool: ", err)
		return nil, status.Errorf(codes.Internal, "database connection error")
	}
	defer conn.Release()

	query := `
        UPDATE auth
        SET 
            name = COALESCE($1, name),
            email = COALESCE($2, email),
            updated_at = $3
        WHERE id = $4;
    `
	_, err = conn.Exec(ctx, query,
		req.Name.GetValue(),
		req.Email.GetValue(),
		time.Now(),
		req.Id,
	)
	if err != nil {
		log.Println("failed to update user: ", err)
		return nil, status.Errorf(codes.Internal, "failed to update user")
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		log.Println("failed to acquire connection from pool: ", err)
		return nil, status.Errorf(codes.Internal, "database connection error")
	}
	defer conn.Release()

	query := `DELETE FROM auth WHERE id = $1`
	_, err = conn.Exec(ctx, query, req.Id)
	if err != nil {
		log.Println("failed to delete user: ", err)
		return nil, status.Errorf(codes.Internal, "failed to delete user")
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		log.Println("failed to acquire connection from pool: ", err)
		return nil, status.Errorf(codes.Internal, "database connection error")
	}
	defer conn.Release()

	id := gofakeit.Number(1000, 9999)
	name := req.GetInfo().GetName()
	email := req.GetInfo().GetEmail()
	role := req.GetInfo().GetRole().String()
	password := req.GetPassword()

	query := `
		INSERT INTO auth (id, name, email, password, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
	`
	_, err = conn.Exec(ctx, query, id, name, email, password, role)
	if err != nil {
		log.Println("failed to insert user: ", err)
		return nil, status.Errorf(codes.Internal, "failed to create user")
	}

	return &desc.CreateResponse{Id: int64(id)}, nil
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
	desc.RegisterAuthV1Server(s, &Server{pool: pool})
	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
