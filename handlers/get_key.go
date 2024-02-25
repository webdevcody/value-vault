package handlers

import (
	"fmt"
	"key-value-app/config"
	"key-value-app/hash"
	"key-value-app/persistence"
	"key-value-app/proxy"
	"key-value-app/util"
	"net/http"
	"os"
	"strings"
)

func GetKeys(w http.ResponseWriter, r *http.Request) {
	context := GetRequestContext(r)

	hostname := os.Getenv("HOSTNAME")

	key := r.PathValue("key")
	force := r.URL.Query().Get("force")

	suffix := ""
	if force == "true" {
		suffix = "?force=true"
	}

	Log(context, fmt.Sprintf("GET /keys/%s%s", key, suffix))

	requestConfiguration := GetConfigurationFromHeaders(r)
	currentConfiguration := config.GetConfiguration()
	if requestConfiguration != nil && requestConfiguration.Version > currentConfiguration.Version {
		config.SetConfiguration(requestConfiguration)
		hash.Reset()
	}

	node := hash.GetCurrentRingNode(key)
	// previousNode := hash.GetPreviousRingNode(key)

	currentNodeHostname := strings.Split(node.Hostname, ".")[0]
	// previousNodeHostname := strings.Split(previousNode.Hostname, ".")[0]

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

			response, err := util.CallWithRetries(10, func() ([]byte, error) {
				return proxy.ForwardGetToNode(oldNode.Hostname, key, true, context.traceId)
			})

			if err != nil {
				fmt.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if err := persistence.WriteJsonToDisk(key, response); err != nil {
				http.Error(w, "could not write to file", http.StatusInternalServerError)
				return
			}

			_, deleteErr := util.CallWithRetries(10, func() ([]byte, error) {
				return proxy.DeleteKeyFromNode(oldNode.Hostname, key, context.traceId)
			})

			if deleteErr != nil {
				http.Error(w, deleteErr.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(response)
		}

	} else {
		response, err := util.CallWithRetries(10, func() ([]byte, error) {
			return proxy.ForwardGetToNode(node.Hostname, key, false, context.traceId)
		})

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
