package proxy

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func ForwardStoreToNode(hostname string, key string, value []byte) error {
	url := fmt.Sprintf("http://%s/keys/%s", hostname, key)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(value))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func ForwardGetToNode(hostname string, key string) ([]byte, error) {
	url := fmt.Sprintf("http://%s/keys/%s", hostname, key)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
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
