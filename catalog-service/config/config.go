package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string

	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string

	RedisHost     string
	RedisPort     string
	RedisPassword string

	ElasticsearchHost     string
	ElasticsearchPort     string
	ElasticsearchUsername string
	ElasticsearchPassword string
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

		RedisHost:     GetEnv("REDIS_HOST", "localhost"),
		RedisPort:     GetEnv("REDIS_PORT", "6379"),
		RedisPassword: GetEnv("REDIS_PASSWORD", ""),

		ElasticsearchHost:     GetEnv("ELASTICSEARCH_HOST", "localhost"),
		ElasticsearchPort:     GetEnv("ELASTICSEARCH_PORT", "9200"),
		ElasticsearchUsername: GetEnv("ELASTICSEARCH_USERNAME", "elastic"),
		ElasticsearchPassword: GetEnv("ELASTICSEARCH_PASSWORD", ""),
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
