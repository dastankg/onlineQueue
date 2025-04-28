package main

import (
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"onlineQueue/internal/offices"
	"onlineQueue/internal/operators"
	"os"
)

func InitializeDefaultAdmin(db *gorm.DB) {
	var adminCount int64
	db.Model(&operators.Operator{}).Where("is_admin = ?", true).Count(&adminCount)

	if adminCount == 0 {
		adminLogin := os.Getenv("ADMIN_LOGIN")
		adminPassword := os.Getenv("ADMIN_PASSWORD")

		if adminLogin == "" || adminPassword == "" {
			log.Fatal("ADMIN_LOGIN и ADMIN_PASSWORD должны быть указаны в переменных окружения")
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("Ошибка при хешировании пароля: %v", err)
		}

		admin := operators.Operator{
			Name:     "Administrator",
			Login:    adminLogin,
			Password: string(hashedPassword),
			IsActive: true,
			IsAdmin:  true,
		}

		result := db.Create(&admin)
		if result.Error != nil {
			log.Fatalf("Не удалось создать администратора: %v", result.Error)
		} else {
			log.Println("Администратор по умолчанию успешно создан")
		}
	}
}
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
	err = db.AutoMigrate(&operators.Operator{}, &offices.Office{})
	if err != nil {
		log.Fatal("Database migration failed:", err)
	}
	log.Println("Database migration completed successfully")
	InitializeDefaultAdmin(db)

}
