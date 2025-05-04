package infrastructure

import (
	"context"
	"fmt"
	"log"
	"thanhldt060802/config"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedisClient() {
	ctx := context.Background()

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.AppConfig.DBRedisHost, config.AppConfig.DBRedisPort),
		Password: config.AppConfig.DBRedisPassword,
		DB:       0,
	})

	if _, err := RedisClient.Ping(ctx).Result(); err != nil {
		log.Fatal("Connect to Redis failed: ", err)
	}
	log.Println("Connect to Redis successful")
}
