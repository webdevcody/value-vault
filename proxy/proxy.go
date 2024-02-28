package proxy

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

var client = http.Client{
	Timeout: time.Second * 5,
}

func setTraceId(req *http.Request, traceId string) {
	req.Header.Set("X-Trace-Id", traceId)
}

func ForwardStoreToNode(hostname string, key string, value []byte, traceId string) ([]byte, error) {
	url := fmt.Sprintf("http://%s/keys/%s", hostname, key)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(value))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	setTraceId(req, traceId)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return []byte(""), nil
}

func ForwardGetToNode(hostname string, key string, traceId string) ([]byte, error) {
	suffix := ""
	url := fmt.Sprintf("http://%s/keys/%s%s", hostname, key, suffix)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	setTraceId(req, traceId)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return body, nil
}

func DeleteKeyFromNode(hostname string, key string, traceId string) ([]byte, error) {
	url := fmt.Sprintf("http://%s/keys/%s", hostname, key)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	setTraceId(req, traceId)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return body, nil
}
