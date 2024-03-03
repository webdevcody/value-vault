package messaging

import (
	"key-value-app/hash"
	"key-value-app/locking"
	"key-value-app/persistence"
	"log"
	"os"
	"strings"
)

func handleRedistributeEvent(message string) {
	log.Printf("Handling redistribute event: %s", message)

	// if we get this event, loop over all files in the data directory and re-publish them to replication queue
	keys, err := persistence.GetAllKeys()
	if err != nil {
		log.Fatalf("could not get all files")
	}

	for i := range keys {
		key := keys[i]
		value, err := persistence.ReadValueFromDisk(key)
		if err != nil {
			log.Fatalf("could not read value from disk")
		}

		if err := PublishEvent(key, string(value)); err != nil {
			log.Fatalf("could not publish event")
		}
	}
}

func handleReplicationEvent(message string) {
	log.Printf("Handling replication event: %s", message)

	parts := strings.Split(message, "|YOLO|")

	fromMode := parts[0]
	key := parts[1]
	value := parts[2]

	// we should skip replication events from the same mode (primary / secondary)
	if fromMode == os.Getenv("MODE") {
		return
	}

	locking.Lock(key)
	defer locking.Unlock(key)

	hostname := os.Getenv("HOSTNAME")
	node := hash.GetCurrentRingNode(key)
	nodeHostname := strings.Split(node.LogicalHostname, ".")[0]

	isDataOnThisNode := nodeHostname == hostname

	if !isDataOnThisNode {
		return
	}

	if err := persistence.WriteJsonToDisk(key, []byte(value)); err != nil {
		log.Fatalf("could not write to disk")
	}
}
