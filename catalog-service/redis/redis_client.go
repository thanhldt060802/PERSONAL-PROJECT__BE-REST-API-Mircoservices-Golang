package redis

import (
	"context"
	"fmt"
	"log"
	"thanhldt060802/config"

	"github.com/redis/go-redis/v9"
)

func NewClient() *redis.Client {
	ctx := context.Background()

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.AppConfig.DBRedisHost, config.AppConfig.DBRedisPort),
		Password: config.AppConfig.DBRedisPassword,
		DB:       0,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Connect to Redis failed: ", err)
	}
	log.Println("Connect to Redis successful")

	return client
}
