package connection

import (
	"github.com/go-redis/redis"
	"os"
)

func GetRedisConnection() *redis.Client {
	host, _ := os.LookupEnv("REDIS_ADDR")
	password, _ := os.LookupEnv("REDIS_PASSWORD")
	client := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       0,
	})
	//if err := client.Ping().Err(); err != nil {
	//	panic(err)
	//}
	return client
}
