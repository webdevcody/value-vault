package hash

import (
	"fmt"
	"key-value-app/config"
	"os"
	"sort"
	"strconv"
	"sync"

	"github.com/spaolacci/murmur3"
)

var weight = 10
var currentRingHash *ConsistentHash
var previousRingHash *ConsistentHash
var currentRingMutex sync.Mutex
var previousRingMutex sync.Mutex

type HashRing []uint32

type Node struct {
	Id       int
	Hostname string
}

type ConsistentHash struct {
	Nodes      map[uint32]Node
	SortedKeys HashRing
	IsSorted   bool
}

func NewConsistentHash() *ConsistentHash {
	return &ConsistentHash{
		Nodes:      make(map[uint32]Node),
		SortedKeys: make(HashRing, 0),
		IsSorted:   false,
	}
}

func (ch *ConsistentHash) AddNode(hostname string) {
	node := Node{
		Hostname: hostname,
	}

	for i := 0; i < weight; i++ {
		hash := ch.hashKey(hostname + strconv.Itoa(i))
		ch.Nodes[hash] = node
	}

	ch.IsSorted = false
}

func (ch *ConsistentHash) hashKey(key string) uint32 {
	hasher := murmur3.New32()
	hasher.Write([]byte(key))
	return hasher.Sum32()
}

func (ch *ConsistentHash) GetNode(key string) Node {
	if len(ch.Nodes) == 0 {
		return Node{}
	}
	hash := ch.hashKey(key)
	if !ch.IsSorted {
		ch.sortNodes()
	}

	for _, k := range ch.SortedKeys {
		if k >= hash {
			return ch.Nodes[k]
		}
	}
	node := ch.Nodes[ch.SortedKeys[0]]
	return node
}

func (ch *ConsistentHash) sortNodes() {
	// Collect keys (hashes) from the map
	keys := make([]uint32, 0, len(ch.Nodes))
	for k := range ch.Nodes {
		keys = append(keys, k)
	}

	// Sort the keys in ascending order
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	sortedKeys := make(HashRing, 0)
	for _, k := range keys {
		sortedKeys = append(sortedKeys, k)
	}

	ch.SortedKeys = sortedKeys

	ch.IsSorted = true
}

func Reset() {
	currentRingMutex.Lock()
	defer currentRingMutex.Unlock()
	previousRingMutex.Lock()
	defer previousRingMutex.Unlock()
	currentRingHash = nil
	previousRingHash = nil
}

func getNodeFromRing(key string, ring **ConsistentHash, totalNodes int) Node {
	if *ring == nil {
		*ring = NewConsistentHash()
		isLocal := os.Getenv("IS_LOCAL") == "true"

		for id := range totalNodes {
			if isLocal {
				(*ring).AddNode(fmt.Sprintf("localhost:%d", 8080+id))
			} else {
				(*ring).AddNode(fmt.Sprintf("api-%d.api.default.svc.cluster.local:8080", id))
			}
		}
	}

	return (*ring).GetNode(key)
}

func GetCurrentRingNode(key string) Node {
	configuration := config.GetConfiguration()
	currentRingMutex.Lock()
	defer currentRingMutex.Unlock()
	return getNodeFromRing(key, &currentRingHash, configuration.CurrentNodeCount)
}

func GetPreviousRingNode(key string) Node {
	configuration := config.GetConfiguration()
	previousRingMutex.Lock()
	defer previousRingMutex.Unlock()
	return getNodeFromRing(key, &previousRingHash, configuration.PreviousNodeCount)
}
