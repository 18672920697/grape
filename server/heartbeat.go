package server

import (
	"github.com/leviathan1995/grape/config"
	"github.com/leviathan1995/grape/cache"
	"github.com/leviathan1995/grape/protocol"
	"github.com/leviathan1995/grape/logger"

	"time"
	"net"
	"fmt"
	"bufio"
	"strconv"
	"strings"
	"log"
	"io"
)

const (
	heartbeatPackageSize = 1024
)

// All servers need to send heartbeat to other servers when it starts
func Heartbeat(config *config.Config, cache *cache.Cache) {
	ticker := time.NewTicker(time.Second * time.Duration(config.HeartbeatInterval))
	for range ticker.C {
		sendHeartbeat(cache)
	}
}

// Send heartbeat and update routetable
func sendHeartbeat(cache *cache.Cache) {
	var routeTable []string

	(*cache).RWMutex.RLock()
	for node := range *cache.RouteTable {
		routeTable = append(routeTable, node)
	}
	(*cache).RWMutex.RUnlock()

	for _, node := range routeTable {
		split := strings.Split(node, ":")
		ip := split[0]
		port, _ := strconv.Atoi(split[1])
		heartbeatPort := port + 1024
		monitorAddr := ip + ":" + strconv.Itoa(heartbeatPort)
		localAddr := (*cache).Config.Address

		tcpAddr, err := net.ResolveTCPAddr("tcp", monitorAddr)
		if err != nil {
			(*cache).RWMutex.Lock()
			(*cache.RouteTable)[node] = false
			(*cache).RWMutex.Unlock()
			continue
		}

		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			(*cache).RWMutex.Lock()
			(*cache.RouteTable)[node] = false
			(*cache).RWMutex.Unlock()
			continue
		}
		defer (*conn).Close()

		request := fmt.Sprintf("*2\r\n$4\r\nPING\r\n$%d\r\n%s\r\n", len(localAddr), localAddr)
		_, err = conn.Write([]byte(request))
		if err != nil {
			(*cache).RWMutex.Lock()
			(*cache.RouteTable)[node] = false
			(*cache).RWMutex.Unlock()
			continue
		}

		reply := make([]byte, heartbeatPackageSize)
		reader := bufio.NewReader(conn)

		_, err = reader.Read(reply)
		command, _ := protocol.Parser(string(reply))
		if err != nil {
			(*cache).RWMutex.Lock()
			(*cache.RouteTable)[node] = false
			(*cache).RWMutex.Unlock()
			continue
		}

		if command.Args[0] == "PONG" {
			(*cache).RWMutex.Lock()
			(*cache.RouteTable)[node] = true
			(*cache).RWMutex.Unlock()
			continue
		} else if command.Args[0] == "Deny heartbeat" {
			logger.Warning.Printf(command.Args[0]+" by %s", node)

			// Join cluster
			joinAddr, err := net.ResolveTCPAddr("tcp", node)
			joinConn, err := net.DialTCP("tcp", nil, joinAddr)
			if err != nil {
				continue
			}
			defer joinConn.Close()

			joinRequest := fmt.Sprintf("*2\r\n$4\r\nJOIN\r\n$%d\r\n%s\r\n", len(localAddr), localAddr)
			_, err = joinConn.Write([]byte(joinRequest))
			if err != nil {
				continue
			}
			logger.Info.Printf("Send join cluster request to %s", node)
			joinReply := make([]byte, heartbeatPackageSize)
			reader := bufio.NewReader(joinConn)

			_, err = reader.Read(joinReply)
			command, _ := protocol.Parser(string(joinReply))
			if err != nil {
				continue
			}
			if command.Args[0] == "OK" {
				logger.Info.Printf("Receive route table infomation from cluster, and update it")
				for index := 1; index < len(command.Args); index++ {
					(*cache).RWMutex.Lock()
					(*cache.RouteTable)[command.Args[index]] = false
					(*cache).RWMutex.Unlock()
				}
				// TODO print route table

			} else if command.Args[0] == "FAIL" {
				// TODO
			}
		}
	}
}

func MonitorHeartbeat(config *config.Config, cache *cache.Cache, monitorStart chan bool) {
	split := strings.Split(config.Address, ":")
	ip := split[0]
	port, _ := strconv.Atoi(split[1])
	heartbeatPort := port + 1024
	monitor := ip + ":" + strconv.Itoa(heartbeatPort)

	listen, err := net.Listen("tcp", fmt.Sprintf("%s", monitor))
	if err != nil {
		panic(err)
	}
	defer listen.Close()
	monitorStart <- true

	for {
		select {
		default:
			conn, err := listen.Accept()
			if err != nil {
				log.Println(err)
				continue
			}
			handleHeartbeat(&conn, cache)
		}
	}
}

func handleHeartbeat(conn *net.Conn, cache *cache.Cache) {
	request := make([]byte, heartbeatPackageSize)
	defer (*conn).Close()

	reader := bufio.NewReader(*conn)
	for {
		_, err := reader.Read(request)
		if err != nil {
			if err == io.EOF {
				return
			}
		}

		command, _ := protocol.Parser(string(request))

		var status protocol.Status
		var resp string
		switch strings.ToUpper(command.Args[0]) {
		case "PING":
			status, resp = cache.HandlePing(command.Args)
		default:
			status, resp = protocol.ProtocolNotSupport, ""
		}

		// Check where is heartbeat coming from
		remote := command.Args[1]

		(*cache).RWMutex.RLock()
		defer (*cache).RWMutex.RUnlock()

		if _, ok := (*cache.RouteTable)[remote]; !ok {
			(*conn).Write([]byte("-Deny heartbeat\r\n"))
			return
		}

		switch status {
		case protocol.RequestFinish:
			(*conn).Write([]byte(resp))
		case protocol.ProtocolNotSupport:
			(*conn).Write([]byte("-Protocol not support\r\n"))
		}
	}
}
