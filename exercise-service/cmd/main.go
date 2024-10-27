package main

import (
	"exercise-service/internal/config"
	exercisepb "exercise-service/internal/genproto/exercises"
	"exercise-service/internal/repository"
	"exercise-service/internal/server"
	"exercise-service/internal/service"
	"exercise-service/pkg/database"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	// Initialize the gRPC server
	grpcServer := grpc.NewServer()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	db := database.InitDB()

	repo := repository.NewExerciseRepository(db)

	exerciseService := service.NewExerciseService(repo)

	exerciseServer := server.NewExerciseServer(exerciseService)

	exercisepb.RegisterExerciseServiceServer(grpcServer, exerciseServer)

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
