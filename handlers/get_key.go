package handlers

import (
	"fmt"
	"key-value-app/config"
	"key-value-app/hash"
	"key-value-app/persistence"
	"key-value-app/proxy"
	"net/http"
	"os"
	"strings"
)

func GetKeys(w http.ResponseWriter, r *http.Request) {
	hostname := os.Getenv("HOSTNAME")

	key := r.PathValue("key")
	force := r.URL.Query().Get("force")

	suffix := ""
	if force == "true" {
		suffix = "?force=true"
	}

	fmt.Printf("\n\nHostname=%s GET /keys/%s%s\n", hostname, key, suffix)

	requestConfiguration := GetConfigurationFromHeaders(r)
	currentConfiguration := config.GetConfiguration()
	if requestConfiguration != nil {
		fmt.Printf("incoming version %d\n", requestConfiguration.Version)
	}
	if requestConfiguration != nil && requestConfiguration.Version > currentConfiguration.Version {
		config.SetConfiguration(requestConfiguration)
		hash.Reset()
	}

	node := hash.GetCurrentRingNode(key)
	previousNode := hash.GetPreviousRingNode(key)

	currentNodeHostname := strings.Split(node.Hostname, ".")[0]
	previousNodeHostname := strings.Split(previousNode.Hostname, ".")[0]

	fmt.Printf("Current Hostname: %s\n", hostname)
	fmt.Printf("Old Target Hostname: %s\n", previousNodeHostname)
	fmt.Printf("New Target Hostname: %s\n", currentNodeHostname)

	if currentNodeHostname == hostname || force == "true" {

		onDisk := persistence.IsKeyOnDisk(key)

		if onDisk || force == "true" {
			jsonData, err := persistence.ReadValueFromDisk(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if jsonData == nil {
				jsonData = []byte("null")
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(jsonData)
		} else {
			oldNode := hash.GetPreviousRingNode(key)
			fmt.Println("FORCE GET TO", oldNode.Hostname)
			response, err := proxy.ForwardGetToNode(oldNode.Hostname, key, true)
			if err != nil {
				fmt.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if err := persistence.WriteJsonToDisk(key, response); err != nil {
				http.Error(w, "could not write to file", http.StatusInternalServerError)
				return
			}

			_, deleteErr := proxy.DeleteKeyFromNode(oldNode.Hostname, key)
			if deleteErr != nil {
				http.Error(w, deleteErr.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(response)
		}

	} else {
		fmt.Println("Forwarding get to", node.Hostname)
		response, err := proxy.ForwardGetToNode(node.Hostname, key, false)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}
