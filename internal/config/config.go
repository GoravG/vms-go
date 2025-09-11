package config

import "os"

type Config struct {
	Addr            string
	HMAC_SECRET_KEY []byte
	DB_USER         string
	DB_PASS         string
	DB_HOST         string
	DB_PORT         string
	DB_NAME         string
	REDIS_HOST      string
	REDIS_PORT      string
	REDIS_CHANNEL   string
	LOG_LEVEL       string
	LOG_FILE        string // Path to log file (optional)
	LOG_FORMAT      string // "plain" or "ecs" (default: plain)
}

func FromEnv() Config {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}
	return Config{
		Addr:            addr,
		HMAC_SECRET_KEY: []byte(os.Getenv("HMAC_SECRET_KEY")),
		DB_USER:         os.Getenv("DB_USER"),
		DB_PASS:         os.Getenv("DB_PASS"),
		DB_HOST:         os.Getenv("DB_HOST"),
		DB_PORT:         os.Getenv("DB_PORT"),
		DB_NAME:         os.Getenv("DB_NAME"),
		REDIS_HOST:      os.Getenv("REDIS_HOST"),
		REDIS_PORT:      os.Getenv("REDIS_PORT"),
		REDIS_CHANNEL:   os.Getenv("REDIS_CHANNEL"),
		LOG_LEVEL:       os.Getenv("LOG_LEVEL"),
		LOG_FILE:        os.Getenv("LOG_FILE"),
		LOG_FORMAT:      os.Getenv("LOG_FORMAT"),
	}
}
