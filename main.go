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
	"strings"
	"syscall"
	"time"
)

func handleEvent(message string) {
	parts := strings.Split(message, "|YOLO|")
	hostname := parts[0]
	key := parts[1]
	body := parts[2]

	log.Printf("%s wrote value for key %s\n", hostname, key)

	if os.Getenv("HOSTNAME") != hostname {
		if err := persistence.WriteJsonToDisk(key, []byte(body)); err != nil {
			log.Fatal("could not write to file", err)
			return
		}
	}
}

func main() {

	messaging.Initialize()

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

	go func() {
		fmt.Println("Starting RabbitMQ listener")
		messaging.InitializeEventListener(handleEvent)
	}()

	// probe.Create()
	messaging.Initialize()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Printf("HTTP server ListenAndServe: %v", err)
	}

	log.Printf("waiting for connections to close\n")

	<-idleConnsClosed

	// log.Printf("removing file\n")

	// probe.Remove()
}
