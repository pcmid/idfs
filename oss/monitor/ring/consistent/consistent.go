package consistent

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"sort"
	"sync"
	"sync/atomic"

	blake2b "github.com/minio/blake2b-simd"
)

const replicationFactor = 10

var ErrNoNodes = errors.New("no nodes added")

type Node struct {
	Name string
	Load int64
}

type Consistent struct {
	nodes     map[uint64]string
	sortedSet []uint64
	loadMap   map[string]*Node
	totalLoad int64

	sync.RWMutex
}

func NewConsistent() *Consistent {
	return &Consistent{
		nodes:     map[uint64]string{},
		sortedSet: []uint64{},
		loadMap:   map[string]*Node{},
	}
}

func (c *Consistent) Add(node string) {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.loadMap[node]; ok {
		return
	}

	c.loadMap[node] = &Node{Name: node, Load: 0}
	for i := 0; i < replicationFactor; i++ {
		h := c.hash(fmt.Sprintf("%s%d", node, i))
		c.nodes[h] = node
		c.sortedSet = append(c.sortedSet, h)

	}
	sort.Slice(c.sortedSet, func(i int, j int) bool {
		if c.sortedSet[i] < c.sortedSet[j] {
			return true
		}
		return false
	})
}

func (c *Consistent) Get(key string) (string, error) {
	c.RLock()
	defer c.RUnlock()

	if len(c.nodes) == 0 {
		return "", ErrNoNodes
	}

	h := c.hash(key)
	idx := c.search(h)
	return c.nodes[c.sortedSet[idx]], nil
}

func (c *Consistent) GetLeast(key string) (string, error) {
	c.RLock()
	defer c.RUnlock()

	if len(c.nodes) == 0 {
		return "", ErrNoNodes
	}

	h := c.hash(key)
	idx := c.search(h)

	i := idx
	for {
		node := c.nodes[c.sortedSet[i]]
		if c.loadOK(node) {
			return node, nil
		}
		i++
		if i >= len(c.nodes) {
			i = 0
		}
	}
}

func (c *Consistent) search(key uint64) int {
	idx := sort.Search(len(c.sortedSet), func(i int) bool {
		return c.sortedSet[i] >= key
	})

	if idx >= len(c.sortedSet) {
		idx = 0
	}
	return idx
}

func (c *Consistent) UpdateLoad(node string, load int64) {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.loadMap[node]; !ok {
		return
	}
	c.totalLoad -= c.loadMap[node].Load
	c.loadMap[node].Load = load
	c.totalLoad += load
}

func (c *Consistent) Inc(node string) {
	c.Lock()
	defer c.Unlock()

	atomic.AddInt64(&c.loadMap[node].Load, 1)
	atomic.AddInt64(&c.totalLoad, 1)
}

func (c *Consistent) Done(node string) {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.loadMap[node]; !ok {
		return
	}
	atomic.AddInt64(&c.loadMap[node].Load, -1)
	atomic.AddInt64(&c.totalLoad, -1)
}

func (c *Consistent) Remove(node string) bool {
	c.Lock()
	defer c.Unlock()

	for i := 0; i < replicationFactor; i++ {
		h := c.hash(fmt.Sprintf("%s%d", node, i))
		delete(c.nodes, h)
		c.delSlice(h)
	}
	delete(c.loadMap, node)
	return true
}

func (c *Consistent) Nodes() (nodes []string) {
	c.RLock()
	defer c.RUnlock()
	for k, _ := range c.loadMap {
		nodes = append(nodes, k)
	}
	return nodes
}

func (c *Consistent) GetLoads() map[string]int64 {
	loads := map[string]int64{}

	for k, v := range c.loadMap {
		loads[k] = v.Load
	}
	return loads
}

func (c *Consistent) MaxLoad() int64 {
	if c.totalLoad == 0 {
		c.totalLoad = 1
	}
	var avgLoadPerNode float64
	avgLoadPerNode = float64(c.totalLoad / int64(len(c.loadMap)))
	if avgLoadPerNode == 0 {
		avgLoadPerNode = 1
	}
	avgLoadPerNode = math.Ceil(avgLoadPerNode * 1.25)
	return int64(avgLoadPerNode)
}

func (c *Consistent) loadOK(node string) bool {
	// a safety check if someone performed c.Done more than needed
	if c.totalLoad < 0 {
		c.totalLoad = 0
	}

	var avgLoadPerNode float64
	avgLoadPerNode = float64((c.totalLoad + 1) / int64(len(c.loadMap)))
	if avgLoadPerNode == 0 {
		avgLoadPerNode = 1
	}
	avgLoadPerNode = math.Ceil(avgLoadPerNode * 1.25)

	bnode, ok := c.loadMap[node]
	if !ok {
		panic(fmt.Sprintf("given node(%server) not in loadsMap", bnode.Name))
	}

	if float64(bnode.Load)+1 <= avgLoadPerNode {
		return true
	}

	return false
}

func (c *Consistent) delSlice(val uint64) {
	idx := -1
	l := 0
	r := len(c.sortedSet) - 1
	for l <= r {
		m := (l + r) / 2
		if c.sortedSet[m] == val {
			idx = m
			break
		} else if c.sortedSet[m] < val {
			l = m + 1
		} else if c.sortedSet[m] > val {
			r = m - 1
		}
	}
	if idx != -1 {
		c.sortedSet = append(c.sortedSet[:idx], c.sortedSet[idx+1:]...)
	}
}

func (c *Consistent) hash(key string) uint64 {
	out := blake2b.Sum512([]byte(key))
	return binary.LittleEndian.Uint64(out[:])
}
