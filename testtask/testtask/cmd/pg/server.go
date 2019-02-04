package main

import (
	"log"
	"testtask"
	"testtask/connection"
)

func main() {
	log.SetFlags(log.Lshortfile)

	pgConnection := connection.GetPgConnection()
	pgRepo := testtask.NewPgStore(pgConnection)
	defer pgConnection.Close()

	redisClient := connection.GetRedisConnection()
	cache := testtask.NewRedisCache(redisClient)
	defer redisClient.Close()

	itemLocations := testtask.NewCachedStore(&testtask.CachedStoreOptions{
		PgStore:           pgRepo,
		ItemLocationCache: cache,
	})
	//itemLocations := pgRepo

	log.Fatal(testtask.ServeStore(itemLocations))
}
