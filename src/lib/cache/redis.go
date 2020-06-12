package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

type RedisConf struct {
	Addrs []string `json:"Addrs"`
	PW    string   `json:"PW"`
}

//RedisClusterClient struct
type RedisClient struct {
	rdb redis.UniversalClient
}

func NewRedisClient(conf RedisConf) (client *RedisClient, err error) {
	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    conf.Addrs,
		Password: conf.PW,
		DB:       0, // use default DB
	})
	rdb.Context()
	pong, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	log.Println(pong, "Redis Success!")

	client.rdb = rdb

	return client, err
}

func (c *RedisClient) Close() error {
	return c.rdb.Close()
}

func (c *RedisClient) GetClient() redis.UniversalClient {
	return c.rdb
}

func (c *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return c.rdb.Get(ctx, key).Result()
}

func (c *RedisClient) GetInt64(ctx context.Context, key string) (int64, error) {
	return c.rdb.Get(ctx, key).Int64()
}

func (c *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.rdb.Set(ctx, key, value, expiration).Err()
}

func (c *RedisClient) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.rdb.SetNX(ctx, key, value, expiration).Err()
}
