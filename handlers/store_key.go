package handlers

import (
	"encoding/json"
	"fmt"
	"key-value-app/cache"
	"key-value-app/config"
	"key-value-app/hash"
	"key-value-app/persistence"
	"key-value-app/proxy"
	"net/http"
	"os"
	"strings"
)

func StoreKey(w http.ResponseWriter, r *http.Request) {
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

	requestConfiguration := GetConfigurationFromHeaders(r)
	currentConfiguration := config.GetConfiguration()
	if requestConfiguration != nil {
		fmt.Printf("incoming version %d\n", requestConfiguration.Version)
	}
	if requestConfiguration != nil && requestConfiguration.Version > currentConfiguration.Version {
		fmt.Printf("new version found\n")

		config.SetConfiguration(requestConfiguration)
		hash.Reset()
	}

	// find the node for the key
	node := hash.GetCurrentRingNode(key)
	oldNode := hash.GetPreviousRingNode(key)
	nodeHostname := strings.Split(node.Hostname, ".")[0]

	fmt.Printf("new node %s\n", node.Hostname)

	// print node hostname and HOSTNAME env
	if nodeHostname == hostname {

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		cache.StoreIntoCache(key, jsonDataBytes)

		// if the key doesn't exist here, we will need to delete from old node
		isOnDisk := persistence.IsKeyOnDisk(key)

		if err := persistence.WriteJsonToDisk(key, jsonDataBytes); err != nil {
			http.Error(w, "could not write to file", http.StatusInternalServerError)
			return
		}

		if !isOnDisk && oldNode.Hostname != node.Hostname {
			_, err := proxy.DeleteKeyFromNode(oldNode.Hostname, key)
			if err != nil {
				fmt.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
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
}
