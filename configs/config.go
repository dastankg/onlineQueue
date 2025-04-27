package configs

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	Db   DbConfig
	Auth AuthConfig
}

type DbConfig struct {
	DatabaseUrl string
}

type AuthConfig struct {
	AccessSecret  string
	RefreshSecret string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Warning: No .env file found. Continuing with system environment variables")
	}

	config := &Config{
		Db: DbConfig{
			DatabaseUrl: os.Getenv("DATABASE_URL"),
		},
		Auth: AuthConfig{
			AccessSecret:  os.Getenv("ACCESS_SECRET"),
			RefreshSecret: os.Getenv("REFRESH_SECRET"),
		},
	}
	if config.Db.DatabaseUrl == "" {
		return nil, ErrMissingEnvVar("DATABASE_URL")
	}
	if config.Auth.AccessSecret == "" {
		return nil, ErrMissingEnvVar("ACCESS_SECRET")
	}
	if config.Auth.RefreshSecret == "" {
		return nil, ErrMissingEnvVar("REFRESH_SECRET")
	}
	return config, nil
}
