package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	DBRedisHost     string
	DBRedisPort     string
	DBRedisPassword string

	ESHost     string
	ESPort     string
	ESUsername string
	ESPassword string
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

		DBRedisHost:     GetEnv("DB_REDIS_HOST", "localhost"),
		DBRedisPort:     GetEnv("DB_REDIS_PORT", "6379"),
		DBRedisPassword: GetEnv("DB_REDIS_PASSWORD", ""),

		ESHost:     GetEnv("ES_HOST", "localhost"),
		ESPort:     GetEnv("ES_PORT", "9200"),
		ESUsername: GetEnv("ES_USERNAME", "elastic"),
		ESPassword: GetEnv("ES_PASSWORD", ""),
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
