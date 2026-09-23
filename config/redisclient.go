package config

import (
	"cmp"
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
)

var (
	Client      *redis.Client
	IsAvailable bool = false
	Ctx              = context.Background()
)

var contentGroup singleflight.Group

func RedisConnect() {
	addr := Env.RedisURL
	password := Env.RedisPassword
	db, _ := strconv.Atoi(cmp.Or(os.Getenv("REDIS_DB"), "0"))

	Client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if err := Client.Ping(context.Background()).Err(); err != nil {
		log.Println("Error connecting to redis: ", err)
	} else {
		IsAvailable = true
		log.Println("Redis Connection Successful")
	}
}

func Set(ctx context.Context, key string, value any, exp time.Duration) {
	_, err, _ := contentGroup.Do(key, func() (any, error) {
		// Double-check if another goroutine has already populated the cache
		// while we were waiting to enter the singleflight.Do call
		if _, err := Client.Get(ctx, key).Result(); err == nil {
			// Another goroutine already did the work
			return nil, nil
		}

		// Set to redis if cache is still empty
		err := Client.Set(ctx, key, value, exp).Err()
		if err != nil {
			log.Println("Error setting value", err)
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		log.Println("Error setting value", err)
	}
}

func Get(ctx context.Context, key string) string {
	val, err := Client.Get(ctx, key).Result()
	if err != nil {
		if err != redis.Nil {
			log.Println("Err Getting value", err)
		}
		return ""
	}
	return val
}
