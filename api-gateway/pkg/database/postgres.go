package database

import (
	models2 "api-gateway/internal/models"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func InitDB() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		viper.GetString("database.host"),
		viper.GetString("database.user"),
		viper.GetString("database.password"),
		viper.GetString("database.dbname"),
		viper.GetString("database.port"),
		viper.GetString("database.sslmode"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	// Auto Migrate
	err = db.AutoMigrate(&models2.User{}, &models2.Exercise{}, &models2.Meal{},
		&models2.Sleep{}, &models2.Hydration{}, &models2.Goal{})
	if err != nil {
		log.Fatal("Failed to auto migrate database: ", err)
	}

	return db
}
