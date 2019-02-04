package main

import (
	"log"
	"testtask"
	"testtask/connection"
)

func main() {
	log.SetFlags(log.Lshortfile)

	redisClient := connection.GetRedisConnection()
	cache := testtask.NewRedisCache(redisClient)
	defer redisClient.Close()

	pgConnection := connection.GetPgConnection()
	pgRepo := testtask.NewPgStore(pgConnection)
	defer pgConnection.Close()

	itemLocations := testtask.NewCachedStore(&testtask.CachedStoreOptions{
		PgStore:           pgRepo,
		ItemLocationCache: cache,
	})


	log.Fatal(testtask.ServeStore(itemLocations))
}
