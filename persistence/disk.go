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

	fmt.Printf("JSON data written to file %s\n", filePath)
	return nil
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
