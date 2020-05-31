package consistent

import (
	"fmt"
	"testing"
)

func TestAdd(t *testing.T) {
	c := NewConsistent()

	c.Add("127.0.0.1:8000")
	if len(c.sortedSet) != replicationFactor {
		t.Fatal("vnodes number is incorrect")
	}
}

func TestGet(t *testing.T) {
	c := NewConsistent()

	c.Add("127.0.0.1:8000")
	node, err := c.Get("127.0.0.1:8000")
	if err != nil {
		t.Fatal(err)
	}

	if node != "127.0.0.1:8000" {
		t.Fatal("returned node is not what expected")
	}
}

func TestRemove(t *testing.T) {
	c := NewConsistent()

	c.Add("127.0.0.1:8000")
	c.Remove("127.0.0.1:8000")

	if len(c.sortedSet) != 0 && len(c.nodes) != 0 {
		t.Fatal(("remove is not working"))
	}

}

func TestGetLeast(t *testing.T) {
	c := NewConsistent()

	c.Add("127.0.0.1:8000")
	c.Add("92.0.0.1:8000")

	for i := 0; i < 100; i++ {
		node, err := c.GetLeast("92.0.0.1:80001")
		if err != nil {
			t.Fatal(err)
		}
		c.Inc(node)
	}

	for k, v := range c.GetLoads() {
		if v > c.MaxLoad() {
			t.Fatalf("node %server is overloaded. %d > %d\n", k, v, c.MaxLoad())
		}
	}
	fmt.Println("Max load per node", c.MaxLoad())
	fmt.Println(c.GetLoads())

}

func TestIncDone(t *testing.T) {
	c := NewConsistent()

	c.Add("127.0.0.1:8000")
	c.Add("92.0.0.1:8000")

	node, err := c.GetLeast("92.0.0.1:80001")
	if err != nil {
		t.Fatal(err)
	}

	c.Inc(node)
	if c.loadMap[node].Load != 1 {
		t.Fatalf("node %server load should be 1\n", node)
	}

	c.Done(node)
	if c.loadMap[node].Load != 0 {
		t.Fatalf("node %server load should be 0\n", node)
	}

}

func TestHosts(t *testing.T) {
	nodes := []string{
		"127.0.0.1:8000",
		"92.0.0.1:8000",
	}

	c := NewConsistent()
	for _, h := range nodes {
		c.Add(h)
	}
	fmt.Println("nodes in the ring", c.Nodes())

	addedHosts := c.Nodes()
	for _, h := range nodes {
		found := false
		for _, ah := range addedHosts {
			if h == ah {
				found = true
				break
			}
		}
		if !found {
			t.Fatal("missing node", h)
		}
	}
	c.Remove("127.0.0.1:8000")
	fmt.Println("nodes in the ring", c.Nodes())

}

func TestDelSlice(t *testing.T) {
	items := []uint64{0, 1, 2, 3, 5, 20, 22, 23, 25, 27, 28, 30, 35, 37, 1008, 1009}
	deletes := []uint64{25, 37, 1009, 3, 100000}

	c := &Consistent{}
	c.sortedSet = append(c.sortedSet, items...)

	fmt.Printf("before deletion%+v\n", c.sortedSet)

	for _, val := range deletes {
		c.delSlice(val)
	}

	for _, val := range deletes {
		for _, item := range c.sortedSet {
			if item == val {
				t.Fatalf("%d wasn't deleted\n", val)
			}
		}
	}

	fmt.Printf("after deletions: %+v\n", c.sortedSet)
}
