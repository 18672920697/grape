package chord

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/leviathan1995/grape/logger"
	"github.com/leviathan1995/grape/proto"
	"log"
)

// Construct a message to request the fingers
func getFingersMessage() []byte {
	msg := new(chordMessages.NetworkMessage)
	msg.Proto = *proto.Uint32(1)

	chordMessage := new(chordMessages.ChordMessage)
	command := chordMessages.ChordMessage_Command(chordMessages.ChordMessage_Command_value["GetFingers"])
	chordMessage.Cmd = command
	chorddata, err := proto.Marshal(chordMessage)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}

	msg.Msg = *proto.String(string(chorddata))
	data, err := proto.Marshal(msg)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}

	return data
}

// Send the fingers
func sendFingersMessage(fingers []Finger) []byte {
	msg := new(chordMessages.NetworkMessage)
	msg.Proto = *proto.Uint32(1)

	chordMessage := new(chordMessages.ChordMessage)
	command := chordMessages.ChordMessage_Command(chordMessages.ChordMessage_Command_value["GetFingers"])
	chordMessage.Cmd = command

	sendFingersMessage := new(chordMessages.SendFingersMessage)
	for _, finger := range fingers {
		if !finger.zero() {
			fingerMsg := new(chordMessages.FingerMessage)
			fingerMsg.Id = *proto.String(string(finger.id[:32]))
			fingerMsg.Address = *proto.String(finger.ipAddr)
			sendFingersMessage.Fingers = append(sendFingersMessage.Fingers, fingerMsg)
		}
	}
	chordMessage.Sfmsg = sendFingersMessage

	chorddata, err := proto.Marshal(chordMessage)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}
	msg.Msg = *proto.String(string(chorddata))

	data, err := proto.Marshal(msg)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}

	return data
}

// getIdMessage constructs a message to ask a server for its chord id
func getIdMessage() []byte {

	msg := new(chordMessages.NetworkMessage)
	msg.Proto = *proto.Uint32(1)

	chordMessage := new(chordMessages.ChordMessage)
	command := chordMessages.ChordMessage_Command(chordMessages.ChordMessage_Command_value["GetId"])
	chordMessage.Cmd = command
	chorddata, err := proto.Marshal(chordMessage)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}
	msg.Msg = *proto.String(string(chorddata))

	data, err := proto.Marshal(msg)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}

	return data
}

// sendIdMessage constructs a message to ask a server for its chord id
func sendIdMessage(id []byte) []byte {

	msg := new(chordMessages.NetworkMessage)
	msg.Proto = *proto.Uint32(1)
	chordMsg := new(chordMessages.ChordMessage)
	command := chordMessages.ChordMessage_Command(chordMessages.ChordMessage_Command_value["GetId"])
	chordMsg.Cmd = command
	sidMsg := new(chordMessages.SendIdMessage)
	sidMsg.Id = *proto.String(string(id))
	chordMsg.Sidmsg = sidMsg
	chorddata, err := proto.Marshal(chordMsg)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}
	msg.Msg = *proto.String(string(chorddata))

	data, err := proto.Marshal(msg)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}

	return data
}

func getPredecessorMessage() []byte {
	msg := new(chordMessages.NetworkMessage)
	msg.Proto = *proto.Uint32(1)
	chordMsg := new(chordMessages.ChordMessage)
	command := chordMessages.ChordMessage_Command(chordMessages.ChordMessage_Command_value["GetPredecessor"])
	chordMsg.Cmd = command
	chorddata, err := proto.Marshal(chordMsg)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}
	msg.Msg = *proto.String(string(chorddata))

	data, err := proto.Marshal(msg)

	if err != nil {
		log.Fatal("marshaling error: ", err)
	}

	return data
}

func sendPredecessorMessage(finger Finger) []byte {
	msg := new(chordMessages.NetworkMessage)
	msg.Proto = *proto.Uint32(1)
	chordMsg := new(chordMessages.ChordMessage)
	command := chordMessages.ChordMessage_Command(chordMessages.ChordMessage_Command_value["GetPredecessor"])
	chordMsg.Cmd = command
	pMsg := new(chordMessages.PredecessorMessage)
	fingerMsg := new(chordMessages.FingerMessage)
	fingerMsg.Id = *proto.String(string(finger.id[:32]))
	fingerMsg.Address = *proto.String(finger.ipAddr)
	pMsg.Predecessor = fingerMsg
	chordMsg.Cpmsg = pMsg

	chorddata, err := proto.Marshal(chordMsg)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}
	msg.Msg = *proto.String(string(chorddata))

	data, err := proto.Marshal(msg)

	if err != nil {
		log.Fatal("marshaling error: ", err)
	}

	return data
}

func claimPredecessorMessage(finger Finger) []byte {
	msg := new(chordMessages.NetworkMessage)
	msg.Proto = *proto.Uint32(1)
	chordMsg := new(chordMessages.ChordMessage)
	command := chordMessages.ChordMessage_Command(chordMessages.ChordMessage_Command_value["ClaimPredecessor"])
	chordMsg.Cmd = command
	predMsg := new(chordMessages.PredecessorMessage)
	fingerMsg := new(chordMessages.FingerMessage)
	fingerMsg.Id = *proto.String(string(finger.id[:32]))
	fingerMsg.Address = *proto.String(finger.ipAddr)
	predMsg.Predecessor = fingerMsg
	chordMsg.Cpmsg = predMsg

	chorddata, err := proto.Marshal(chordMsg)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}
	msg.Msg = *proto.String(string(chorddata))

	data, err := proto.Marshal(msg)

	if err != nil {
		log.Fatal("marshaling error: ", err)
	}

	return data
}

func pingMessage() []byte {

	msg := new(chordMessages.NetworkMessage)
	msg.Proto = *proto.Uint32(1)
	chordMsg := new(chordMessages.ChordMessage)
	command := chordMessages.ChordMessage_Command(chordMessages.ChordMessage_Command_value["Ping"])
	chordMsg.Cmd = command
	chorddata, err := proto.Marshal(chordMsg)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}
	msg.Msg = *proto.String(string(chorddata))

	data, err := proto.Marshal(msg)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}

	return data
}

func pongMessage() []byte {

	msg := new(chordMessages.NetworkMessage)
	msg.Proto = *proto.Uint32(1)
	chordMsg := new(chordMessages.ChordMessage)
	command := chordMessages.ChordMessage_Command(chordMessages.ChordMessage_Command_value["Pong"])
	chordMsg.Cmd = command
	chorddata, err := proto.Marshal(chordMsg)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}
	msg.Msg = *proto.String(string(chorddata))

	data, err := proto.Marshal(msg)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}

	return data
}

func getSuccessorsMessage() []byte {
	msg := new(chordMessages.NetworkMessage)
	msg.Proto = *proto.Uint32(1)
	chordMsg := new(chordMessages.ChordMessage)
	command := chordMessages.ChordMessage_Command(chordMessages.ChordMessage_Command_value["GetSuccessor"])
	chordMsg.Cmd = command
	chorddata, err := proto.Marshal(chordMsg)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}
	msg.Msg = *proto.String(string(chorddata))

	data, err := proto.Marshal(msg)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}

	return data

}

func nullMessage() []byte {
	msg := new(chordMessages.NetworkMessage)
	msg.Proto = *proto.Uint32(1)

	data, err := proto.Marshal(msg)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}

	return data
}

// parseMessage takes as input an unmarshalled message and
// performs actions based on what the message contains.
func (node *ChordNode) parseMessage(data []byte) []byte {
	msg := new(chordMessages.NetworkMessage)

	err := proto.Unmarshal(data, msg)
	if err != nil {
		logger.Error.Printf("%s,%s ", err.Error(), node.ipAddr)
		return nil
	}

	chorddata := []byte(msg.GetMsg())
	chordmsg := new(chordMessages.ChordMessage)
	err = proto.Unmarshal(chorddata, chordmsg)
	if err != nil {
		logger.Error.Printf("%s,%s ", err.Error(), node.ipAddr)
		return nil
	}

	cmd := int32(chordmsg.GetCmd())
	switch {
	case cmd == chordMessages.ChordMessage_Command_value["Ping"]:
	  return pongMessage()
	case cmd == chordMessages.ChordMessage_Command_value["GetPredecessor"]:
		predecessor := *node.predecessor
		if predecessor.zero() {
			return nullMessage()
		} else {
		  return sendPredecessorMessage(predecessor)
		}
	case cmd == chordMessages.ChordMessage_Command_value["GetId"]:
		return sendIdMessage(node.id[:32])
	case cmd == chordMessages.ChordMessage_Command_value["GetFingers"]:
		table := make([]Finger, 32*8+1)
		for i := range table {
			table[i] = node.fingerTable[i]
		}

		return sendFingersMessage(table)
	case cmd == chordMessages.ChordMessage_Command_value["ClaimPredecessor"]:
		newPredecessor, err := parseFinger(data)
		if err != nil {
			logger.Error.Printf("%s", err.Error())
			return nil
		}
		if err != nil {
			return nullMessage()
			break
		}
		predecessor := node.predecessor

		if predecessor.zero() || InRange(newPredecessor.id, predecessor.id, node.id) {
			node.notify(newPredecessor)
		}
		return nullMessage()
	case cmd == chordMessages.ChordMessage_Command_value["GetSuccessor"]:
		table := make([]Finger, 32*8)
		for i := range table {
			table[i] = node.fingerTable[i]
		}
		return sendFingersMessage(table)
	}
	fmt.Printf("No matching commands.\n")
  return nil
}

//parseFingers can be called to return a finger table from a received
//message after a getfingers call.
func parseFingers(data []byte) (ft []Finger, err error) {
	msg := new(chordMessages.NetworkMessage)
	err = proto.Unmarshal(data, msg)
	if msg.GetProto() != 1 {
		return
	}
	chorddata := []byte(msg.GetMsg())
	chordmsg := new(chordMessages.ChordMessage)
	err = proto.Unmarshal(chorddata, chordmsg)
	if err != nil {
		logger.Error.Printf("%s", err.Error())
		return
	}
	if chordmsg == nil {
		return
	}
	sfmsg := chordmsg.GetSfmsg()
	fingers := sfmsg.GetFingers()
	prevfinger := new(Finger)
	for _, finger := range fingers {
		newfinger := new(Finger)
		copy(newfinger.id[:], []byte(finger.Id))
		newfinger.ipAddr = finger.Address
		if !newfinger.zero() && newfinger.ipAddr != prevfinger.ipAddr {
			ft = append(ft, *newfinger)
		}
		*prevfinger = *newfinger
	}
	return
}

func parseFinger(data []byte) (f Finger, err error) {
	msg := new(chordMessages.NetworkMessage)
	err = proto.Unmarshal(data, msg)
	if err != nil {
		logger.Error.Printf("%s", err.Error())
		return
	}

	if msg.GetMsg() == "" {
		return
	}

	chorddata := []byte(msg.GetMsg())
	chordmsg := new(chordMessages.ChordMessage)
	err = proto.Unmarshal(chorddata, chordmsg)
	if err != nil {
		logger.Error.Printf("%s", err.Error())
		return
	}

	cpmsg := chordmsg.GetCpmsg()
	finger := cpmsg.GetPredecessor()
	if finger == nil {
		return
	}
	copy(f.id[:], []byte(finger.Id))
	f.ipAddr = finger.Address

	return
}

func parseId(data []byte) (id [32]byte, err error) {
	msg := new(chordMessages.NetworkMessage)
	err = proto.Unmarshal(data, msg)
	if msg.GetProto() != 1 {
		return
	}

	chorddata := []byte(msg.GetMsg())
	chordmsg := new(chordMessages.ChordMessage)
	err = proto.Unmarshal(chorddata, chordmsg)
	if err != nil {
		logger.Error.Printf("%s", err.Error())
		return
	}

	if chordmsg == nil {
		return
	}

	idmsg := chordmsg.GetSidmsg()
	arr := []byte(idmsg.GetId())
	copy(id[:], arr[:32])
	return
}

func parsePong(data []byte) (success bool, err error) {

	msg := new(chordMessages.NetworkMessage)
	err = proto.Unmarshal(data, msg)
	if err != nil {
		logger.Error.Printf("%s", err.Error())
		return false, err
	}

	if msg.GetProto() != 1 {
		return
	}

	chorddata := []byte(msg.GetMsg())
	chordmsg := new(chordMessages.ChordMessage)
	err = proto.Unmarshal(chorddata, chordmsg)
	if err != nil {
		logger.Error.Printf("%s", err.Error())
		return
	}

	if chordmsg == nil {
		return
	}

	command := int32(chordmsg.GetCmd())
	if command == chordMessages.ChordMessage_Command_value["Pong"] {
		success = true
	} else {
		success = false
	}

	return
}
