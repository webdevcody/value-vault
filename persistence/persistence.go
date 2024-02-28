package persistence

import (
	"fmt"
	"os"
)

var directoryPath string

func WriteJsonToDisk(key string, jsonData []byte) error {
	filePath := getFilePathUsingKey(key)

	fmt.Printf("%s\n", directoryPath)

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
	return directoryPath + "/" + key + ".json"
}

func Initialize() error {
	// Check if the directory already exists
	directoryPath = fmt.Sprintf("%s/%s", os.Getenv("FILE_PATH_PREFIX"), os.Getenv("HOSTNAME"))
	fmt.Printf("%s\n", directoryPath)
	_, err := os.Stat(directoryPath)
	if os.IsNotExist(err) {
		// Directory doesn't exist, create it
		err := os.MkdirAll(directoryPath, 0755) // 0755 is the default permissions
		if err != nil {
			return err
		}
	} else if err != nil {
		// Error occurred while checking directory existence
		return err
	}

	// Directory already exists or created successfully
	return nil
}
