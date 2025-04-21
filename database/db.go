package database

import (
	"fmt"
	"log"
	"os"
	"tracker/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB(){

	// make connection using env file

	connectStr:= fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	// Open database connection
	var err error
	DB, err := gorm.Open(postgres.Open(connectStr), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to db: %v", err)
	}

	// migrate models to tables in db

	err = DB.AutoMigrate(&models.Budget{}, &models.Transaction{}, &models.User{})
	if err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}
	fmt.Println("Connection to databse was a success")
}