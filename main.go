package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var store map[string]any = make(map[string]any)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /key/{key}", func(w http.ResponseWriter, r *http.Request) {
		key := r.PathValue("key")

		value, err := json.Marshal(store[key])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(value)
	})

	mux.HandleFunc("POST /key/{key}", func(w http.ResponseWriter, r *http.Request) {
		key := r.PathValue("key")

		fmt.Println(r.Body)
		var jsonData any
		err := json.NewDecoder(r.Body).Decode(&jsonData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		store[key] = jsonData

		w.WriteHeader(http.StatusCreated)
	})

	fmt.Println("Starting server on port 8080")
	http.ListenAndServe(":8080", mux)
}
