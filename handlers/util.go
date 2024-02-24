package handlers

import (
	"key-value-app/config"
	"net/http"
	"strconv"
)

func GetConfigurationFromHeaders(r *http.Request) *config.Configuration {
	if r.Header.Get("X-Configuration-Version") == "" {
		return nil
	}

	// TODO: error handle, for now assume headers will always be set

	version, _ := strconv.Atoi(r.Header.Get("X-Configuration-Version"))
	nodes, _ := strconv.Atoi(r.Header.Get("X-Configuration-Nodes"))
	previousNodes, _ := strconv.Atoi(r.Header.Get("X-Configuration-Previous-Nodes"))
	return &config.Configuration{
		CurrentNodeCount:  nodes,
		PreviousNodeCount: previousNodes,
		Version:           version,
	}
}
