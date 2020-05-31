package ring

import "idfs/oss/monitor/ring/consistent"

const (
	RingSize  = 1 << 32
	VNodeSize = 1 << 24
)

type Ring struct {
	//Nodes map[string]*Node
	c *consistent.Consistent
}

func NewRing(ringFile string) *Ring {
	// ring file
	// host:port:disk
	ringData := []string{
		"127.0.0.1:9001:/tmp/disk1",
		"127.0.0.1:9001:/tmp/disk2",
	}

	r := &Ring{}

	r.c = consistent.NewConsistent()

	for _, n := range ringData {
		r.c.Add(n)
	}

	return r
}

func (r *Ring) Get(item string) *Node {
	var node Node

	nodeS, _ := r.c.Get(item)

	node.Unmarshal(nodeS)
	return &node
}
