package config

import (
	"fmt"
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
			CurrentNodeCount:  getIntFromEnv("NODES"),
			PreviousNodeCount: getIntFromEnv("PREVIOUS_NODES"),
			Version:           getIntFromEnv("CONFIG_VERSION"),
		}
	}
	return *configuration
}

// this will be called if a header comes in with data containing a new configuration version
func SetConfiguration(config *Configuration) {
	mutex.Lock()
	defer mutex.Unlock()
	configuration.CurrentNodeCount = config.CurrentNodeCount
	configuration.PreviousNodeCount = config.PreviousNodeCount
	configuration.Version = config.Version
	fmt.Printf("%d %d %d\n", configuration.CurrentNodeCount, configuration.PreviousNodeCount, configuration.Version)
}
