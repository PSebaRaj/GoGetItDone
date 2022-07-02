package cache

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-redis/redis"
	"github.com/psebaraj/gogetitdone/utils"
)

var REDIS *redis.Client

func getCacheUri() string {
	utils.LoadEnvVars()

	var cacheHost = os.Getenv("REDIS_HOST")
	var cachePort = os.Getenv("REDIS_PORT")
	return fmt.Sprintf("%s:%s", cacheHost, cachePort)
}

func ConnectRedisCache() {
	utils.LoadEnvVars()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     getCacheUri(),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	if _, redis_err := redisClient.Ping().Result(); redis_err != nil {
		fmt.Println(redis_err.Error())
		panic("Error: Unable to connect to Redis")
	}
	REDIS = redisClient
	fmt.Println("Connected to Redis cache successfully")
}

func SetInCache(c *redis.Client, key string, value interface{}) bool {
	marshalledValue, err := json.Marshal(value)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Unable to set element in cache")
		return false
	}
	c.Set(key, marshalledValue, 0)
	return true
}

func GetFromCache(c *redis.Client, key string) interface{} {
	value, err := c.Get(key).Result()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return value
}

func DeleteFromCache(c *redis.Client, key string) {
	c.Del(key)
}
