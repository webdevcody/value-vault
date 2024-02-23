package main

import (
	"encoding/json"
	"fmt"
	"key-value-app/cache"
	"key-value-app/hash"
	"key-value-app/persistence"
	"key-value-app/proxy"
	"net/http"
	"os"
	"strings"
)

func handleEvent(eventName string) {
	cache.SetDirty(eventName)
}

func main() {

	// messaging.Initialize()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /status", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("GET /keys/{key}", func(w http.ResponseWriter, r *http.Request) {
		hostname := os.Getenv("HOSTNAME")

		key := r.PathValue("key")

		fmt.Printf("Hostname=%s GET /keys/%s\n", hostname, key)

		node := hash.GetNode(key)

		nodeHostname := strings.Split(node.Hostname, ".")[0]

		if nodeHostname == hostname {

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

		} else {
			fmt.Println("Forwarding get to", node.Hostname)
			response, err := proxy.ForwardGetToNode(node.Hostname, key)
			if err != nil {
				fmt.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(response)
		}
	})

	mux.HandleFunc("POST /keys/{key}", func(w http.ResponseWriter, r *http.Request) {
		hostname := os.Getenv("HOSTNAME")

		key := r.PathValue("key")

		fmt.Printf("Hostname=%s POST /keys/%s\n", hostname, key)

		var jsonData any
		err := json.NewDecoder(r.Body).Decode(&jsonData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonDataBytes, err := json.Marshal(jsonData)

		// find the node for the key
		node := hash.GetNode(key)

		nodeHostname := strings.Split(node.Hostname, ".")[0]

		// print node hostname and HOSTNAME env
		if nodeHostname == hostname {

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			cache.StoreIntoCache(key, jsonDataBytes)

			if err := persistence.WriteJsonToDisk(key, jsonDataBytes); err != nil {
				http.Error(w, "could not write to file", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusCreated)

		} else {
			fmt.Println("Forwarding store to", node.Hostname)
			err := proxy.ForwardStoreToNode(node.Hostname, key, jsonDataBytes)
			if err != nil {
				// print error
				fmt.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)
		}

		// messaging.PublishEvent(key)
	})

	// go func() {
	// 	fmt.Println("Starting RabbitMQ listener")
	// 	messaging.InitializeEventListener(handleEvent)
	// }()

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	fmt.Println("Starting server on port %s", port)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), mux)
	if err != nil {
		panic(err)
	}

}
