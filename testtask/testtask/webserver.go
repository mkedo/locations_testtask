package testtask

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"testtask/store"
)

func findLocationHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "To be implemented", http.StatusNotFound)
	})
}

func putItemLocationsHandler(itemLocations store.ItemLocations) http.Handler {
	type putItemLocationRequest struct {
		ItemId      int64
		LocationIds []int64
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "POST required", http.StatusBadRequest)
			return
		}
		var request putItemLocationRequest

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil || request.ItemId == 0 {
			http.Error(w, "Wrong request format", http.StatusBadRequest)
			return
		}
		//TODO: если есть не уникальные LocationId то удалить дубликаты или кинуть ошибку?

		err = itemLocations.PutContext(r.Context(), request.ItemId, request.LocationIds)
		if err != nil {
			log.Println(err)
			http.Error(w, "Couldn't put items", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, `{"ok": 1}`)
	})
}

func getItemLocationsHandler(itemLocations store.ItemLocations) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "GET required", http.StatusBadRequest)
			return
		}
		itemIdStr, ok := r.URL.Query()["ItemId"]
		if !ok || len(itemIdStr[0]) < 1 {
			http.Error(w, "ItemId required", http.StatusBadRequest)
			return
		}
		itemId, err := strconv.ParseInt(itemIdStr[0], 10, 64)
		if err != nil {
			http.Error(w, "Wrong ItemId format", http.StatusBadRequest)
			return
		}
		locations, err := itemLocations.GetContext(r.Context(), itemId)
		if err != nil {
			log.Println(err)
			http.Error(w, "Couldn't get locations", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(locations)
	})
}

func ServeStore(itemLocations store.ItemLocations) error {
	addr := ":8080"
	//http.Handle("/findLocation", findLocationHandler())
	http.Handle("/putItemLocations", putItemLocationsHandler(itemLocations))
	http.Handle("/getItemLocations", getItemLocationsHandler(itemLocations))
	log.Printf("Serving at %s\n", addr)
	return http.ListenAndServe(addr, nil)
}
