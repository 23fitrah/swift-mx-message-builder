package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	OutputDir   string
	WorkerCount int
	QueueSize   int
}

func Load() Config {
	_ = godotenv.Load()
	return Config{
		Port:        os.Getenv("APP_PORT"),
		OutputDir:   os.Getenv("MX_OUTPUT_DIR"),
		WorkerCount: 5,
		QueueSize:   100,
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
