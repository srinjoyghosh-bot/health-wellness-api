package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"user-service/internal/config"
	"user-service/internal/models"
)

func InitDB(cfg *config.DBConfig) *gorm.DB {

	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	// Auto Migrate
	err = db.AutoMigrate(models.User{})
	if err != nil {
		log.Fatal("Failed to auto migrate database: ", err)
	}

	return db
}
