package main

import (
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"onlineQueue/internal/operators"
	"onlineQueue/internal/registers"
	"os"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Successfully connected to database")
	err = db.AutoMigrate(&operators.Operator{}, &registers.Register{})
	if err != nil {
		log.Fatal("Database migration failed:", err)
	}
	log.Println("Database migration completed successfully")
}
