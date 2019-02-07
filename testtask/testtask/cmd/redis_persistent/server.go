package main

import (
	"log"
	"testtask"
	"testtask/connection"
)

func main() {
	log.SetFlags(log.Lshortfile)

	pgConnection := connection.GetPgConnection()
	defer pgConnection.Close()

	redisClient := connection.GetRedisConnection()
	defer redisClient.Close()

	itemLocations := testtask.NewPgRedisPersistent(
		pgConnection,
		redisClient,
	)

	if err := testtask.ServeStore(itemLocations); err != nil {
		log.Fatal(err)
	}
}
