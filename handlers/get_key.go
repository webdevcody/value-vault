package handlers

import (
	"fmt"
	"key-value-app/hash"
	"key-value-app/locking"
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

	locking.RLock(key)
	defer locking.RUnlock(key)

	Log(context, fmt.Sprintf("GET /keys/%s", key))

	// currentConfiguration := config.GetConfiguration()

	node := hash.GetCurrentRingNode(key)
	// previousNode := hash.GetPreviousRingNode(key)

	currentNodeHostname := strings.Split(node.LogicalHostname, ".")[0]
	// previousNodeHostname := strings.Split(previousNode.Hostname, ".")[0]

	if currentNodeHostname == hostname {

		onDisk := persistence.IsKeyOnDisk(key)

		if onDisk {
			jsonData, err := persistence.ReadValueFromDisk(key)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-Stored-On", hostname)
			w.WriteHeader(http.StatusOK)
			w.Write(jsonData)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-Stored-On", hostname)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("null"))
		}
	} else {
		response, err := util.CallWithRetries(10, func() ([]byte, error) {
			return proxy.ForwardGetToNode(node.PhysicalHostname, key, context.traceId)
		})

		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		shortHostname := strings.Split(node.PhysicalHostname, ".")[0]
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Stored-On", shortHostname)
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}
