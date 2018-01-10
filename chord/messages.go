package cache

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/leviathan1995/grape/logger"
	"github.com/leviathan1995/grape/proto"
	"log"
)

//lookupMsg constructs a message to perform the lookup of a key and returns the
//marshalled redis buffer
func getfingersMsg() []byte {

	msg := new(chordMsgs.NetworkMessage)
	msg.Proto = *proto.Uint32(1)
	chordMsg := new(chordMsgs.ChordMessage)
	command := chordMsgs.ChordMessage_Command(chordMsgs.ChordMessage_Command_value["GetFingers"])
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

func sendfingersMsg(fingers []Finger) []byte {

	msg := new(chordMsgs.NetworkMessage)
	msg.Proto = *proto.Uint32(1)
	chordMsg := new(chordMsgs.ChordMessage)
	command := chordMsgs.ChordMessage_Command(chordMsgs.ChordMessage_Command_value["GetFingers"])
	chordMsg.Cmd = command
	sfMsg := new(chordMsgs.SendFingersMessage)
	for _, finger := range fingers {
		if !finger.zero() {
			fingerMsg := new(chordMsgs.FingerMessage)
			fingerMsg.Id = *proto.String(string(finger.id[:32]))
			fingerMsg.Address = *proto.String(finger.ipAddr)
			sfMsg.Fingers = append(sfMsg.Fingers, fingerMsg)
		}
	}
	chordMsg.Sfmsg = sfMsg
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

//getidMsg constructs a message to ask a server for its chord id
func getidMsg() []byte {

	msg := new(chordMsgs.NetworkMessage)
	msg.Proto = *proto.Uint32(1)
	chordMsg := new(chordMsgs.ChordMessage)
	command := chordMsgs.ChordMessage_Command(chordMsgs.ChordMessage_Command_value["GetId"])
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

//sendidMsg constructs a message to ask a server for its chord id
func sendidMsg(id []byte) []byte {

	msg := new(chordMsgs.NetworkMessage)
	msg.Proto = *proto.Uint32(1)
	chordMsg := new(chordMsgs.ChordMessage)
	command := chordMsgs.ChordMessage_Command(chordMsgs.ChordMessage_Command_value["GetId"])
	chordMsg.Cmd = command
	sidMsg := new(chordMsgs.SendIdMessage)
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

//TODO: rewrite
func getpredMsg() []byte {
	msg := new(chordMsgs.NetworkMessage)
	msg.Proto = *proto.Uint32(1)
	chordMsg := new(chordMsgs.ChordMessage)
	command := chordMsgs.ChordMessage_Command(chordMsgs.ChordMessage_Command_value["GetPred"])
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

//TODO: rewrite
func sendpredMsg(finger Finger) []byte {
	msg := new(chordMsgs.NetworkMessage)
	msg.Proto = *proto.Uint32(1)
	chordMsg := new(chordMsgs.ChordMessage)
	command := chordMsgs.ChordMessage_Command(chordMsgs.ChordMessage_Command_value["GetPred"])
	chordMsg.Cmd = command
	pMsg := new(chordMsgs.PredMessage)
	fingerMsg := new(chordMsgs.FingerMessage)
	fingerMsg.Id = *proto.String(string(finger.id[:32]))
	fingerMsg.Address = *proto.String(finger.ipAddr)
	pMsg.Pred = fingerMsg
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

func claimpredMsg(finger Finger) []byte {
	msg := new(chordMsgs.NetworkMessage)
	msg.Proto = *proto.Uint32(1)
	chordMsg := new(chordMsgs.ChordMessage)
	command := chordMsgs.ChordMessage_Command(chordMsgs.ChordMessage_Command_value["ClaimPred"])
	chordMsg.Cmd = command
	predMsg := new(chordMsgs.PredMessage)
	fingerMsg := new(chordMsgs.FingerMessage)
	fingerMsg.Id = *proto.String(string(finger.id[:32]))
	fingerMsg.Address = *proto.String(finger.ipAddr)
	predMsg.Pred = fingerMsg
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

//pingMsg constructs a message to ping a server
func pingMsg() []byte {

	msg := new(chordMsgs.NetworkMessage)
	msg.Proto = *proto.Uint32(1)
	chordMsg := new(chordMsgs.ChordMessage)
	command := chordMsgs.ChordMessage_Command(chordMsgs.ChordMessage_Command_value["Ping"])
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

//pongMsg constructs a message to reply to a ping
func pongMsg() []byte {

	msg := new(chordMsgs.NetworkMessage)
	msg.Proto = *proto.Uint32(1)
	chordMsg := new(chordMsgs.ChordMessage)
	command := chordMsgs.ChordMessage_Command(chordMsgs.ChordMessage_Command_value["Pong"])
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

func getsuccessorsMsg() []byte {
	msg := new(chordMsgs.NetworkMessage)
	msg.Proto = *proto.Uint32(1)
	chordMsg := new(chordMsgs.ChordMessage)
	command := chordMsgs.ChordMessage_Command(chordMsgs.ChordMessage_Command_value["GetSucc"])
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

func nullMsg() []byte {
	msg := new(chordMsgs.NetworkMessage)
	msg.Proto = *proto.Uint32(1)

	data, err := proto.Marshal(msg)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}

	return data
}

//parseMessage takes as input an unmarshalled redis buffer and
//performs actions based on what the message contains.
func (node *ChordNode) parseMessage(data []byte, c chan []byte) {

	msg := new(chordMsgs.NetworkMessage)

	err := proto.Unmarshal(data, msg)
	if err != nil {
		logger.Error.Printf("%s,%s ", err.Error(), node.ipAddr)
		return
	}

	protocol := msg.GetProto()
	if protocol != 1 {
		/*
		if app, ok := node.applications[byte(protocol)]; ok {
			c <- app.Message([]byte(msg.GetMsg()))
		}
		return
		*/
	}

	chorddata := []byte(msg.GetMsg())
	chordmsg := new(chordMsgs.ChordMessage)
	err = proto.Unmarshal(chorddata, chordmsg)
	if err != nil {
		logger.Error.Printf("%s,%s ", err.Error(), node.ipAddr)
		return
	}

	cmd := int32(chordmsg.GetCmd())
	switch {
	case cmd == chordMsgs.ChordMessage_Command_value["Ping"]:
		c <- pongMsg()
		return
	case cmd == chordMsgs.ChordMessage_Command_value["GetPred"]:
		pred := *node.predecessor
		if pred.zero() {
			c <- nullMsg()
		} else {
			c <- sendpredMsg(pred)
		}
		return
	case cmd == chordMsgs.ChordMessage_Command_value["GetId"]:
		c <- sendidMsg(node.id[:32])
		return
	case cmd == chordMsgs.ChordMessage_Command_value["GetFingers"]:
		table := make([]Finger, 32*8+1)
		//fmt.Printf("Fingers of node %s:\n", node.ipaddr)
		for i := range table {
			node.request <- request{false, false, i}
			f := <-node.finger
			//fmt.Printf("\t%s\n", f.String())
			table[i] = f
		}

		c <- sendfingersMsg(table)
		return
	case cmd == chordMsgs.ChordMessage_Command_value["ClaimPred"]:
		//extract finger
		newPred, err := parseFinger(data)
		if err != nil {
			logger.Error.Printf("%s", err.Error())
			return
		}
		if err != nil {
			c <- nullMsg()
			break
		}
		node.request <- request{false, false, -1}
		pred := <-node.finger

		if pred.zero() || InRange(newPred.id, pred.id, node.id) {
			node.notify(newPred)
		}
		c <- nullMsg()
		//update finger table
		return
	case cmd == chordMsgs.ChordMessage_Command_value["GetSucc"]:
		table := make([]Finger, 32*8)
		for i := range table {
			node.request <- request{false, true, i}
			f := <-node.finger
			table[i] = f
		}

		c <- sendfingersMsg(table)
		return

	}
	fmt.Printf("No matching commands.\n")
}

//parseFingers can be called to return a finger table from a received
//parseFingers can be called to return a finger table from a received
//message after a getfingers call.
func parseFingers(data []byte) (ft []Finger, err error) {
	msg := new(chordMsgs.NetworkMessage)
	err = proto.Unmarshal(data, msg)
	if msg.GetProto() != 1 {
		//TODO: return non-nil error
		return
	}
	chorddata := []byte(msg.GetMsg())
	chordmsg := new(chordMsgs.ChordMessage)
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
	msg := new(chordMsgs.NetworkMessage)
	err = proto.Unmarshal(data, msg)
	if err != nil {
		logger.Error.Printf("%s", err.Error())
		return
	}

	if msg.GetMsg() == "" { //then received null msg instead. return nil
		return
	}

	chorddata := []byte(msg.GetMsg())
	chordmsg := new(chordMsgs.ChordMessage)
	err = proto.Unmarshal(chorddata, chordmsg)
	if err != nil {
		logger.Error.Printf("%s", err.Error())
		return
	}

	cpmsg := chordmsg.GetCpmsg()
	finger := cpmsg.GetPred()
	if finger == nil {
		return
	}
	copy(f.id[:], []byte(finger.Id))
	f.ipAddr = finger.Address

	return
}

func parseId(data []byte) (id [32]byte, err error) {
	msg := new(chordMsgs.NetworkMessage)
	err = proto.Unmarshal(data, msg)
	if msg.GetProto() != 1 {
		//TODO: return non-nil error
		return
	}

	chorddata := []byte(msg.GetMsg())
	chordmsg := new(chordMsgs.ChordMessage)
	err = proto.Unmarshal(chorddata, chordmsg)
	if err != nil {
		logger.Error.Printf("%s", err.Error())
		return
	}

	if chordmsg == nil { //then received null msg instead. return nil
		return
	}

	idmsg := chordmsg.GetSidmsg()
	arr := []byte(idmsg.GetId())
	copy(id[:], arr[:32])
	return
}

func parsePong(data []byte) (success bool, err error) {

	msg := new(chordMsgs.NetworkMessage)
	err = proto.Unmarshal(data, msg)
	if err != nil {
		logger.Error.Printf("%s", err.Error())
		return false, err
	}

	if msg.GetProto() != 1 {
		//TODO: return non-nil error
		fmt.Printf("Something went wrong!\n")
		return
	}

	chorddata := []byte(msg.GetMsg())
	chordmsg := new(chordMsgs.ChordMessage)
	err = proto.Unmarshal(chorddata, chordmsg)
	if err != nil {
		logger.Error.Printf("%s", err.Error())
		return
	}

	if chordmsg == nil { //then received null msg instead. return nil
		return
	}

	command := int32(chordmsg.GetCmd())
	if command == chordMsgs.ChordMessage_Command_value["Pong"] {
		success = true
	} else {
		success = false
	}

	return
}
