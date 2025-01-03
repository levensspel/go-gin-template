package server

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/levensspel/go-gin-template/config"
	"github.com/levensspel/go-gin-template/helper"
	"github.com/levensspel/go-gin-template/middleware"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Start() error {
	db, err := config.NewDbInit()

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	helper.WORK_DIR = wd

	r := gin.Default()
	r.Use(middleware.EnableCORS)

	extGrpcHost := os.Getenv("EXTERNAL_GRPC_HOST")
	extGrpcPort := os.Getenv("EXTERNAL_GRPC_PORT")
	if extGrpcHost == "" || extGrpcPort == "" {
		log.Fatalf("did not connect: %v", helper.ErrInvalidEnvKey)
	}
	extAddr := fmt.Sprintf("%s:%s", extGrpcHost, extGrpcPort)

	grpcClient, err := grpc.NewClient(extAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer grpcClient.Close()

	grpcServer := grpc.NewServer()

	NewRouter(r, db, grpcClient, grpcServer)

	r.Use(gin.Recovery())

	go runGrpcServer(grpcServer)

	httpPort := os.Getenv("PORT")
	if len(httpPort) == 0 {
		httpPort = "8080"
	}

	appEnv := os.Getenv("MODE")

	switch appEnv {
	case "PRODUCTION":
		gin.SetMode(gin.ReleaseMode)

		sslCert := os.Getenv("SSL_CERT_PATH")
		sslKey := os.Getenv("SSL_KEY_PATH")

		if sslCert == "" || sslKey == "" {
			log.Fatal("SSL certificates not configured")
		}

		host := os.Getenv("PROD_HOST")
		err := r.RunTLS(
			fmt.Sprintf("%s:%s", host, httpPort),
			sslCert,
			sslKey,
		)
		if err != nil {
			log.Fatalf("Failed to start HTTPS server: %v", err)
		}
	default:
		gin.SetMode(gin.DebugMode)
		host := os.Getenv("DEBUG_HOST")
		r.Run(fmt.Sprintf("%s:%s", host, httpPort))
	}

	return nil
}

func runGrpcServer(s *grpc.Server) {
	grpcHost := os.Getenv("GRPC_HOST")
	if grpcHost == "" {
		grpcHost = "localhost"
	}

	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}

	address := fmt.Sprintf("%s:%s", grpcHost, grpcPort)

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
