package hash

import (
	"fmt"
	"key-value-app/config"
	"os"
	"sort"
	"strconv"
	"strings"
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
	Id               int
	LogicalHostname  string
	PhysicalHostname string // this is a hack because I can't think of a better way to get this working locally
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

func (ch *ConsistentHash) AddNode(logicalHostname string, physicalHostname string) {
	fmt.Printf("adding node %s %s\n", logicalHostname, physicalHostname)
	node := Node{
		LogicalHostname:  logicalHostname,
		PhysicalHostname: physicalHostname,
	}

	for i := 0; i < weight; i++ {
		hash := ch.hashKey(physicalHostname + strconv.Itoa(i))
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
		mode := os.Getenv("MODE")

		for id := range totalNodes {
			if isLocal {
				localPorts := os.Getenv("LOCAL_PORTS")
				fmt.Printf("%s\n", localPorts)
				parts := strings.Split(os.Getenv("LOCAL_PORTS"), ";")
				logicalHostname := os.Getenv("HOSTNAME")
				lastIndex := len(logicalHostname) - 1
				strippedString := logicalHostname[:lastIndex]
				logicalHostname = fmt.Sprintf("%s%d", strippedString, id)

				port, _ := strconv.Atoi(parts[id])

				fmt.Printf("PORT = %d\n", port)
				fmt.Printf("ID = %d\n", id)
				fmt.Printf("SPLIT LEN = %d\n", len(parts))

				(*ring).AddNode(logicalHostname, fmt.Sprintf("localhost:%d", port))
			} else {
				logicalHostname := fmt.Sprintf("api-%s-%d.api-%s.default.svc.cluster.local:8080", mode, id, mode)
				physicalHostname := logicalHostname
				(*ring).AddNode(logicalHostname, physicalHostname)
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
