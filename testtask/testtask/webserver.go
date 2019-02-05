package testtask

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strconv"
	"testtask/store"
)

func findLocationHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "To be implemented", http.StatusNotFound)
	})
}
// Парсит параметр с id объявления.
func getItemId(r *http.Request) (int64, error) {
	vars := mux.Vars(r)
	itemIdStr := vars["ItemId"]
	itemId, err := strconv.ParseInt(itemIdStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return itemId, nil
}

func putItemLocationsHandler(itemLocations store.ItemLocations) http.Handler {
	type putItemLocationRequest struct {
		LocationIds []int64
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		itemId, err := getItemId(r)
		if err != nil {
			http.Error(w, "Wrong ItemId format", http.StatusBadRequest)
			return
		}
		var request putItemLocationRequest

		err = json.NewDecoder(r.Body).Decode(&request)
		if err != nil || itemId == 0 {
			http.Error(w, "Wrong request format", http.StatusBadRequest)
			return
		}
		//TODO: если есть не уникальные LocationId то удалить дубликаты или кинуть ошибку?

		err = itemLocations.PutContext(r.Context(), itemId, request.LocationIds)
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
		itemId, err := getItemId(r)
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

// Запустить веб-север с заданным хранилищем.
func ServeStore(itemLocations store.ItemLocations) error {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)

	r := mux.NewRouter()
	r.Handle("/item/{ItemId}/locations", getItemLocationsHandler(itemLocations)).Methods("GET")
	r.Handle("/item/{ItemId}/locations", putItemLocationsHandler(itemLocations)).Methods("POST")
	//r.Handle("/findLocation", findLocationHandler()).Methods("GET")
	log.Printf("Serving at %s\n", addr)
	return http.ListenAndServe(addr, r)
}
