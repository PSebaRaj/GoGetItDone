package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis"
	"github.com/psebaraj/gogetitdone/utils"
)

// DOCS: https://github.com/go-redis/redis/blob/master/commands.go

var REDIS *redis.Client // global redis cache

// Internal: gets the URI of the cache from env vars
func getCacheURI() string {
	utils.LoadEnvVars()

	cacheHost := os.Getenv("REDIS_HOST")
	cachePort := os.Getenv("REDIS_PORT")

	if cacheHost == "" || cachePort == "" {
		panic("Invalid cache URI")
	}

	return fmt.Sprintf("%s:%s", cacheHost, cachePort)
}

// from URI, connect to cache and create new redis client
func ConnectRedisCache() {
	utils.LoadEnvVars()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     getCacheURI(),
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

// set data in redis cache with specific key
// data expires after 60 mins
func SetInCache(c *redis.Client, key string, value interface{}) bool {

	// DO I need to check if redis client is valid/not bad?

	marshalledValue, err := json.Marshal(value)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Unable to marshal element to JSON")
		return false
	}

	_, err = c.Set(key, marshalledValue, 1*time.Hour).Result()
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Unable to set element in cache")
		return false

	}

	return true
}

// get data from redis cache with specific key
func GetFromCache(c *redis.Client, key string) interface{} {
	value, err := c.Get(key).Result()
	if err != nil {
		fmt.Println(err)
		fmt.Println("Unable to get element from cache")
		return nil
	}

	return value
}

// delete data from redis cache with specific key
func DeleteFromCache(c *redis.Client, key string) {
	_, err := c.Del(key).Result()
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error: Unable to delete element from cache")
	}

	return
}
