package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	APP_PORT                 string
	DB_URL                   string
	REDIS_HOST               string
	REDIS_PORT               string
	PAYMENT_SERVICE_GRPC_URL string
	KAFKA_BROKERS            string
	PROVIDER_API_URL         string
	PROVIDER_API_KEY         string
	JWT_SECRET               string
}

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env file not found, using system env")
	}
	fmt.Println("✅ .env loaded successfully")
}
func LoadConfig() *Config {
	LoadEnv()

	return &Config{
		APP_PORT:                 os.Getenv("APP_PORT"),
		DB_URL:                   os.Getenv("DB_URL"),
		REDIS_HOST:               os.Getenv("REDIS_HOST"),
		REDIS_PORT:               os.Getenv("REDIS_PORT"),
		PAYMENT_SERVICE_GRPC_URL: os.Getenv("PAYMENT_SERVICE_GRPC_URL"),
		KAFKA_BROKERS:            os.Getenv("KAFKA_BROKERS"),
		PROVIDER_API_URL:         os.Getenv("PROVIDER_API_URL"),
		PROVIDER_API_KEY:         os.Getenv("PROVIDER_API_KEY"),
		JWT_SECRET:               os.Getenv("JWT_SECRET"),
	}
}
