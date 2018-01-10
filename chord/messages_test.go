package chord

import (
	"testing"
	"github.com/leviathan1995/grape/proto"
	"github.com/golang/protobuf/proto"
)

func Test_getFingersMessage(t *testing.T) {
	encodeMsg := getFingersMessage()
	decodeMsg := &chordMessages.NetworkMessage{}
	err := proto.Unmarshal(encodeMsg, decodeMsg)
	if err != nil {
		t.Errorf("Unmarshaling getFingersMessage error: %s.", err.Error())
	}
	chordMsg := &chordMessages.ChordMessage{}
	err = proto.Unmarshal([]byte(decodeMsg.GetMsg()), chordMsg)
	if err != nil {
		t.Errorf("Unmarshaling getFingersMessage error: %s.", err.Error())
	}
  if int32(chordMsg.GetCmd()) != chordMessages.ChordMessage_Command_value["GetFingers"] {
		t.Errorf("Construct getFingersMessage command error.%d,%d", int32(chordMsg.GetCmd()), chordMessages.ChordMessage_Command_value["GetPredecessor"])
	}
}

func Test_sendFingersMessage(t *testing.T) {
	var fingers [5]Finger
	fingers[0].ipAddr = "127.0.0.1:1000"
	fingers[1].ipAddr = "127.0.0.1:1001"
	fingers[2].ipAddr = "127.0.0.1:1002"
	fingers[3].ipAddr = "127.0.0.1:1003"
	fingers[4].ipAddr = "127.0.0.1:1004"

	encodeMsg := sendFingersMessage(fingers[:5])
	decodeMsg := &chordMessages.NetworkMessage{}
	err := proto.Unmarshal(encodeMsg, decodeMsg)
	if err != nil {
		t.Errorf("Unmarshaling SendFingersMessage error: %s.", err.Error())
	}
	chordMsg := &chordMessages.ChordMessage{}
	err = proto.Unmarshal([]byte(decodeMsg.GetMsg()), chordMsg)
	if err != nil {
		t.Errorf("Unmarshaling SendFingersMessage error: %s.", err.Error())
	} else if chordMsg.Cmd != 4 {
		t.Errorf("Construct SendFingersMessage command error.")
	}

	// Parse fingers
	fingerMsg := chordMsg.GetSfmsg().GetFingers()
	for index,finger  := range fingerMsg {
		if finger.Address != fingers[index].ipAddr {
			t.Errorf("Parse fingers error: %s-%s", finger.Address, fingers[index].ipAddr)
		}
	}
}

func Test_parseMessage(t *testing.T) {
  var node1Addr = "127.0.0.1:12601"
	var node2Addr = "127.0.0.1:12502"
	var node3Addr = "127.0.0.1:12603"
	var node4Addr = "127.0.0.1:12604"

	var node1 = Create(node1Addr)
	var node2 = Create(node2Addr)
	var node3 = Create(node3Addr)
	var node4 = Create(node4Addr)

	sucessor,_ := node1.Join(node2.ipAddr)
	node2.afterJoin(sucessor)

	sucessor,_ = node1.Join(node3.ipAddr)
	node3.afterJoin(sucessor)

	sucessor,_ = node1.Join(node4.ipAddr)
	node4.afterJoin(sucessor)

  // Ping
  pingMsg := pingMessage()
  parseResp := node1.parseMessage(pingMsg)
  success, err := parsePong(parseResp)
  if err != nil {
    t.Errorf("Parse Ping message error.")
  }

  if !success {
    t.Errorf("pingMessage error.")
  }

  // GetPredecessor
  getPredecessorMsg := getPredecessorMessage()
  parseResp = node1.parseMessage(getPredecessorMsg)
  finger, err := parseFinger(parseResp)
  if err != nil {
    t.Errorf("Parse Predecessor message error.")
  }

  if finger.ipAddr == "" {
    t.Errorf("getPredecessorMessage error.")
  }

  // GetId
  getIdMsg := getIdMessage()
  parseResp = node1.parseMessage(getIdMsg)
  id,err := parseId(parseResp)
  if err != nil {
    t.Errorf("Parse Id message error.")
  }

  if string(id[:]) == "" {
    t.Errorf("getIdMessage error.")
  }

  // GetSuccessors
  getSuccessorMsg := getSuccessorsMessage()
  parseResp = node1.parseMessage(getSuccessorMsg)
  fingers, err := parseFingers(parseResp)
  if err != nil {
    t.Errorf("Parse Successors error.")
  }
  for _, finger := range fingers {
    if finger.ipAddr == "" {
      t.Errorf("getSuccessorsMessage error.")
    }
  }
}
