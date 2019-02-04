package main

import (
	"context"
	"fmt"
	"log"
	"testtask"
	"testtask/connection"
	"testtask/populate"
	"testtask/store"
)

func main() {
	log.SetFlags(log.Lshortfile)
	pgConnection := connection.GetPgConnection()
	defer pgConnection.Close()
	itemLocations := testtask.NewPgStore(pgConnection)

	maxLocationId := int64(100000)
	fmt.Println("Population started...")
	bulkSize := int64(100)
	fmt.Println("Populating locations...")
	for current := int64(1); current <= maxLocationId; current += bulkSize {
		locations := populate.GetRandomLocations(current, current+bulkSize-1)
		err := itemLocations.Add(context.Background(), locations)
		if err != nil {
			panic(err)
		}
	}

	itemLocations = nil

	redisClient := connection.GetRedisConnection()
	defer redisClient.Close()

	locations := testtask.NewPgRedisPersistent(
		pgConnection,
		redisClient,
	)

	fmt.Println("Populating item-location links...")
	for current := int64(1); current <= 500000; current++ {
		locIds := populate.GetRandomLocationIds(1, maxLocationId)
		err := locations.PutContext(context.Background(), store.ItemId(current), locIds)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Population done")
}
