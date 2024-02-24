package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type Configuration struct {
	CurrentNodeCount  int
	PreviousNodeCount int
	Version           int
}

var configuration *Configuration

func getIntFromEnv(envString string) int {
	nodes := os.Getenv(envString)
	nodesInt, err := strconv.Atoi(nodes)
	if err != nil {
		log.Fatal(err)
	}
	return nodesInt
}

func GetConfiguration() *Configuration {
	if configuration == nil {
		configuration = &Configuration{
			CurrentNodeCount:  getIntFromEnv("NODES"),
			PreviousNodeCount: getIntFromEnv("PREVIOUS_NODES"),
			Version:           getIntFromEnv("CONFIG_VERSION"),
		}
	}
	return configuration
}

// this will be called if a header comes in with data containing a new configuration version
func SetConfiguration(config *Configuration) {
	configuration.CurrentNodeCount = config.CurrentNodeCount
	configuration.PreviousNodeCount = config.PreviousNodeCount
	configuration.Version = config.Version
	fmt.Printf("%d %d %d\n", configuration.CurrentNodeCount, configuration.PreviousNodeCount, configuration.Version)
}
