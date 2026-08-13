package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_URL     string
	JWT_SECRET string
	PORT       string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		DB_URL:     os.Getenv("DB_URL"),
		JWT_SECRET: os.Getenv("JWT_SECRET"),
		PORT:       os.Getenv("PORT"),
	}

	if cfg.DB_URL == "" || cfg.JWT_SECRET == "" {
		return nil, fmt.Errorf("Required environment variables are missing")
	}

	return cfg, nil
}
