package handlers

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type RequestContext struct {
	traceId  string
	hostname string
}

func Log(context RequestContext, message string) {
	fmt.Printf("Trace=%s Hostname=%s Log=%s\n", context.traceId, context.hostname, message)
}

func GenerateRandomString() string {
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
		traceId = GenerateRandomString()
	}

	return traceId
}

func GetRequestContext(r *http.Request) RequestContext {
	return RequestContext{
		traceId:  getTraceId(r),
		hostname: os.Getenv("HOSTNAME"),
	}
}
