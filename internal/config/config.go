package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseUrl string
	AppPort     string
}

func Load() *Config {
	err := godotenv.Load("../../.env")
	if err != nil {
		fmt.Println("Error loading .env file, using environment variables")
	}

	dbUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	return &Config{
		DatabaseUrl: dbUrl,
		AppPort:     os.Getenv("APP_PORT"),
	}
}
