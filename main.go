package main

import (
	"encoding/json"
	"fmt"
	"key-value-app/cache"
	"key-value-app/messaging"
	"key-value-app/persistence"
	"net/http"
)

func handleEvent(eventName string) {
	cache.SetDirty(eventName)
}

func main() {

	messaging.Initialize()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /keys/{key}", func(w http.ResponseWriter, r *http.Request) {
		key := r.PathValue("key")

		fmt.Printf("GET /keys/%s\n", key)

		var cacheValue []byte

		isDirty := cache.IsDirty(key)
		if isDirty {
			jsonData, err := persistence.ReadValueFromDisk(key)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			cache.StoreIntoCache(key, jsonData)

			cache.UnsetDirty(key)

			cacheValue = jsonData
		} else {
			cacheValue = cache.GetFromCache(key)

			if cacheValue == nil {
				jsonData, err := persistence.ReadValueFromDisk(key)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				if jsonData == nil {
					jsonData = []byte("null")
				}
				cache.StoreIntoCache(key, jsonData)
				cacheValue = jsonData
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(cacheValue)
	})

	mux.HandleFunc("POST /keys/{key}", func(w http.ResponseWriter, r *http.Request) {
		key := r.PathValue("key")

		fmt.Printf("POST /keys/%s\n", key)

		var jsonData any
		err := json.NewDecoder(r.Body).Decode(&jsonData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonDataBytes, err := json.Marshal(jsonData)

		cache.StoreIntoCache(key, jsonDataBytes)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := persistence.WriteJsonToDisk(key, jsonDataBytes); err != nil {
			http.Error(w, "could not write to file", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		messaging.PublishEvent(key)
	})

	go func() {
		fmt.Println("Starting RabbitMQ listener")
		messaging.InitializeEventListener(handleEvent)
	}()

	fmt.Println("Starting server on port 8080")

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}

}
