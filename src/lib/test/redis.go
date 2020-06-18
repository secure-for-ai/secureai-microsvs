package main

import (
	"context"
	"fmt"
	"os"
	"template2/lib/cache"
	"template2/lib/util"
)

var redisClient *cache.RedisClient

func initRedis() {
	var err error
	redisConf := cache.RedisConf{}
	redisConf.Addrs = []string{"localhost:6379"}
	redisConf.PW = "password"
	redisClient, err = cache.NewRedisClient(redisConf)

	if err != nil {
		fmt.Println("cannot connect to redis")
		os.Exit(1)
	}

	fmt.Println("connect to redis")
}

func main() {
	initRedis()

	key := "test"
	value, _ := util.GenerateRandomKey(10)
	ctx := context.Background()
	_, err := redisClient.Set(ctx, "test", value, 0)

	if err != nil {
		fmt.Println("cannot set", key, value)
	}

	fmt.Println("set", key, value)

	result, err := redisClient.Get(ctx, "test")

	if err != nil {
		fmt.Println("cannot get", key)
	}

	fmt.Println("get", key, result)

	count, err := redisClient.Del(ctx, "test")

	if err != nil {
		fmt.Println("cannot delete", key)
	}

	fmt.Println("delete", count, "key")

	result, err = redisClient.Get(ctx, "test")

	if err != nil {
		fmt.Println("cannot get", key)
	}

	fmt.Println("get", key, result)
}
