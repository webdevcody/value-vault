package main

import (
	"context"
	"fmt"
	"key-value-app/handlers"
	"key-value-app/messaging"
	"key-value-app/persistence"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	persistence.Initialize()

	messaging.Initialize()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /keys/{key}", handlers.GetKeys)
	mux.HandleFunc("POST /keys/{key}", handlers.StoreKeyHandler)
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
		signal.Notify(sigint, os.Interrupt)
		signal.Notify(sigint, syscall.SIGTERM)
		<-sigint

		// important, we need this sleep to wait for the k8s iptables to update before it
		// starts to tear this pod down...
		time.Sleep(10 * time.Second)

		// We received an interrupt signal, shut down.
		if err := server.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	go func() {
		fmt.Println("Starting Replication listener")
		messaging.InitializeReplicationEventListener()
	}()

	go func() {
		fmt.Println("Listening to Redistribute Listener")
		messaging.InitializeRedistributedEventListener()
	}()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("HTTP server ListenAndServe: %v", err)
	}

	log.Printf("waiting for connections to close!\n")

	<-idleConnsClosed

}
