package main

import (
	"context"
	"fmt"
	"key-value-app/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /keys/{key}", handlers.GetKeys)
	mux.HandleFunc("POST /keys/{key}", handlers.StoreKey)
	mux.HandleFunc("DELETE /keys/{key}", handlers.DeleteKey)
	mux.HandleFunc("GET /status", handlers.StatusHandler)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)

		// interrupt signal sent from terminal
		signal.Notify(sigint, os.Interrupt)
		// sigterm signal sent from kubernetes
		signal.Notify(sigint, syscall.SIGTERM)

		<-sigint

		// sleep
		time.Sleep(10 * time.Second)

		// We received an interrupt signal, shut down.
		if err := server.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	// probe.Create()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Printf("HTTP server ListenAndServe: %v", err)
		// messaging.Initialize()
	}

	log.Printf("waiting for connections to close\n")

	<-idleConnsClosed

	// log.Printf("removing file\n")

	// probe.Remove()
}
