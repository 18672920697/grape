package chord

import (
	"testing"
	"time"
)

func init() {

}

/*
	node1 - node4 - node3 - node5 - node2
 */
func TestPredecessorAndSuccessor(t *testing.T) {
	var node1Addr = "127.0.0.1:12601"
	var node2Addr = "127.0.0.1:12502"
	var node3Addr = "127.0.0.1:12603"
	var node4Addr = "127.0.0.1:12604"
	var node5Addr = "127.0.0.1:12605"

	var node1 = Create(node1Addr)
	var node2 = Create(node2Addr)
	var node3 = Create(node3Addr)
	var node4 = Create(node4Addr)
	var node5 = Create(node5Addr)

	sucessor,_ := node1.Join(node2.ipAddr)
	node2.afterJoin(sucessor)

	sucessor,_ = node1.Join(node3.ipAddr)
	node3.afterJoin(sucessor)

	sucessor,_ = node1.Join(node4.ipAddr)
	node4.afterJoin(sucessor)

	time.Sleep(5 * time.Second)

	sucessor,_ = node1.Join(node5.ipAddr)
	node5.afterJoin(sucessor)

	time.Sleep(5 * time.Second)

	node3.stabilize()
	// Check node1
	if node1.successor.ipAddr != node4.ipAddr {
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
	if node3.successor.ipAddr != node5.ipAddr {
		t.Errorf("the predecessor of %s is wrong.", node3.ipAddr)
	}

	if node3.predecessor.ipAddr != node4.ipAddr {
		t.Errorf("the successor of %s is wrong.", node3.ipAddr)
	}

	// Check node4
	if node4.successor.ipAddr != node3.ipAddr {
		t.Errorf("the predecessor of %s is wrong.", node4.ipAddr)
	}

	if node4.predecessor.ipAddr != node1.ipAddr {
		t.Errorf("the successor of %s is wrong.", node4.ipAddr)
	}

	// Check node5
	if node5.successor.ipAddr != node2.ipAddr {
		t.Errorf("the predecessor of %s is wrong.", node5.ipAddr)
	}

	if node5.predecessor.ipAddr != node3.ipAddr {
		t.Errorf("the successor of %s is wrong.", node5.ipAddr)
	}
}
