package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	return nil
}

func DbUrl() string {
	u := os.Getenv("DB_URL")
	if u == "" {
		log.Println("URL DB is not found")
	}
	return u
}
