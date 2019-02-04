package main

import (
	"fmt"
	"log"
	"testtask"
	"testtask/connection"
	"testtask/populate"
	"testtask/store"
)

func main() {
	log.SetFlags(log.Lshortfile)
	client := connection.GetTntConnection()
	defer client.Close()
	itemLocations := testtask.NewTntStore(client)

	maxLocationId := int64(100000)
	fmt.Println("Population started...")
	bulkSize := int64(100)
	fmt.Println("Populating locations...")
	for current := int64(1); current <= maxLocationId; current += bulkSize {
		locations := populate.GetRandomLocations(current, current+bulkSize-1)
		err := itemLocations.Add(locations)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Populating item-location links...")
	for current := int64(1); current <= 500000; current++ {
		locIds := populate.GetRandomLocationIds(1, maxLocationId)
		err := itemLocations.Put(store.ItemId(current), locIds)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Population done")
}
