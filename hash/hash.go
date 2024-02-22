package hash

import (
	"sort"
	"strconv"

	"github.com/spaolacci/murmur3"
)

var weight = 10
var consistentHash *ConsistentHash

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
	node := ch.Nodes[ch.firstKey()]
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

func (ch *ConsistentHash) firstKey() uint32 {
	for k := range ch.Nodes {
		return k
	}
	return 0
}

func GetNode(key string) Node {
	if consistentHash == nil {
		consistentHash = NewConsistentHash()
		consistentHash.AddNode("localhost:8080")
		consistentHash.AddNode("localhost:8081")
		consistentHash.AddNode("localhost:8082")
	}

	return consistentHash.GetNode(key)
}
