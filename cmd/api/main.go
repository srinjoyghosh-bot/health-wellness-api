package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"healthApi/internal/config"
	"healthApi/pkg/database"
	"log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	db := database.InitDB()
	db.Begin()
	router := gin.Default()
	port := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Print(port)
	if err := router.Run(port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
