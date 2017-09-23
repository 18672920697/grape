// consistent
package consistent

import (
	"errors"
	"hash/crc32"
	"sort"
	"sync"
)

type uints []uint32

func (x uints) Len() int { return len(x) }

func (x uints) Less(i, j int) bool { return x[i] < x[j] }

func (x uints) Swap(i, j int) { x[i], x[j] = x[j], x[i] }

var ErrEmptyCircle = errors.New("empty circle")

type Consistent struct {
	nodes       map[uint32]string
	sortedNodes uints
	sync.RWMutex
}

func New() *Consistent {
	c := new(Consistent)
	c.nodes = make(map[uint32]string)
	return c
}

func (c *Consistent) AddNode(node string) {
	c.Lock()
	defer c.Unlock()
	c.nodes[c.HashKey(node)] = node
	c.sortCircle()
}

func (c *Consistent) RemoveNode(node string) {
	c.Lock()
	defer c.Unlock()
	delete(c.nodes, c.HashKey(node))
	c.sortCircle()
}

func (c *Consistent) HashKey(key string) uint32 {
	if len(key) < 64 {
		var scratch [64]byte
		copy(scratch[:], key)
		return crc32.ChecksumIEEE(scratch[:len(key)])
	}
	return crc32.ChecksumIEEE([]byte(key))
}

func (c *Consistent) sortCircle() {
	hashes := c.sortedNodes[:0]
	for k := range c.nodes {
		hashes = append(hashes, k)
	}
	sort.Sort(hashes)
	c.sortedNodes = hashes
}

func (c *Consistent) SetKey(key string) (string, error) {
	c.RLock()
	defer c.RUnlock()
	if len(c.nodes) == 0 {
		return "", ErrEmptyCircle
	}
	hashKey := c.HashKey(key)
	i := c.search(hashKey)
	return c.nodes[c.sortedNodes[i]], nil

}

func (c *Consistent) search(key uint32) (i int) {
	f := func(x int) bool {
		return c.sortedNodes[x] > key
	}

	i = sort.Search(len(c.sortedNodes), f)
	if i >= len(c.sortedNodes) {
		i = 0
	}
	return
}
