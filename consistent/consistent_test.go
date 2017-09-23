package consistent

import (
	"sort"
	"strconv"
	"testing"
)

func checkNum(num, expected int, t *testing.T) {
	if num != expected {
		t.Errorf("get %d, expected %d", num, expected)
	}
}

func TestNew(t *testing.T) {
	x := New()
	if x == nil {
		t.Errorf("expected object")
	}
}

func TestAddNode(t *testing.T) {
	x := New()
	x.AddNode("node1")
	x.AddNode("node2")
	x.AddNode("node3")
	x.AddNode("node4")
	x.AddNode("node5")

	checkNum(len(x.nodes), 5, t)
	checkNum(len(x.sortedNodes), 5, t)
	if sort.IsSorted(x.sortedNodes) == false {
		t.Errorf("expected sorted hashes to be sorted")
	}
}

func TestRemoveNode(t *testing.T) {
	x := New()
	x.AddNode("node1")
	checkNum(len(x.nodes), 1, t)
	x.RemoveNode("node1")
	checkNum(len(x.nodes), 0, t)
}

func TestSetKey(t *testing.T) {
	x := New()
	x.AddNode("node1")
	x.AddNode("node2")
	x.AddNode("node3")
	x.AddNode("node4")
	x.AddNode("node5")

	for i := 0; i < 10; i++ {
		server, err := x.SetKey("key" + strconv.Itoa(i))
		if _, ok := x.nodes[x.HashKey(server)]; !ok {
			t.Error(err)
		}
	}
}
