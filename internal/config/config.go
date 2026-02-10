package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	DatabaseURL          string
	JWTSecret            string
	AccessTokenTTL       time.Duration
	Addr                 string
}

func MustLoad() Config {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL переменная окружения не установлена")
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET переменная окружения не установлена")
	}

	ttlStr := os.Getenv("ACCESS_TOKEN_TTL_MINUTES")
	if ttlStr == "" {
		ttlStr = "15"
	}
	ttlMin, err := strconv.Atoi(ttlStr)
	if err != nil || ttlMin <= 0 {
		log.Fatal("ACCESS_TOKEN_TTL_MINUTES должна быть положительным целым числом")
	}

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}

	return Config{
		DatabaseURL:    dbURL,
		JWTSecret:      secret,
		AccessTokenTTL: time.Duration(ttlMin) * time.Minute,
		Addr:           addr,
	}
}
