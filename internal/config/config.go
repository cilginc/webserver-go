package config

import (
	"os"
	"runtime"
	"strconv"
)

type Config struct {
	Port        string
	StaticDir   string
	WorkerCount int
}

func Load() *Config {
	port := getEnv("PORT", "8080")
	staticDir := getEnv("STATIC_DIR", ".")

	workerCount := runtime.NumCPU()
	if wcStr := os.Getenv("WORKER_COUNT"); wcStr != "" {
		if wc, err := strconv.Atoi(wcStr); err == nil && wc > 0 {
			workerCount = wc
		}
	}

	return &Config{
		Port:        port,
		StaticDir:   staticDir,
		WorkerCount: workerCount,
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
