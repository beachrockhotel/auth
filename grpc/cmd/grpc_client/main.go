package main

import (
	"context"
	"log"
	"time"

	desc "github.com/beachrockhotel/pkg/auth_v1"
	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	address = "localhost:50051"
	authID  = 12
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	c := desc.NewAuthV1Client(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Get(ctx, &desc.GetRequest{Id: authID})
	if err != nil {
		log.Fatalf("failed to get user by id: %d", authID)
	}

	log.Printf(color.RedString("User info:\n") + color.GreenString("%+v", r.GetInfo()))
}
