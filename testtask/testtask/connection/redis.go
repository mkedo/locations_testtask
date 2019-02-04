package connection

import (
	"github.com/go-redis/redis"
	"os"
	"time"
)

func GetRedisConnection() *redis.Client {
	host, _ := os.LookupEnv("REDIS_ADDR")
	password, _ := os.LookupEnv("REDIS_PASSWORD")
	client := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       0,
		DialTimeout: 10 * time.Millisecond,
		ReadTimeout: 10 * time.Millisecond,
		WriteTimeout: 10 * time.Millisecond,
	})
	//if err := client.Ping().Err(); err != nil {
	//	panic(err)
	//}
	return client
}
