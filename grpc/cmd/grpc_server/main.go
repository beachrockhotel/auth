package main

import (
	"context"
	"flag"
	"log"
	"net"

	"github.com/beachrockhotel/internal/client/db"
	"github.com/beachrockhotel/internal/config"
	"github.com/beachrockhotel/internal/config/env"
	"github.com/beachrockhotel/internal/model"
	authRepoPkg "github.com/beachrockhotel/internal/repository/auth"
	"github.com/beachrockhotel/internal/repository/converter"
	desc "github.com/beachrockhotel/pkg/auth_v1"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var configPath string

type server struct {
	desc.UnimplementedAuthV1Server
	authRepository authRepoPkg.AuthRepository
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	info := req.GetInfo()
	id, err := s.authRepository.Create(ctx, info.GetName(), info.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, err
	}
	return &desc.CreateResponse{Id: id}, nil
}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	auth, err := s.authRepository.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	// Здесь предполагается, что репозиторий возвращает структуру *model.Auth
	authModel, ok := auth.(*model.Auth)
	if !ok {
		return nil, err
	}

	return &desc.GetResponse{
		Id: authModel.ID,
		Info: &desc.UserInfo{
			Name:  authModel.Info.Title,
			Email: authModel.Info.Content,
			Role:  desc.Role_USER, // TODO: маппить из строки "role" при необходимости
		},
		CreatedAt: timestamppb.New(authModel.CreatedAt),
		UpdatedAt: timestamppb.New(authModel.UpdatedAt.Time),
	}, nil
}

func (s *server) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	log.Printf("Update not implemented, id: %d, name: %v, email: %v", req.GetId(), req.GetName(), req.GetEmail())
	return &emptypb.Empty{}, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	log.Printf("Delete not implemented, id: %d", req.GetId())
	return &emptypb.Empty{}, nil
}

func main() {
	flag.Parse()
	ctx := context.Background()

	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcConfig, err := env.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	pgConfig, err := env.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	pool, err := pgxpool.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	dbClient := db.NewClient(pool)
	authRepo := authRepoPkg.NewRepository(dbClient)

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterAuthV1Server(s, &server{authRepository: authRepo})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
