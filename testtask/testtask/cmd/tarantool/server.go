package main

import (
	"log"
	"testtask"
	"testtask/connection"
)

func main() {
	log.SetFlags(log.Lshortfile)
	client := connection.GetTntConnection()
	defer client.Close()
	itemLocations := testtask.NewTntStore(client)
	if err := testtask.ServeStore(itemLocations); err != nil {
		log.Fatal(err)
	}
}
