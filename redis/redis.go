package redis

import (
	"context"
	"log"
	"time"

	"github.com/PI-Team04-GameClub/gameclub-backend/config"
	goredis "github.com/redis/go-redis/v9"
)

var Client *goredis.Client

func Connect(cfg *config.Config) {
	Client = goredis.NewClient(&goredis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := Client.Ping(ctx).Err(); err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v. Caching disabled.", err)
		Client = nil
		return
	}
	log.Println("Redis connected successfully")
}

func Close() {
	if Client != nil {
		if err := Client.Close(); err != nil {
			log.Printf("Error closing Redis connection: %v", err)
		}
	}
}
