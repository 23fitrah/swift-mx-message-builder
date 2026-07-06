package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config holds runtime settings for the MX builder service.
type Config struct {
	Port        string
	OutputDir   string // root folder where generated MX files are written
	WorkerCount int    // number of goroutines in the file-writer worker pool
	QueueSize   int    // buffered channel size for pending jobs
}

// Load reads configuration from environment variables, falling back to
// sensible defaults for local development.
func Load() Config {
	_ = godotenv.Load() // Load .env file
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
