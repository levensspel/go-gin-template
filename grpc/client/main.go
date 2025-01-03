package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	pb "github.com/levensspel/go-gin-template/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	grpcHost := os.Getenv("GRPC_HOST")
	if grpcHost == "" {
		grpcHost = "localhost"
	}

	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}

	address := fmt.Sprintf("%s:%s", grpcHost, grpcPort)

	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	response, err := c.RegisterUser(ctx, &pb.RequestRegister{
		Id:       "generated_uuid_here",
		Username: "grpc_username",
		Email:    "grpc.client@email.com",
		Password: "plain_password_sample",
	})
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	log.Printf("Returned ID: %s", response.GetUserId())
	log.Printf("Returned Status Code: %d", response.GetStatusCode())
	log.Printf("Returned Message: %s", response.GetMessage())
}
