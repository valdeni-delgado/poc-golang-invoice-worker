package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Profile                  string
	DatabaseURL              string
	StoragePath              string
	S3Bucket                 string
	InvoiceProcessingTimeout time.Duration
	MaxFailures              int
}

func Load() (*Config, error) {
	profile := getEnv("PROFILE", "local")

	var dbURL string
	if profile == "local" {
		dbURL = getEnv("PG_JDBC_URL", "postgres://postgres:postgres@localhost:5432/faturador_local")
	} else {
		dbURL = requireEnv("PG_JDBC_URL")
	}

	maxFailuresStr := getEnv("MAX_FAILURES", "5")
	maxFailures, err := strconv.Atoi(maxFailuresStr)
	if err != nil {
		return nil, fmt.Errorf("MAX_FAILURES inválido: %w", err)
	}

	timeoutSecondsStr := getEnv("MAX_EXECUTION_TIME_SECONDS", "600")
	timeoutSeconds, err := strconv.Atoi(timeoutSecondsStr)
	if err != nil {
		return nil, fmt.Errorf("MAX_EXECUTION_TIME_SECONDS inválido: %w", err)
	}

	return &Config{
		Profile:                  profile,
		DatabaseURL:              dbURL,
		StoragePath:              getEnv("STORAGE_PATH", "/tmp/storage"),
		S3Bucket:                 getEnv("S3_BUCKET", ""),
		InvoiceProcessingTimeout: time.Duration(timeoutSeconds) * time.Second,
		MaxFailures:              maxFailures,
	}, nil
}

func requireEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("Variável de ambiente obrigatória não definida: %s", key))
	}
	return value
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
