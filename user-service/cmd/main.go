package main

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"user-service/internal/config"
	userpb "user-service/internal/genproto/user"
	"user-service/internal/repository"
	"user-service/internal/server"
	"user-service/internal/service"
	"user-service/pkg/database"
)

func main() {
	// Initialize the gRPC server
	grpcServer := grpc.NewServer()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	db := database.InitDB(&cfg.Database)

	repo := repository.NewUserRepository(db)

	userService := service.NewUserService(repo)

	userServer := server.NewUserServer(userService)

	userpb.RegisterUserServiceServer(grpcServer, userServer)

	// Start listening for requests
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println(fmt.Sprintf("Starting GRPC server on :%d", cfg.Server.Port))
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
