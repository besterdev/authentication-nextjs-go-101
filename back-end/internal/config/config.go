package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv              string
	Port                string
	DatabaseURL         string
	DatabaseReadURL     string
	RedisURL            string
	DBMaxOpenConns      int
	DBMaxIdleConns      int
	DBConnMaxLifetime   time.Duration
	BcryptMaxConcurrent int
	JWTAccessSecret     string
	JWTRefreshSecret    string
	JWTAccessExpiry     time.Duration
	JWTRefreshExpiry    time.Duration
	CORSOrigin          string
	OTLPEndpoint        string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	accessExpiry, err := time.ParseDuration(getEnv("JWT_ACCESS_EXPIRY", "15m"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_ACCESS_EXPIRY: %w", err)
	}

	refreshExpiry, err := time.ParseDuration(getEnv("JWT_REFRESH_EXPIRY", "168h"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_REFRESH_EXPIRY: %w", err)
	}

	connMaxLifetime, err := time.ParseDuration(getEnv("DB_CONN_MAX_LIFETIME", "30m"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_CONN_MAX_LIFETIME: %w", err)
	}

	cfg := &Config{
		AppEnv:              getEnv("APP_ENV", "development"),
		Port:                getEnv("PORT", "8080"),
		DatabaseURL:         getEnv("DATABASE_URL", ""),
		DatabaseReadURL:     getEnv("DATABASE_READ_URL", ""),
		RedisURL:            getEnv("REDIS_URL", ""),
		DBMaxOpenConns:      getEnvInt("DB_MAX_OPEN_CONNS", 25),
		DBMaxIdleConns:      getEnvInt("DB_MAX_IDLE_CONNS", 10),
		DBConnMaxLifetime:   connMaxLifetime,
		BcryptMaxConcurrent: getEnvInt("BCRYPT_MAX_CONCURRENT", 4),
		JWTAccessSecret:     getEnv("JWT_ACCESS_SECRET", ""),
		JWTRefreshSecret:    getEnv("JWT_REFRESH_SECRET", ""),
		JWTAccessExpiry:     accessExpiry,
		JWTRefreshExpiry:    refreshExpiry,
		CORSOrigin:          getEnv("CORS_ORIGIN", "http://localhost:3000"),
		OTLPEndpoint:        getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", ""),
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.DatabaseReadURL == "" {
		cfg.DatabaseReadURL = cfg.DatabaseURL
	}
	if len(cfg.JWTAccessSecret) < 32 {
		return nil, fmt.Errorf("JWT_ACCESS_SECRET must be at least 32 characters")
	}
	if len(cfg.JWTRefreshSecret) < 32 {
		return nil, fmt.Errorf("JWT_REFRESH_SECRET must be at least 32 characters")
	}

	return cfg, nil
}

func (c *Config) IsProduction() bool {
	return c.AppEnv == "production"
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		n, err := strconv.Atoi(value)
		if err == nil {
			return n
		}
	}
	return fallback
}
