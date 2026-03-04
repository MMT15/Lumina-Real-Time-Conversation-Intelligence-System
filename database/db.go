package database

import (
	"fmt"
	"log"
	"os"

	"github.com/MMT15/Lumina/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	host := getEnv("DB_HOST", "localhost")
	user := getEnv("DB_USER", "user")
	password := getEnv("DB_PASSWORD", "password")
	dbname := getEnv("DB_NAME", "lumina")
	port := getEnv("DB_PORT", "5432")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		host, user, password, dbname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto Migrate the models
	err = db.AutoMigrate(&models.Conversation{}, &models.Ticket{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	DB = db
	fmt.Println("Database connection established and migrated successfully")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
