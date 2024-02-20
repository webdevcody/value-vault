package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/streadway/amqp"
)

var (
	rabbitConn *amqp.Connection
	rabbitCh   *amqp.Channel
	connErr    error
)

// InitializeRabbitMQ initializes the RabbitMQ connection and channel.
func InitializeRabbitMQ() {
	rabbitConn, connErr = amqp.Dial(getRabbitMqUrl())
	if connErr != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", connErr)
	}
	rabbitCh, connErr = rabbitConn.Channel()
	if connErr != nil {
		log.Fatalf("Failed to open a channel: %v", connErr)
	}
}

var cache map[string][]byte = make(map[string][]byte)
var dirty map[string]bool = make(map[string]bool)
var cacheMutex = sync.RWMutex{}
var dirtyMutex = sync.RWMutex{}

func getRabbitMqUrl() string {
	var rabbitMqHost = os.Getenv("RABBIT_MQ_HOST")
	var rabbitMqPass = os.Getenv("RABBIT_MQ_PASSWORD")

	if rabbitMqPass == "" {
		log.Fatal("RABBIT_MQ_PASSWORD env variable is not set")
	}

	return "amqp://user:" + rabbitMqPass + "@" + rabbitMqHost + ":5672/"
}

func handleEvent(body []byte) {
	eventName := string(body)
	// Add your logic to handle the event here
	fmt.Println("setting cache as dirty for", eventName)
	dirtyMutex.Lock()
	dirty[eventName] = true
	dirtyMutex.Unlock()
}

func PublishEvent(topicName, event string) error {
	if rabbitConn == nil || rabbitCh == nil {
		InitializeRabbitMQ()
	}

	// Declare the topic exchange if it's not already declared
	err := rabbitCh.ExchangeDeclare(
		topicName, // exchange name
		"topic",   // exchange type
		true,      // durable
		false,     // auto-deleted
		false,     // internal
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return err
	}

	// Publish the event message to the topic exchange
	err = rabbitCh.Publish(
		topicName, // exchange
		"",        // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		},
	)
	if err != nil {
		return err
	}

	log.Printf("Event '%s' published to topic '%s'\n", event, topicName)
	return nil
}

func listenToRabbitMQ() {
	if rabbitConn == nil || rabbitCh == nil {
		InitializeRabbitMQ()
	}

	// Declare the exchange
	err := rabbitCh.ExchangeDeclare(
		"events", // exchange name
		"topic",  // exchange type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare the exchange: %v", err)
	}

	// Declare a queue
	q, err := rabbitCh.QueueDeclare(
		"",    // queue name (empty to let RabbitMQ generate a unique name)
		false, // durable
		true,  // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	// Bind the queue to the exchange
	err = rabbitCh.QueueBind(
		q.Name,   // queue name
		"#",      // routing key (listen to all topics)
		"events", // exchange name
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		log.Fatalf("Failed to bind the queue to the exchange: %v", err)
	}

	// Consume messages from the queue
	msgs, err := rabbitCh.Consume(
		q.Name, // queue name
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // arguments
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	// Read messages from the channel
	for msg := range msgs {
		handleEvent(msg.Body)
	}
}

func writeToJSONFile(filePath string, jsonData []byte) error {
	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("error writing JSON data to file: %v", err)
	}

	fmt.Printf("JSON data written to file %s\n", filePath)
	return nil
}

func readValueFromDisk(filePath string) ([]byte, error) {
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading JSON data from file: %v", err)
	}
	return jsonData, nil
}

func getKeyPath(key string) (string, error) {
	filePathPrefix := os.Getenv("FILE_PATH_PREFIX")
	if filePathPrefix == "" {
		return "", fmt.Errorf("FILE_PATH_PREFIX environment variable is not set")
	}
	return filePathPrefix + "/" + key + ".json", nil
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /keys/{key}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("\nGET /keys/{key}")

		key := r.PathValue("key")

		cacheMutex.RLock()
		cacheValue := cache[key]
		cacheMutex.RUnlock()

		if cacheValue == nil {
			keyPath, err := getKeyPath(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// check if keyPath exists on disk, and if it does load the value into cache
			if _, err := os.Stat(keyPath); err == nil {
				fmt.Println("keyPath exists on disk")
				jsonData, err := readValueFromDisk(keyPath)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				cacheMutex.Lock()
				cache[key] = jsonData
				cacheMutex.Unlock()

				fmt.Printf("Loaded data from file %s\n", keyPath)
			}
		}

		dirtyMutex.RLock()
		isDirty := dirty[key]
		dirtyMutex.RUnlock()

		if isDirty {
			keyPath, err := getKeyPath(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			jsonData, err := readValueFromDisk(keyPath)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			cacheMutex.Lock()
			cache[key] = jsonData
			cacheMutex.Unlock()

			dirtyMutex.Lock()
			dirty[key] = false
			dirtyMutex.Unlock()

			fmt.Printf("Loaded data from file %s\n", keyPath)
		}

		cacheMutex.RLock()
		value := cache[key]
		cacheMutex.RUnlock()

		if value == nil {
			http.Error(w, "key not found", http.StatusNotFound)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(value)
	})

	mux.HandleFunc("POST /keys/{key}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("\nPOST /keys/{key}")
		key := r.PathValue("key")

		var jsonData any
		err := json.NewDecoder(r.Body).Decode(&jsonData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// is it possible to convert jsonData to []byte?
		jsonDataBytes, err := json.Marshal(jsonData)

		cacheMutex.Lock()
		cache[key] = jsonDataBytes
		cacheMutex.Unlock()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		filePathPrefix := os.Getenv("FILE_PATH_PREFIX")
		if filePathPrefix == "" {
			http.Error(w, "FILE_PATH_PREFIX environment variable not set", http.StatusInternalServerError)
			return
		}

		filePath := filePathPrefix + "/" + key + ".json"

		if err := writeToJSONFile(filePath, jsonDataBytes); err != nil {
			http.Error(w, "could not write to file", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)

		PublishEvent("events", key)
	})

	go func() {
		fmt.Println("Starting RabbitMQ listener")
		listenToRabbitMQ()
	}()

	fmt.Println("Starting server on port 8080")
	http.ListenAndServe(":8080", mux)
}
