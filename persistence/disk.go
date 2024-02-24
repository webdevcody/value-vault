package persistence

import (
	"fmt"
	"log"
	"os"
)

func WriteJsonToDisk(key string, jsonData []byte) error {
	filePath := getFilePathUsingKey(key)

	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("error writing JSON data to file: %v", err)
	}

	return nil
}

func DeleteKey(key string) error {
	filePath := getFilePathUsingKey(key)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil
	}

	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("error deleting key: %v", err)
	}

	return nil
}

func IsKeyOnDisk(key string) bool {
	filePath := getFilePathUsingKey(key)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
}

func ReadValueFromDisk(key string) ([]byte, error) {
	filePath := getFilePathUsingKey(key)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, nil
	}

	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading JSON data from file: %v", err)
	}

	return jsonData, nil
}

func getFilePathUsingKey(key string) string {
	filePathPrefix := os.Getenv("FILE_PATH_PREFIX")
	if filePathPrefix == "" {
		log.Fatal("FILE_PATH_PREFIX environment variable is not set")
	}
	return filePathPrefix + "/" + key + ".json"
}
