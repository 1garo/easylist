package env

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Env struct {
	Port        int
	DatabaseURL string
	DatabaseToken string
	DatabaseName string
}

func Load() Env {
	// Load .env file (ignore if missing, system env takes precedence)
	_ = godotenv.Load()

	portStr := getEnv("PORT", "3000")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid PORT: %v", err)
	}

	cfg := Env{
		Port:        port,
		DatabaseURL: getEnv("DATABASE_URL", ""),
		DatabaseToken: getEnv("AUTH_TOKEN", ""),
		DatabaseName: getEnv("DATABASE_NAME", "easylist"),
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
