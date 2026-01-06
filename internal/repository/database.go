package repository

import (
	"fmt"
	"log"
	"todo-api/internal/config"
	"todo-api/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB(cfg *config.Config) error {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	// automatic migrate
	err = DB.AutoMigrate(&models.User{}, &models.Todo{})
	if err != nil {
		return err
	}

	log.Println("✅ Connected to database successfully")
	return nil
}

func GetDB() *gorm.DB {
	return DB
}
