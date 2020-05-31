package ring

import (
	"fmt"
	"github.com/google/uuid"
	"testing"
)

func TestRing_Get(t *testing.T) {
	r := NewRing("")

	nodes := make(map[string]int)

	for _, n := range r.c.Nodes() {
		nodes[n] = 0
	}

	for i := 0; i < 1000000; i++ {
		item := uuid.New()
		nodes[r.Get(item.String()).String()]++
	}

	avg := 0

	for _, i := range nodes {
		avg += i
	}

	avg /= len(nodes)

	fmt.Printf("avg: %d\n", avg)

	for n, i := range nodes {
		fmt.Printf("%server: %d, %.2f%%\n", n, i, float32(i)*100/float32(avg))
	}

}
