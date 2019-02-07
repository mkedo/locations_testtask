package testtask

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"testtask/store"
	"time"
)

type itemHandler func(w http.ResponseWriter, r *http.Request, itemId int64)

func handlerWithItemId(handler itemHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		itemIdStr := vars["ItemId"]
		itemId, err := strconv.ParseInt(itemIdStr, 10, 64)
		if err != nil {
			http.Error(w, "Wrong ItemId format", http.StatusBadRequest)
		} else {
			handler(w, r, itemId)
		}
	})
}

var itemRandom = rand.New(rand.NewSource(time.Now().UnixNano()))
var randomLock = sync.Mutex{}

func handlerWithRandomId(handler itemHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		randomLock.Lock()
		itemId := itemRandom.Int63n(200000)
		randomLock.Unlock()
		handler(w, r, itemId)
	})
}

func putItemLocationsHandler(itemLocations store.ItemLocations) itemHandler {
	type putItemLocationRequest struct {
		LocationIds []int64
	}
	return itemHandler(func(w http.ResponseWriter, r *http.Request, itemId int64) {
		var request putItemLocationRequest

		err := json.NewDecoder(r.Body).Decode(&request)
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

func getItemLocationsHandler(itemLocations store.ItemLocations) itemHandler {
	return itemHandler(func(w http.ResponseWriter, r *http.Request, itemId int64) {
		locations, err := itemLocations.GetContext(r.Context(), itemId)
		if err != nil {
			log.Println(err)
			http.Error(w, "Couldn't get locations", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
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
	r.Handle("/item/{ItemId}/locations", handlerWithItemId(getItemLocationsHandler(itemLocations))).Methods("GET")
	r.Handle("/item/{ItemId}/locations", handlerWithItemId(putItemLocationsHandler(itemLocations))).Methods("POST")

	// для теста
	r.Handle("/random_item/locations", handlerWithRandomId(getItemLocationsHandler(itemLocations))).Methods("GET")
	r.Handle("/random_item/locations", handlerWithRandomId(putItemLocationsHandler(itemLocations))).Methods("POST")

	srv := http.Server{
		Addr:    addr,
		Handler: r,
	}

	// graceful shutdown
	waitForShutdownChan := make(chan struct{})
	go func() {
		sigInt := make(chan os.Signal, 1)
		defer close(sigInt)
		signal.Notify(sigInt, os.Interrupt, syscall.SIGTERM)
		<-sigInt
		signal.Stop(sigInt)

		log.Println("Shutdown requested")
		// 3 секунды должно хватить всем
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("HTTP server Shutdown: %v\n", err)
		}
		close(waitForShutdownChan)
	}()

	log.Printf("Serving at %s\n", addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("HTTP server ListenAndServe: %v\n", err)
		return err
	} else {
		<-waitForShutdownChan
		log.Println("Web server is down")
		return nil
	}
}
