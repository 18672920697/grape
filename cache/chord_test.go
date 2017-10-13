package cache

import (
	"testing"
)

func TestChord(t *testing.T) {
	node1Addr := "127.0.0.1:10001"
	node2Addr := "127.0.0.1:10002"
	node1 := Create(node1Addr)
	node2 := Create(node2Addr)

	node1.Join(node1.ipaddr, node2.ipaddr)

	if node1.predecessor.ipaddr != node2Addr {
		t.Errorf("node2 finger table error")
	}
	for _, f := range node1.fingerTable[1:] {
		if f.ipaddr != node2Addr {
			t.Errorf("node1 finger table error")
		}
	}

	if node1.predecessor.ipaddr != node2Addr {
		t.Errorf("node2 finger table error")
	}
	for _, f := range node2.fingerTable[1:] {
		if f.ipaddr != node1Addr {
			t.Errorf("node2 finger table error")
		}
	}

}
