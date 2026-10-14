package auth

import (
	"os"
	"time"
)

type Config struct {
	DatabaseURL string
	GRPCAddress string
	JWTSecret   string
	JWTTTL      time.Duration
}

func LoadConfig() Config {
	ttl := 24 * time.Hour

	if value := os.Getenv("JWT_TTL"); value != "" {
		if parsed, err := time.ParseDuration(value); err == nil {
			ttl = parsed
		}
	}
	return Config{
		DatabaseURL: os.Getenv("AUTH_DATABASE_URL"),
		GRPCAddress: getEnv(
			"AUTH_GRPC_ADDRESS",
			":50051",
		),
		JWTSecret: os.Getenv("JWT_SECRET"),
		JWTTTL:    ttl,
	}
}

func getEnv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
