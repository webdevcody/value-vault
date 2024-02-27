package handlers

import (
	"encoding/json"
	"fmt"
	"key-value-app/cache"
	"key-value-app/config"
	"key-value-app/hash"
	"key-value-app/messaging"
	"key-value-app/persistence"
	"key-value-app/proxy"
	"key-value-app/util"
	"net/http"
	"os"
	"strings"
)

func StoreKey(w http.ResponseWriter, r *http.Request) {
	context := GetRequestContext(r)
	hostname := os.Getenv("HOSTNAME")

	key := r.PathValue("key")

	Log(context, fmt.Sprintf("POST /keys/%s", key))

	var jsonData any
	err := json.NewDecoder(r.Body).Decode(&jsonData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonDataBytes, err := json.Marshal(jsonData)

	requestConfiguration := GetConfigurationFromHeaders(r)
	currentConfiguration := config.GetConfiguration()
	if requestConfiguration != nil && requestConfiguration.Version > currentConfiguration.Version {
		fmt.Printf("new version found\n")

		config.SetConfiguration(requestConfiguration)
		hash.Reset()
	}

	// find the node for the key
	node := hash.GetCurrentRingNode(key)
	oldNode := hash.GetPreviousRingNode(key)
	nodeHostname := strings.Split(node.LogicalHostname, ".")[0]

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

		if !isOnDisk && oldNode.LogicalHostname != node.LogicalHostname {

			_, err := util.CallWithRetries(10, func() ([]byte, error) {
				return proxy.DeleteKeyFromNode(oldNode.PhysicalHostname, key, context.traceId)
			})

			if err != nil {
				fmt.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusCreated)

		messaging.PublishEvent(key, string(jsonDataBytes))
	} else {
		_, err := util.CallWithRetries(10, func() ([]byte, error) {
			return proxy.ForwardStoreToNode(node.PhysicalHostname, key, jsonDataBytes, context.traceId)
		})

		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}

	// messaging.PublishEvent(key)
}
