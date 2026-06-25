package main

import (
	"context"
	"log"
	"time"

	"github.com/Nishishei01/Go_Hexagonal/proto/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Error to connect with gRPC server: %v", err)
	}
	defer conn.Close()

	authClient := auth.NewAuthServiceClient(conn)

	callRegister(authClient)
}

func callRegister(client auth.AuthServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &auth.RegisterRequest{
		Username:  "test123",
		Password:  "test12",
		Email:     "test123@gmail.com",
		FirstName: "test",
		LastName:  "eiei",
	}

	res, err := client.Register(ctx, req)
	if err != nil {
		log.Printf("Register Failed: %v", err)
		return
	}

	log.Printf("Register Successfully!: %s", res.GetMessage())
}
