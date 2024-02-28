package config

import (
	"log"
	"os"
	"strconv"
	"sync"
)

type Configuration struct {
	CurrentNodeCount  int
	PreviousNodeCount int
	Version           int
}

var configuration *Configuration
var mutex sync.Mutex

func getIntFromEnv(envString string) int {
	stringValue := os.Getenv(envString)
	value, err := strconv.Atoi(stringValue)
	if err != nil {
		log.Fatal(err)
	}
	return value
}

func GetConfiguration() Configuration {
	if configuration != nil {
		return *configuration
	}

	mutex.Lock()
	defer mutex.Unlock()

	if configuration == nil {
		configuration = &Configuration{
			CurrentNodeCount: getIntFromEnv("NODES"),
		}
	}
	return *configuration
}
