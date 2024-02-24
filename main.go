package main

import (
	"fmt"
	"key-value-app/handlers"
	"net/http"
	"os"
)

func main() {
	// messaging.Initialize()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /keys/{key}", handlers.GetKeys)
	mux.HandleFunc("POST /keys/{key}", handlers.StoreKey)
	mux.HandleFunc("DELETE /keys/{key}", handlers.DeleteKey)

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
