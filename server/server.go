package server

import (
	"github.com/leviathan1995/grape/cache"
	"github.com/leviathan1995/grape/config"
	"github.com/leviathan1995/grape/logger"
	"github.com/leviathan1995/grape/protocol"

	"bufio"
	"fmt"
	"io"
	"net"
	"time"
)

const (
	receiveBufferSize = 1024 * 4
	sendBufferSize    = 1024 * 4
)

func StartServer(config *config.Config, cache *cache.Cache) {
	listen, err := net.Listen("tcp", fmt.Sprintf("%s", config.Address))
	if err != nil {
		panic(err)
	}
	defer listen.Close()

	/*
		monitorStart := make(chan bool)
		// Monitor heartbeat port
		go MonitorHeartbeat(config, cache, monitorStart)

		select {
		case start := <-monitorStart:
			if start {
				logger.Info.Printf("Heartbeat monitor start...")
			}
		}

		// Check the networks of cluster
		logger.Info.Printf("Wait for all nodes connected ")
		for !ClusterConnected(cache) {
			sendHeartbeat(cache)
		}

		(*cache).RLock()
		for peer, status := range *cache.RouteTable {
			if status {
				logger.Info.Printf("Connecting to node %s OK", peer)
			}
		}
		(*cache).RUnlock()
		logger.Info.Printf("Create cluster success...")

		// Send heartbeat to others at a fixed interval of time
		go Heartbeat(config, cache)
	*/

	joinCluster(config, cache)

	// Start service
	logger.Info.Printf("Start service...")
	for {
		select {
		default:
			conn, err := listen.Accept()
			if err != nil {
				logger.Error.Printf("%s", err)
				continue
			}
			go handleConnection(&conn, cache)
		}
	}
}

func joinCluster(config *config.Config, cache *cache.Cache) {
	for _, peers := range config.RemotePeers {
		err := cache.Chord.Join(config.Address, peers)
		if err != nil {
			logger.Error.Print(err)
		}
	}
}

func handleConnection(conn *net.Conn, cache *cache.Cache) {
	request := make([]byte, receiveBufferSize)
	defer (*conn).Close()

	reader := bufio.NewReader(*conn)

	for {
		len, err := reader.Read(request)
		if err != nil {
			if err == io.EOF {
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
			resp = resendRequest(string(request[:len]), resp, cache)
			(*conn).Write([]byte(resp))
		}
	}
}

// The server could not response the client's request, maybe need to send to other servers
func resendRequest(request, addr string, cache *cache.Cache) string {
	conn, ok := cache.Connections[addr]
	if !ok {
		tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
		if err != nil {
			logger.Error.Printf("ResolveTCPAddr failed: %s", err.Error())
			return string("-Can not connect to destination Node\r\n")
		}
		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			logger.Error.Printf("Dial failed: %s", err.Error())
			return string("-Can not connect to destination Node\r\n")
		}
		conn.SetDeadline(time.Now().Add(3 * time.Minute))
		cache.Connections[addr] = *conn

		_, err = conn.Write([]byte(request))
		if err != nil {
			logger.Error.Printf("Write to the peer-server failed: %s", err.Error())
			return string("-Can not connect to destination Node\r\n")
		}
		reply := make([]byte, sendBufferSize)
		length, err := conn.Read(reply)
		if err != nil {
			logger.Error.Printf("Read from the peer-server failed: %s", err.Error())
			return string("-Can not connect to destination Node\r\n")
		}
		return string(reply[0:length])
	}

	_, err := conn.Write([]byte(request))
	if err != nil {
		tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
		if err != nil {
			logger.Error.Printf("ResolveTCPAddr failed: %s", err.Error())
			return string("-Can not connect to destination Node\r\n")
		}
		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			logger.Error.Printf("Dial failed: %s", err.Error())
			return string("-Can not connect to destination Node\r\n")
		}
		conn.SetDeadline(time.Now().Add(3 * time.Minute))
		cache.Connections[addr] = *conn

		_, err = conn.Write([]byte(request))
		if err != nil {
			logger.Error.Printf("Write to the peer-server failed: %s", err.Error())
			return string("-Can not connect to destination Node\r\n")
		}
		reply := make([]byte, sendBufferSize)
		length, err := conn.Read(reply)
		if err != nil {
			logger.Error.Printf("Read from the peer-server failed: %s", err.Error())
			return string("-Can not connect to destination Node\r\n")
		}
		return string(reply[0:length])
	}
	reply := make([]byte, sendBufferSize)
	length, err := conn.Read(reply)
	if err != nil {
		logger.Error.Printf("Read from the peer-server failed: %s", err.Error())
		return string("-Can not connect to destination Node\r\n")
	}
	return string(reply[0:length])
}

func ClusterConnected(cache *cache.Cache) bool {
	(*cache).RLock()
	defer (*cache).RUnlock()

	for _, status := range *cache.RouteTable {
		if status == false {
			return status
		}
	}

	return true
}
