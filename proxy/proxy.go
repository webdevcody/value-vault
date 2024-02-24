package proxy

import (
	"bytes"
	"fmt"
	"io"
	"key-value-app/config"
	"net/http"
)

func setConfigurationOnHeaders(req *http.Request) {
	req.Header.Set("X-Configuration-Version", fmt.Sprintf("%d", config.GetConfiguration().Version))
	req.Header.Set("X-Configuration-Nodes", fmt.Sprintf("%d", config.GetConfiguration().CurrentNodeCount))
	req.Header.Set("X-Configuration-Previous-Nodes", fmt.Sprintf("%d", config.GetConfiguration().PreviousNodeCount))
}

func ForwardStoreToNode(hostname string, key string, value []byte) error {
	url := fmt.Sprintf("http://%s/keys/%s", hostname, key)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(value))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	setConfigurationOnHeaders(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func ForwardGetToNode(hostname string, key string, force bool) ([]byte, error) {
	suffix := ""
	if force {
		suffix = "?force=true"
	}
	url := fmt.Sprintf("http://%s/keys/%s%s", hostname, key, suffix)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	setConfigurationOnHeaders(req)

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

func DeleteKeyFromNode(hostname string, key string) ([]byte, error) {
	url := fmt.Sprintf("http://%s/keys/%s", hostname, key)

	req, err := http.NewRequest("DELETE", url, nil)
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
