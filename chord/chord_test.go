package chord

import (
	"testing"
	"time"
)

func init() {

}

/*
	node2 - node1 - node3
*/
func TestPredecessorAndSuccessor(t *testing.T) {
	var node1Addr = "127.0.0.1:16001"
	var node2Addr = "127.0.0.1:16002"
	var node3Addr = "127.0.0.1:16003"

	var node1 = Create(node1Addr)
	var node2 = Create(node2Addr)
	var node3 = Create(node3Addr)

	sucessor, _ := node1.Join(node2.ipAddr)
	node2.afterJoin(sucessor)

	sucessor, _ = node1.Join(node3.ipAddr)
	node3.afterJoin(sucessor)

	time.Sleep(3 * time.Second)
	// Check node1
	if node1.successor.ipAddr != node3.ipAddr {
		t.Errorf("the predecessor of %s is wrong.", node1.ipAddr)
	}

	if node1.predecessor.ipAddr != node2.ipAddr {
		t.Errorf("the successor of %s is wrong.", node1.ipAddr)
	}

	// Check node2
	if node2.successor.ipAddr != node1.ipAddr {
		t.Errorf("the predecessor of %s is wrong.", node2.ipAddr)
	}

	if node2.predecessor.ipAddr != node3.ipAddr {
		t.Errorf("the successor of %s is wrong.", node2.ipAddr)
	}

	// Check node3
	if node3.successor.ipAddr != node2.ipAddr {
		t.Errorf("the predecessor of %s is wrong.", node3.ipAddr)
	}

	if node3.predecessor.ipAddr != node1.ipAddr {
		t.Errorf("the successor of %s is wrong.", node3.ipAddr)
	}
}
