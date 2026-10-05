package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	RedisAddr     string
	RedisUsername string
	RedisPassword string
	RabbitMQURL   string
}

func Load() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found")
	}

	return Config{
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        os.Getenv("DB_PORT"),
		DBUser:        os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBName:        os.Getenv("DB_NAME"),
		RedisAddr:     os.Getenv("REDIS_ADDR"),
		RedisUsername: os.Getenv("REDIS_USERNAME"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RabbitMQURL:   os.Getenv("RABBITMQ_URL"),
	}
}