package probe

import "os"

const liveFile = "/tmp/live"

// Create will create a file for the liveness check
func Create() error {
	_, err := os.Create(liveFile)
	return err
}

// Remove will remove the file create for the liveness probe
func Remove() error {
	return os.Remove(liveFile)
}
