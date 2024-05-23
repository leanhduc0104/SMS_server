package cache

import (
	"context"
	"encoding/json"
	"log"
	"time"
	"vcs_server/entity"

	"github.com/go-redis/redis/v8"
)

type redisCache struct {
	host        string
	db          int
	expire_time time.Duration
}

func NewRedisCache(host string, db int, exp time.Duration) ServerCache {
	return &redisCache{
		host:        host,
		db:          db,
		expire_time: exp,
	}
}

func (cache *redisCache) getClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cache.host,
		Password: "", // no password set
		DB:       cache.db,
	})
}

func (cache *redisCache) Set(key string, value *entity.Server) {
	log.Println("Using Redis Set", key)
	ctx := context.Background()
	client := cache.getClient()

	json, err := json.Marshal(value)
	if err != nil {
		panic(json)
	}

	err = client.Set(ctx, key, string(json), 10*cache.expire_time*time.Second).Err()
	if err != nil {
		panic(err)
	}
}

func (cache *redisCache) Get(key string) *entity.Server {
	log.Println("Using Redis Get", key)
	ctx := context.Background()
	client := cache.getClient()
	val, err := client.Get(ctx, key).Result()
	if err != nil {
		if err != redis.Nil {
			log.Panic(err)
		}
		return nil
	}

	server := entity.Server{}
	err = json.Unmarshal([]byte(val), &server)
	if err != nil {
		panic(err)
	}

	return &server

}

func (cache *redisCache) ASet(key string, value []entity.Server) {
	log.Println("Using Redis ASet", key)
	ctx := context.Background()
	client := cache.getClient()

	json, err := json.Marshal(value)
	if err != nil {
		panic(json)
	}

	err = client.Set(ctx, key, string(json), cache.expire_time*time.Second).Err()
	if err != nil {
		panic(err)
	}
}

func (cache *redisCache) AGet(key string) []entity.Server {
	log.Println("Using Redis AGet", key)
	ctx := context.Background()
	client := cache.getClient()
	val, err := client.Get(ctx, key).Result()
	if err != nil {
		if err != redis.Nil {
			log.Panic(err)
		}
		return nil
	}

	var servers []entity.Server
	err = json.Unmarshal([]byte(val), &servers)
	if err != nil {
		panic(err)
	}

	return servers
}

func (cache *redisCache) Ping() {
	log.Println("Using Redis Ping")
	ctx := context.Background()
	client := cache.getClient()
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	log.Printf("Connected to Redis: %v\n", pong)
}

func (cache *redisCache) Del(key string) error {
	log.Println("Using Redis Del")
	ctx := context.Background()
	client := cache.getClient()
	err := client.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	log.Println("Del ", key)
	return nil
}
