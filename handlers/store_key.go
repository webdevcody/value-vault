package handlers

import (
	"encoding/json"
	"fmt"
	"key-value-app/hash"
	"key-value-app/messaging"
	"key-value-app/persistence"
	"key-value-app/proxy"
	"key-value-app/util"
	"net/http"
	"os"
	"strings"
)

func StoreKeyValue(key string, value []byte, traceId string) error {
	hostname := os.Getenv("HOSTNAME")

	if traceId == "" {
		traceId = GenerateRandomString()
	}

	node := hash.GetCurrentRingNode(key)
	nodeHostname := strings.Split(node.LogicalHostname, ".")[0]

	isDataOnThisNode := nodeHostname == hostname

	if !isDataOnThisNode {
		_, err := util.CallWithRetries(10, func() ([]byte, error) {
			return proxy.ForwardStoreToNode(node.PhysicalHostname, key, value, traceId)
		})

		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	if err := persistence.WriteJsonToDisk(key, value); err != nil {
		return err
	}

	messaging.PublishEvent(key, string(value))

	return nil
}

func StoreKeyHandler(w http.ResponseWriter, r *http.Request) {
	context := GetRequestContext(r)

	key := r.PathValue("key")

	Log(context, fmt.Sprintf("POST /keys/%s", key))

	var jsonData any
	err := json.NewDecoder(r.Body).Decode(&jsonData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonDataBytes, err := json.Marshal(jsonData)

	if err != nil {
		http.Error(w, "could not marshal json", http.StatusInternalServerError)
		return
	}

	err = StoreKeyValue(key, jsonDataBytes, context.traceId)

	if err != nil {
		http.Error(w, "could not store key value", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
