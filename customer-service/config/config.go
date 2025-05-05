package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string

	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string

	JWTSecret          string
	TokenExpireMinutes string

	RedisHost     string
	RedisPort     string
	RedisPassword string
}

var AppConfig *Config

func InitConfig() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	AppConfig = &Config{
		AppPort: GetEnv("APP_PORT", "8080"),

		PostgresHost:     GetEnv("POSTGRES_HOST", "localhost"),
		PostgresPort:     GetEnv("POSTGRES_PORT", "5432"),
		PostgresUser:     GetEnv("POSTGRES_USER", "postgres"),
		PostgresPassword: GetEnv("POSTGRES_PASSWORD", ""),
		PostgresDB:       GetEnv("POSTGRES_DB", "my_db"),

		JWTSecret:          GetEnv("JWT_SECRET", "123"),
		TokenExpireMinutes: GetEnv("TOKEN_EXPIRE_MINUTES", "30"),

		RedisHost:     GetEnv("REDIS_HOST", "localhost"),
		RedisPort:     GetEnv("REDIS_PORT", "6379"),
		RedisPassword: GetEnv("REDIS_PASSWORD", ""),
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

func (config *Config) GetTokenExpireMinutes() *time.Duration {
	tokenExpireMinutes, err := strconv.Atoi(AppConfig.TokenExpireMinutes)
	if err != nil {
		log.Fatal("Value of environment variable ACCESS_TOKEN_EXPIRE_MINUTES is not valid")
		return nil
	}

	expireDuration := time.Duration(tokenExpireMinutes) * time.Minute
	return &expireDuration
}
