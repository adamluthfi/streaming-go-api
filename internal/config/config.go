package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseDSN string
	JwtSecret   []byte
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		Port:        getEnv("APP_PORT", "8080"),
		DatabaseDSN: getEnv("DATABASE_DSN", ""),
		JwtSecret:   []byte(getEnv("JWT_SECRET", "SUPER_SECRET_KEY")),
	}

	if cfg.DatabaseDSN == "" {
		log.Fatal("DATABASE_DSN is required")
	}

	if len(cfg.JwtSecret) == 0 {
		log.Fatal("JWT_SECRET is required")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
