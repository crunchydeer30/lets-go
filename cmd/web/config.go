package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
	DB   DbConfig
}

type DbConfig struct {
	driver string
	url    string
}

func LoadConfig() (cfg *Config) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg = &Config{
		Port: os.Getenv("PORT"),
		DB: DbConfig{
			driver: os.Getenv("DB_DRIVER"),
			url:    os.Getenv("DB_URL"),
		},
	}
	return
}
