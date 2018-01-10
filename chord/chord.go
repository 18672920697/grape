package chord

import (
	"crypto/sha256"
	"fmt"
	"github.com/leviathan1995/grape/logger"
	"math/big"
	"net"
	"time"
)

//	Finger type denoting identifying information about a ChordNode
type Finger struct {
	id     [sha256.Size]byte
	ipAddr string
}

type request struct {
	write bool
	successor  bool
	index int
}

//	ChordNode type denoting a Chord server.
type ChordNode struct {
	predecessor   *Finger
	successor     *Finger
	successorList [sha256.Size * 8]Finger
	fingerTable   [sha256.Size*8 + 1]Finger

	finger       chan Finger
	request      chan request

	id     [sha256.Size]byte
	ipAddr string

	connections  map[string]net.TCPConn
}

type PeerError struct {
	Address string
	Err     error
}

func (e *PeerError) Error() string {
	return fmt.Sprintf("Failed to connect to peer: %s. Cause of failure: %s.", e.Address, e.Err)
}

//	Lookup returns the address of the successor of key in the Chord DHT.
//	The lookup process is iterative. Beginning with the address of a
//	Chord node, start, this function will request the finger tables of
//	the closest preceeding Chord node to key until the successor is found.
//
//	If the start address is unreachable, the error is of type PeerError.
func Lookup(key [sha256.Size]byte, start string) (addr string, err error) {

	addr = start

	msg := getFingersMessage()
	reply, err := Send(msg, start)
	if err != nil {
		err = &PeerError{start, err}
		return
	}

	ft, err := parseFingers(reply)
	if err != nil {
		err = &PeerError{start, err}
		return
	}
	if len(ft) < 2 {
		return
	}

	current := ft[0]

	if key == current.id {
		addr = current.ipAddr
		return
	}

	// Loop through finger table and see what the closest finger is
	for i := len(ft) - 1; i > 0; i-- {
		f := ft[i]
		if i == 0 {
			break
		}
		if InRange(f.id, current.id, key) { // See if f.id is closer than I am.
			addr, err = Lookup(key, f.ipAddr)
			if err != nil {
				continue
			}
			return
		}
	}
	addr = ft[1].ipAddr
	msg = pingMessage()
	reply, err = Send(msg, addr)

	// If the current node's successor has gone missing
	if err != nil {
		// Ask node for its successor list
		msg = getSuccessorsMessage()
		reply, err = Send(msg, current.ipAddr)
		if err != nil {
			addr = current.ipAddr
			return
		}

		ft, err = parseFingers(reply)
		if err != nil {
			addr = current.ipAddr
			return
		}

		for i := 0; i < len(ft); i++ {
			f := ft[i]
			if i == 0 {
				break
			}
			msg = pingMessage()
			reply, err = Send(msg, f.ipAddr)
			if err != nil { // Closest next successor that responds
				addr = f.ipAddr
				return
			}
		}

		addr = current.ipAddr
		return
	}

	return
}

//	Lookup returns the address of the ChordNode that is responsible
//	for the key. The procedure begins at the address denoted by start.
func (node *ChordNode) lookup(key [sha256.Size]byte, start string) (addr string, err error) {

	addr = start

	msg := getFingersMessage()
	reply, err := node.send(msg, start)
	if err != nil {
		err = &PeerError{start, err}
		return
	}

	ft, err := parseFingers(reply)
	if err != nil {
		err = &PeerError{start, err}
		return
	}
	if len(ft) < 2 {
		return
	}

	current := ft[0]

	if key == current.id {
		addr = current.ipAddr
		return
	}

	//loop through finger table and see what the closest finger is
	for i := len(ft) - 1; i > 0; i-- {
		f := ft[i]
		if i == 0 {
			break
		}
		if InRange(f.id, current.id, key) { //see if f.id is closer than I am.
			addr, err = node.lookup(key, f.ipAddr)
			if err != nil { //node failed
				continue
			}
			return
		}
	}
	addr = ft[1].ipAddr
	msg = pingMessage()
	reply, err = node.send(msg, addr)

	//this code is executed if the id's successor has gone missing
	if err != nil {
		//ask node for its successor list
		msg = getSuccessorsMessage()
		reply, err = node.send(msg, current.ipAddr)
		if err != nil {
			addr = current.ipAddr
			return
		}
		ft, err = parseFingers(reply)
		if err != nil {
			addr = current.ipAddr
			return
		}

		for i := 0; i < len(ft); i++ {
			f := ft[i]
			if i == 0 {
				break
			}
			msg = pingMessage()
			reply, err = node.send(msg, f.ipAddr)
			if err != nil { //closest next successor that responds
				addr = f.ipAddr
				return
			}
		}

		addr = current.ipAddr
		return
	}

	return
}

//	Create will start a new Chord DHT and return the original ChordNode
func Create(myaddr string) *ChordNode {
	node := new(ChordNode)
	// Initialize node information
	node.id = sha256.Sum256([]byte(myaddr))
	node.ipAddr = myaddr
	me := new(Finger)
	me.id = node.id
	me.ipAddr = node.ipAddr
	node.fingerTable[0] = *me
	succ := new(Finger)
	node.successor = succ
	pred := new(Finger)
	node.predecessor = pred

	//	Set up channels for finger manager
	c := make(chan Finger)
	c2 := make(chan request)
	node.finger = c
	node.request = c2

	//	Initialize listener and network manager threads
	node.listen(myaddr)
	node.connections = make(map[string]net.TCPConn)

	// Stabilize
	node.stabilize()
	// Check predecessor
	node.checkPredecessor()
	// Update fingers
	node.fixFinger()

	//	Initialize maintenance and finger manager threads
	go node.infoRouteTable()
	go node.maintainRouteTable()

	return node
}

//	Join will add a new ChordNode to an existing DHT. It looks up the successor
//	of the new node starting at an existing Chord node specified by addr.
//
//	If the start address is unreachable, the error is of type PeerError.
func (node *ChordNode) Join(peerAddr string) (*Finger, error) {
	newId := sha256.Sum256([]byte(peerAddr))
	successor, err := Lookup(newId, node.ipAddr)
	if err != nil || successor == "" {
		return nil, &PeerError{peerAddr, err}
	}

	//	find the id of successor node
	msg := getIdMessage()
	reply, err := Send(msg, successor)
	if err != nil {
		return nil, &PeerError{peerAddr, err}
	}

	//	update node info to include successor
	succ := new(Finger)
	succ.id, err = parseId(reply)
	if err != nil {
		return nil, &PeerError{peerAddr, err}
	}
	succ.ipAddr = successor

	return succ, nil
}

func (node *ChordNode) afterJoin(successor *Finger) {
	node.successor = successor
	node.fingerTable[1] = *successor
	node.successorList[0] = *successor
	node.stabilize()
	node.fixFinger()
}
//	Manages reads and writes to the route table of node
func (node *ChordNode) infoRouteTable() {
	for {
		req := <-node.request
		if req.write {	// write
			if req.successor {
				node.successorList[req.index] = <-node.finger
			} else {
				if req.index < 0 {
					*node.predecessor = <-node.finger
				} else if req.index == 1 {
					*node.successor = <-node.finger
					node.fingerTable[1] = *node.successor
					node.successorList[0] = *node.successor
				} else {
					node.fingerTable[req.index] = <-node.finger
				}
			}
		} else {	// read
			if req.successor {
				node.finger <- node.successorList[req.index]
			} else {
				if req.index < 0 {
					node.finger <- *node.predecessor
				} else {
					node.finger <- node.fingerTable[req.index]
				}
			}
		}
	}
}

//query allows functions to read from or write to the node object
func (node *ChordNode) query(write bool, succ bool, index int, newf *Finger) Finger {
	f := new(Finger)
	req := request{write, succ, index}
	node.request <- req
	if write {
		node.finger <- *newf
	} else {
		*f = <-node.finger
	}

	return *f
}

//	maintain will periodically perform maintenance operations
func (node *ChordNode) maintainRouteTable() {
	for {
		time.Sleep(1 * time.Duration(time.Millisecond))

		node.stabilize()

		node.checkPredecessor()

		node.fixFinger()
	}
}

// 1. Check to see if the successor is still around, if not, find the next available successor
// 2. Update successor list
// 3. Stablize ensures that the node's successor's predecessor is itself
// 	  If not, it updates its successor's predecessor.
func (node *ChordNode) stabilize() {
	successor := node.fingerTable[1]
	if successor.zero() {
		return
	}

	// Check to see if the successor is still around
	msg := pingMessage()
	reply, err := node.send(msg, successor.ipAddr)
	if err != nil {
		// If successor failed to respond
		// Check in successor list for next available successor.
		for i := 1; i < sha256.Size*8; i++ {
			//successor = node.query(false, true, i, nil)
			successor = node.successorList[i]
			if successor.ipAddr == node.ipAddr {
				continue
			}
			msg := pingMessage()
			reply, err = node.send(msg, successor.ipAddr)
			if err == nil {
				break
			} else {
				successor.ipAddr = ""
			}
		}
		//node.query(true, false, 1, &successor)
		node.successor = &successor
		node.fingerTable[1] = successor
		node.successorList[0] = successor
		if successor.ipAddr == "" {
			return
		}
	}

	// Update successor list
	msg = getSuccessorsMessage()
	reply, err = node.send(msg, successor.ipAddr)
	if err != nil {
		return
	}
	ft, err := parseFingers(reply)
	if err != nil {
		return
	}
	for i := range ft {
		if i < sha256.Size*8-1 {
			//node.query(true, true, i+1, &ft[i])
			node.successorList[i+1] = ft[i]
		}
	}

	//	Ask successor for predecessor
	msg = getPredecessorMessage()
	reply, err = node.send(msg, successor.ipAddr)
	if err != nil {
		return
	}

	predecessorOfSuccessor, err := parseFinger(reply)
	if err != nil { // Node failed
		return
	}
	if predecessorOfSuccessor.ipAddr != "" {
		if predecessorOfSuccessor.id != node.id {
			if InRange(predecessorOfSuccessor.id, node.id, successor.id) {
				node.successor = &predecessorOfSuccessor
				node.fingerTable[1] = predecessorOfSuccessor
				node.successorList[0] = predecessorOfSuccessor
			}
		} else { // Everything is fine
			return
		}
	}

	// Claim to be predecessor of successor
	me := new(Finger)
	me.id = node.id
	me.ipAddr = node.ipAddr
	msg = claimPredecessorMessage(*me)
	node.send(msg, successor.ipAddr)
}

func (node *ChordNode) notify(newPred Finger) {
	//node.query(true, false, -1, &newPred)
	node.predecessor = &newPred
	//update predecessor
	//successor := node.query(false, false, 1, nil)
	successor := node.fingerTable[1]
	if successor.zero() { //TODO: so if you get here, you were probably the first node.
		//node.query(true, false, 1, &newPred)
		node.successor = &newPred
		node.fingerTable[1] = newPred
		node.successorList[0] = newPred
	}
	//node.stabilize()
	//node.fixFinger()
}

func (node *ChordNode) checkPredecessor() {
	//predecessor := node.query(false, false, -1, nil)
	predecessor := *node.predecessor
	if predecessor.zero() {
		return
	}

	msg := pingMessage()
	reply, err := node.send(msg, predecessor.ipAddr)
	if err != nil {
		predecessor.ipAddr = ""
		//node.query(true, false, -1, &predecessor)
		node.predecessor = &predecessor
	} else {
		return
	}

	if success, err := parsePong(reply); !success || err != nil {
		predecessor.ipAddr = ""
		//node.query(true, false, -1, &predecessor)
		node.predecessor = &predecessor
	}

	return
}

//
func (node *ChordNode) fixFinger() {
	for which := 0; which < 257; which++ {
		successor := node.fingerTable[1]
		if which == 0 || which == 1 || successor.zero() {
			continue
		}
		var targetId [sha256.Size]byte
		copy(targetId[:sha256.Size], target(node.id, which)[:sha256.Size])
		newip, err := node.lookup(targetId, successor.ipAddr)
		if err != nil {
			logger.Error.Printf("%s\n", err.Error())
			continue
		}
		if newip == node.ipAddr {
			if err != nil {
				logger.Error.Printf("%s\n", err.Error())
			}
			continue
		}

		// Find the id of node
		msg := getIdMessage()
		reply, err := node.send(msg, newip)
		if err != nil {
			logger.Error.Printf("%s\n", err.Error())
			continue
		}

		newfinger := new(Finger)
		newfinger.ipAddr = newip
		newfinger.id, _ = parseId(reply)
		if which == 1 {
			node.successor = newfinger
			node.fingerTable[1] = *newfinger
			node.successorList[0] = *newfinger
		} else {
			node.fingerTable[which] = *newfinger
		}
	}
}

// Finalize stops all communication and removes the ChordNode from the DHT.
func (node *ChordNode) Finalize() {
	//send message to all children to terminate

	fmt.Printf("Exiting...\n")
}

// InRange is a helper function that returns true if the value x is between the values (min, max)
func InRange(x [sha256.Size]byte, min [sha256.Size]byte, max [sha256.Size]byte) bool {
	//There are 3 cases: min < x and x < max,
	//x < max and max < min, max < min and min < x
	xint := new(big.Int)
	maxint := new(big.Int)
	minint := new(big.Int)
	xint.SetBytes(x[:sha256.Size])
	minint.SetBytes(min[:sha256.Size])
	maxint.SetBytes(max[:sha256.Size])

	if xint.Cmp(minint) == 1 && maxint.Cmp(xint) == 1 {
		return true
	}

	if maxint.Cmp(xint) == 1 && minint.Cmp(maxint) == 1 {
		return true
	}

	if minint.Cmp(maxint) == 1 && xint.Cmp(minint) == 1 {
		return true
	}

	return false
}

// target() returns the target id used by the fix function
func target(me [sha256.Size]byte, which int) []byte {
	meint := new(big.Int)
	meint.SetBytes(me[:sha256.Size])

	baseint := new(big.Int)
	baseint.SetUint64(2)

	powint := new(big.Int)
	powint.SetInt64(int64(which - 1))

	var biggest [sha256.Size + 1]byte
	for i := range biggest {
		biggest[i] = 255
	}

	tmp := new(big.Int)
	tmp.SetInt64(1)

	modint := new(big.Int)
	modint.SetBytes(biggest[:sha256.Size])
	modint.Add(modint, tmp)

	target := new(big.Int)
	target.Exp(baseint, powint, modint)
	target.Add(meint, target)
	target.Mod(target, modint)

	bytes := target.Bytes()
	diff := sha256.Size - len(bytes)
	if diff > 0 {
		tmp := make([]byte, sha256.Size)
		//pad with zeros
		for i := 0; i < diff; i++ {
			tmp[i] = 0
		}
		for i := diff; i < sha256.Size; i++ {
			tmp[i] = bytes[i-diff]
		}
		bytes = tmp
	}
	return bytes[:sha256.Size]
}

func (f Finger) String() string {
	return fmt.Sprintf("%s", f.ipAddr)
}

func (f Finger) zero() bool {
	if f.ipAddr == "" {
		return true
	} else {
		return false
	}
}

/** Printouts of information **/

//String returns a string containing the node's ip address, sucessor, and predecessor.
func (node *ChordNode) String() string {
	var succ, pred string
	successor := node.query(false, false, 1, nil)
	predecessor := node.query(false, false, -1, nil)
	if !successor.zero() {
		succ = successor.String()
	} else {
		succ = "Unknown"
	}
	if !predecessor.zero() {
		pred = predecessor.String()
	} else {
		pred = "Unknown"
	}
	return fmt.Sprintf("%s\t%s\t%s\n", node.ipAddr, succ, pred)
}

//ShowFingers returns a string representation of the ChordNode's finger table.
func (node *ChordNode) ShowFingers() string {
	retval := ""
	finger := new(Finger)
	prevfinger := new(Finger)
	ctr := 0
	for i := 0; i < sha256.Size*8+1; i++ {
		*finger = node.query(false, false, i, nil)
		if !finger.zero() {
			ctr += 1
			if i == 0 || finger.ipAddr != prevfinger.ipAddr {
				retval += fmt.Sprintf("%d %s\n", i, finger.String())
			}
		}
		*prevfinger = *finger
	}
	return retval + fmt.Sprintf("Total fingers: %d.\n", ctr)
}

// ShowSucc returns a string representation of the ChordNode's successor list.
func (node *ChordNode) ShowSucc() string {
	table := ""
	finger := new(Finger)
	prevfinger := new(Finger)
	for i := 0; i < sha256.Size*8; i++ {
		*finger = node.query(false, true, i, nil)
		if finger.ipAddr != "" {
			if i == 0 || finger.ipAddr != prevfinger.ipAddr {
				table += fmt.Sprintf("%s\n", finger.String())
			}
		}
		*prevfinger = *finger
	}
	return table
}
