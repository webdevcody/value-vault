package handlers

import (
	"fmt"
	"key-value-app/config"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

type RequestContext struct {
	traceId  string
	hostname string
}

func Log(context RequestContext, message string) {
	fmt.Printf("Trace=%s Hostname=%s Log=%s\n", context.traceId, context.hostname, message)
}

func generateRandomString() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 32
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func getTraceId(r *http.Request) string {
	traceId := r.Header.Get("X-Trace-Id")

	if traceId == "" {
		traceId = generateRandomString()
	}

	return traceId
}

func GetRequestContext(r *http.Request) RequestContext {
	return RequestContext{
		traceId:  getTraceId(r),
		hostname: os.Getenv("HOSTNAME"),
	}
}

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
