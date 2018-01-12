package server

import (
	"github.com/leviathan1995/grape/cache"
	"github.com/leviathan1995/grape/config"
	"github.com/leviathan1995/grape/logger"
	"github.com/leviathan1995/grape/redis"

	"bufio"
	"fmt"
	"io"
	"net"
	"time"
)

const (
	receiveBufferSize = 1024 * 10 * 2
	sendBufferSize    = 1024 * 10 * 2
)

func StartServer(config *config.Config, cache *cache.Cache) {
	listen, err := net.Listen("tcp", fmt.Sprintf("%s", config.Address))
	if err != nil {
		logger.Error.Printf("Listen in %s failed.", config.Address)
	}
	defer listen.Close()

	joinCluster(config, cache)

	// Start service
	logger.Info.Printf("Start service.")
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
		successor, err := cache.Chord.Join(peers)
		if err != nil {
			logger.Error.Print(err)
		}
		if cache.Chord.GetSuccessorAddr() == "" && cache.Chord.GetNodeAddr() == successor.IpAddr() {
			cache.Chord.AfterJoin(peers)
		}
	}
}

func handleConnection(conn *net.Conn, cache *cache.Cache) {
	request := make([]byte, receiveBufferSize)

	reader := bufio.NewReader(*conn)
	for {
		length, err := reader.Read(request)
		if err != nil {
			if err == io.EOF {
				(*conn).Close()
				return
			}
		}

		command, _ := redis.Parser(string(request))
		status, resp := cache.HandleCommand(command)
		switch status {
		case redis.RequestFinish:
			(*conn).Write([]byte(resp))
		case redis.RequestNotFound:
			(*conn).Write([]byte("-Not found\r\n"))
		case redis.ProtocolNotSupport:
			(*conn).Write([]byte("-Protocol not support\r\n"))
		case redis.ProtocolOtherNode:
			resp = ForwardRequest(string(request[:length]), resp, cache)
			(*conn).Write([]byte(resp))
		}
	}
}

// The server could not response the client's request, maybe need to send to other servers
func ForwardRequest(request, addr string, cache *cache.Cache) string {
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
		conn.SetDeadline(time.Now().Add(1 * time.Minute))
		cache.Connections[addr] = *conn

		_, err = conn.Write([]byte(request))
		if err != nil {
			logger.Error.Printf("Write to the peer-server failed: %s", err.Error())
			return string("-Can not connect to destination Node\r\n")
		}
		reader := bufio.NewReader(conn)
		reply := make([]byte, sendBufferSize)
		length, err := reader.Read(reply)
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
		conn.SetDeadline(time.Now().Add(1 * time.Minute))
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
	reader := bufio.NewReader(&conn)
	reply := make([]byte, sendBufferSize)
	length, err := reader.Read(reply)
	if err != nil {
		logger.Error.Printf("Read from the peer-server failed: %s", err.Error())
		return string("-Can not connect to destination Node\r\n")
	}
	return string(reply[0:length])
}
