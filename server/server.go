package server

import (
	"github.com/leviathan1995/grape/cache"
	"github.com/leviathan1995/grape/config"
	"github.com/leviathan1995/grape/logger"
	"github.com/leviathan1995/grape/protocol"

	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

func StartServer(config *config.Config, cache *cache.Cache) {
	listen, err := net.Listen("tcp", fmt.Sprintf("%s", config.Address))
	if err != nil {
		panic(err)
	}
	defer listen.Close()

	// Monitor heartbeat port
	go MonitorHeartbeat(config, cache)

	// Check the networks of cluster
	logger.Info.Printf("Wait for all nodes connected ")
	for !ClusterConnected(cache.RouteTable) {
		sendHeartbeat(cache)
	}

	for peer, status := range *cache.RouteTable {
		if status {
			logger.Info.Printf("Connecting to node %s OK", peer)
		}
	}
	logger.Info.Printf("Create cluster success...")

	// Send heartbeat to others at a fixed interval
	go Heartbeat(config, cache)

	// Start service
	logger.Info.Printf("Start service...")
	for {
		select {
		default:
			conn, err := listen.Accept()
			if err != nil {
				log.Println(err)
				continue
			}
			//logger.Info.Println("Connect to", conn.RemoteAddr())
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
				//logger.Info.Printf("Close connection")
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

// The server could not response the client's request, maybe need to send to other servers
func resendRequest(request, addr string) string {
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
	_, err = conn.Write([]byte(request))
	if err != nil {
		logger.Error.Printf("Write to the peer-server failed: %s", err.Error())
		return string("-Can not connect to destination Node\r\n")
	}
	reply := make([]byte, 1024)
	_, err = conn.Read(reply)
	if err != nil {
		logger.Error.Printf("Read from the peer-server failed: %s", err.Error())
		return string("-Can not connect to destination Node\r\n")
	}
	return string(reply)
}

// All servers need to send heartbeat to other servers when it starts
func Heartbeat(config *config.Config, cache *cache.Cache) {
	ticker := time.NewTicker(time.Second * time.Duration(config.HeartbeatInterval))

	for range ticker.C {
		sendHeartbeat(cache)
		ClusterConnected(cache.RouteTable)
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
		heartbeartPort := port + 1024
		monitor := ip + ":" + strconv.Itoa(heartbeartPort)
		localAddr := (*cache).Config.Address

		tcpAddr, err := net.ResolveTCPAddr("tcp", monitor)
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

		reply := make([]byte, 1024)
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
			joinRequest := fmt.Sprintf("*2\r\n$4\r\nJOIN\r\n$%d\r\n%s\r\n", len(localAddr), localAddr)
			_, err = joinConn.Write([]byte(joinRequest))
			if err != nil {
				continue
			}
			logger.Info.Printf("Send join cluster request to %s", node)
			joinReply := make([]byte, 1024)
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
			joinConn.Close()
		}
	}
}

func ClusterConnected(RouteTable *map[string]bool) bool {
	for _, status := range *RouteTable {
		if status == false {
			//logger.Warning.Printf("Can not connect to %s", peer)
			return false
		}
	}
	return true
}

func MonitorHeartbeat(config *config.Config, cache *cache.Cache) {
	split := strings.Split(config.Address, ":")
	ip := split[0]
	port, _ := strconv.Atoi(split[1])
	heartbeartPort := port + 1024
	monitor := ip + ":" + strconv.Itoa(heartbeartPort)

	listen, err := net.Listen("tcp", fmt.Sprintf("%s", monitor))
	logger.Info.Printf("Monitor heartbeat address: %s:%d", ip, heartbeartPort)
	if err != nil {
		panic(err)
	}
	defer listen.Close()

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
	request := make([]byte, 1024)
	defer (*conn).Close()

	reader := bufio.NewReader(*conn)
	for {
		_, err := reader.Read(request)
		if err != nil {
			if err == io.EOF {
				//logger.Info.Printf("Close connection: %s", (*conn).LocalAddr())
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
