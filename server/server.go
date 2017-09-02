// server
package server

import (
	"github.com/leviathan1995/grape/cache"
	"github.com/leviathan1995/grape/config"
	"github.com/leviathan1995/grape/protocol"

	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

var PeerStatus map[string]bool

func StartServer(config *config.Config, cache *cache.Cache) {
	listen, err := net.Listen("tcp", fmt.Sprintf("%s", config.Address))
	if err != nil {
		panic(err)
	}
	defer listen.Close()

	// Send heartbeat
	PeerStatus = make(map[string]bool)
	sendHeartbeat(config)

	// The network checker
	// Start service

	// Send heartbeat to others at a fixed interval
	go Heartbeat(config)

	for {
		select {
		default:
			conn, err := listen.Accept()
			if err != nil {
				log.Println(err)
				continue
			}
			log.Println("Connect to", conn.RemoteAddr())
			go handleConnection(&conn, cache)
		}
	}
}

func handleConnection(conn *net.Conn, cache *cache.Cache) {
	request := make([]byte, 1024)
	defer (*conn).Close()

	reader := bufio.NewReader(*conn)

	for {
		_, err := reader.Read(request)
		if err != nil {
			if err == io.EOF {
				log.Printf("Close connection")
				(*conn).Close()
				return
			}
		}

		command, _ := protocol.Parser(string(request))
		status, resp := cache.HandleCommand(command)
		switch status {
		case protocol.RequestFinish:
			(*conn).Write([]byte(resp))
		case protocol.RequestNotFound:
			(*conn).Write([]byte("-Not found\r\n"))
		case protocol.ProtocolNotSupport:
			(*conn).Write([]byte("-Protocol not support\r\n"))
		case protocol.ProtocolOtherNode:
			resp = resendRequest(string(request), resp)
			(*conn).Write([]byte(resp))
		}
	}
}

// the server could not response the client's request, maybe need to send to other servers
func resendRequest(request, addr string) string {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Printf("ResolveTCPAddr failed: %s", err.Error())
		return string("-Can not connect to destination Node\r\n")
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Printf("Dial failed: %s", err.Error())
		return string("-Can not connect to destination Node\r\n")
	}
	_, err = conn.Write([]byte(request))
	if err != nil {
		log.Printf("Write to the peer-server failed: %s", err.Error())
		return string("-Can not connect to destination Node\r\n")
	}
	reply := make([]byte, 1024)
	_, err = conn.Read(reply)
	if err != nil {
		log.Printf("Read from the peer-server failed: %s", err.Error())
		return string("-Can not connect to destination Node\r\n")
	}
	return string(reply)
}

// all servers need to send heartbeat to other servers when it starts
func Heartbeat(config *config.Config)  {
	ticker := time.NewTicker(time.Second * time.Duration(config.HeartbeatInterval))
	for _ = range ticker.C {
		sendHeartbeat(config)
	}
}

func sendHeartbeat(config *config.Config) {
	for _, node := range config.RemotePeers {
		tcpAddr, err := net.ResolveTCPAddr("tcp", node)
		if err != nil {
			PeerStatus[node] = false
			continue
		}
		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			log.Printf("Dial failed: %s", err.Error())
			PeerStatus[node] = false
			continue
		}

		request := "*3\r\nPING\r\n" + config.Address + "\r\n"
		_, err = conn.Write([]byte(request))
		if err != nil {
			PeerStatus[node] = false
		}
		reply := make([]byte, 1024)
		_, err = conn.Read(reply)
		if err != nil {
			PeerStatus[node] = false
			continue
		}
		rep, _ := protocol.Parser(string(reply))
		if (rep.Args[0] == "PONG") {
			PeerStatus[node] = true
			continue
		}
		if (PeerStatus[node]) {
			log.Printf("Send heartbeat to %s succeeded", node)
		} else {
			log.Printf("Send heartbeat to %s failed", node)
		}
	}
}