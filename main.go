package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/streadway/amqp"
)

var cache map[string][]byte = make(map[string][]byte)
var dirty map[string]bool = make(map[string]bool)

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
	dirty[eventName] = true
}

func PublishEvent(topicName, event string) error {
	// Establish connection to RabbitMQ server
	conn, err := amqp.Dial(getRabbitMqUrl())
	if err != nil {
		return err
	}
	defer conn.Close()

	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// Declare the topic exchange
	err = ch.ExchangeDeclare(
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
	err = ch.Publish(
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
	conn, err := amqp.Dial(getRabbitMqUrl())
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Open a channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// Declare the exchange
	err = ch.ExchangeDeclare(
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
	q, err := ch.QueueDeclare(
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
	err = ch.QueueBind(
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
	msgs, err := ch.Consume(
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
	// jsonData, err := json.MarshalIndent(jsonString, "", "  ")
	// if err != nil {
	// 	return fmt.Errorf("error marshaling JSON data: %v", err)
	// }
	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("error writing JSON data to file: %v", err)
	}

	fmt.Printf("JSON data written to file %s\n", filePath)
	return nil
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
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

		if cache[key] == nil {
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
				cache[key] = jsonData
				fmt.Printf("Loaded data from file %s\n", keyPath)
			}
		}

		if dirty[key] {
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
			cache[key] = jsonData
			dirty[key] = false
			fmt.Printf("Loaded data from file %s\n", keyPath)
		}

		value := cache[key]
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
		cache[key] = jsonDataBytes

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

		fmt.Printf("Event '%s' published to topic 'events'\n", key)
		PublishEvent("events", key)
	})

	go func() {
		fmt.Println("Starting RabbitMQ listener")
		listenToRabbitMQ()
	}()

	fmt.Println("Starting server on port 8080")
	http.ListenAndServe(":8080", mux)
}
