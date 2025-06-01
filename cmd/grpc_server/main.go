package main

import (
	"context"
	"database/sql"
	"flag"
	"github.com/beachrockhotel/auth/internal/config"
	"github.com/beachrockhotel/auth/internal/config/env"
	"github.com/beachrockhotel/auth/internal/converter"
	"github.com/beachrockhotel/auth/internal/service"
	desc "github.com/beachrockhotel/auth/pkg/auth_v1"
	"github.com/brianvoe/gofakeit"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net"
	"strings"
	"time"
)

type Server struct {
	desc.UnimplementedAuthV1Server
	authService service.AuthService
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
	authObj, err := s.authService.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	log.Printf("id: %d, name: %s, email: %s, role: %s, password: %s, created_at: %v, updated_at: %v\n", authObj.ID, authObj.Info)

	return &desc.GetResponse{
		Auth: converter.ToAuthInfoFromService(authObj),
	}, nil
}

func (s *Server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	id, err := s.authService.Create(ctx, converter.ToAuthInfoFromDesc(req.GetInfo()))
	if err != nil {
		return nil, err
	}

	log.Printf("inserted auth with id: %d", id)

	return &desc.CreateResponse{
		Id: id,
	}, nil
}

func main() {
	ctx := context.Background()

	err := config.Load(".env")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcConfig, err := config.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	pgConfig, err := config.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	lis, err := net.Listen("tcp", grpcConfig.GRPCAddress())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	pool, err := pgxpool.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	authRepo := authRepository.NewRepository(pool)
	authSrv := authService.NewService(authRepo)

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterAuthV1Server(s, &Server{pool: pgConfig})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
