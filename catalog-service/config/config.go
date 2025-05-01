package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	JWTSecret string

	DBRedisHost     string
	DBRedisPort     string
	DBRedisPassword string

	TokenExpireMinutes string
}

var AppConfig *Config

func InitConfig() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	AppConfig = &Config{
		AppPort: GetEnv("APP_PORT", "8080"),

		DBHost:     GetEnv("DB_HOST", "localhost"),
		DBPort:     GetEnv("DB_PORT", "5432"),
		DBUser:     GetEnv("DB_USER", "postgres"),
		DBPassword: GetEnv("DB_PASSWORD", ""),
		DBName:     GetEnv("DB_NAME", "my_db"),

		JWTSecret: GetEnv("JWT_SECRET", "123"),

		DBRedisHost:     GetEnv("DB_REDIS_HOST", "localhost"),
		DBRedisPort:     GetEnv("DB_REDIS_PORT", "6379"),
		DBRedisPassword: GetEnv("DB_REDIS_PASSWORD", ""),

		TokenExpireMinutes: GetEnv("TOKEN_EXPIRE_MINUTES", "30"),
	}

	log.Println("Loading .env file successful")
}

func GetEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	} else {
		return defaultValue
	}
}

func (config *Config) GetTokenExpireMinutes() (*time.Duration, error) {
	tokenExpireMinutes, err := strconv.Atoi(AppConfig.TokenExpireMinutes)
	if err != nil {
		return nil, fmt.Errorf("value environment variable ACCESS_TOKEN_EXPIRE_MINUTES is not valid")
	}

	expireDuration := time.Duration(tokenExpireMinutes) * time.Minute
	return &expireDuration, nil
}
