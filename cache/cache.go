package cache

import (
	"github.com/leviathan1995/grape/config"
	"github.com/leviathan1995/grape/consistent"
	"github.com/leviathan1995/grape/protocol"

	"bytes"
	"fmt"
	"github.com/leviathan1995/grape/logger"
	"net"
	"strings"
	"sync"
)

type Cache struct {
	storage     *map[string]string
	Config      *config.Config
	consistency *consistent.Consistent
	RouteTable  *map[string]bool
	sync.Mutex
	sync.RWMutex
}

func NewCache(config *config.Config, consistency *consistent.Consistent) *Cache {
	storage := make(map[string]string)

	route := make(map[string]bool)
	for _, node := range config.RemotePeers {
		route[node] = false
	}
	cache := &Cache{
		storage:     &storage,
		Config:      config,
		consistency: consistency,
		RouteTable:  &route,
	}

	return cache
}

// Check this key whether store in node
func (cache *Cache) CheckKey(key string) (bool, string) {
	server, _ := cache.consistency.SetKey(key)
	if server != cache.Config.Address {
		return false, server
	} else {
		return true, ""
	}
}

func (cache *Cache) HandleCommand(data protocol.CommandData) (protocol.Status, string) {
	switch strings.ToUpper(data.Args[0]) {
	case "COMMAND":
		return protocol.ProtocolNotSupport, ""
	case "SET":
		return cache.HandleSet(data.Args)
	case "GET":
		return cache.HandleGet(data.Args)
	case "PING":
		return cache.HandlePing(data.Args)
	case "INFO":
		return cache.HandleInfo(data.Args)
	case "JOIN":
		return cache.HandleJoin(data.Args)
	default:
		return protocol.ProtocolNotSupport, ""
	}
}

func (cache *Cache) HandleSet(args []string) (protocol.Status, string) {
	key := args[1]
	// Check this key whether store in node
	if store, server := cache.CheckKey(key); !store {
		return protocol.ProtocolOtherNode, server
	}
	value := args[2]

	cache.Mutex.Lock()
	defer cache.Mutex.Unlock()
	(*cache.storage)[key] = value

	resp := fmt.Sprintf("+OK\r\n")
	return protocol.RequestFinish, resp
}

func (cache *Cache) HandleGet(args []string) (protocol.Status, string) {
	key := args[1]
	// Check this key whether store in node
	if store, server := cache.CheckKey(key); !store {
		return protocol.ProtocolOtherNode, server
	}

	if value, ok := (*cache.storage)[key]; ok {
		resp := fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)
		return protocol.RequestFinish, resp
	} else {
		return protocol.RequestNotFound, ""
	}
}

func (cache *Cache) HandlePing(args []string) (protocol.Status, string) {
	resp := fmt.Sprintf("+PONG\r\n")
	return protocol.RequestFinish, resp
}

func (cache *Cache) HandleInfo(args []string) (protocol.Status, string) {
	var resp bytes.Buffer

	num := fmt.Sprintf("*%d\r\n", len(*cache.RouteTable)+1)
	resp.WriteString(num)

	title := fmt.Sprintf("$%d\r\n%s\r\n", len("Connect status:"), "Connect status:")
	resp.WriteString(title)
	for peer, status := range *cache.RouteTable {
		var str_status string
		if status {
			str_status = "Up"
		} else {
			str_status = "Down"
		}
		peer_status := fmt.Sprintf("$%d\r\n%s\r\n", len(peer+": "+str_status), peer+": "+str_status)
		resp.WriteString(peer_status)
	}
	return protocol.RequestFinish, resp.String()
}

func (cache *Cache) HandleJoin(args []string) (protocol.Status, string) {
	joinAddr := args[1]

	routeResp := fmt.Sprintf("*%d\r\n$2\r\nOK\r\n", len(*cache.RouteTable)+1)

	var routeTable []string
	(*cache).RWMutex.RLock()
	for node, _ := range *cache.RouteTable {
		routeTable = append(routeTable, node)
	}
	(*cache).RWMutex.RUnlock()

	// Broadcast
	for _, node := range routeTable {
		nodeResp := fmt.Sprintf("$%d\r\n%s\r\n", len(node), node)
		routeResp += nodeResp

		nodeAddr, _ := net.ResolveTCPAddr("tcp", node)
		conn, err := net.DialTCP("tcp", nil, nodeAddr)
		if err != nil {
			continue
		}
		defer (*conn).Close()

		request := fmt.Sprintf("*2\r\n$4\r\nJOIN\r\n$%d\r\n%s\r\n", len(joinAddr), joinAddr)
		_, err = conn.Write([]byte(request))
		if err != nil {
			continue
		}
	}
	if joinAddr != (*cache).Config.Address {
		(*cache).RWMutex.Lock()
		if _, ok := (*cache.RouteTable)[joinAddr]; !ok {
			(*cache.RouteTable)[joinAddr] = false
			logger.Info.Printf("Add %s to route table", joinAddr)
		}
		(*cache).RWMutex.Unlock()
	}
	return protocol.RequestFinish, routeResp
}

func (cache *Cache) HandleRemove(args []string) (protocol.Status, string) {
	removeAddr := args[1]

	// Broadcast
	for node, _ := range *cache.RouteTable {
		nodeAddr, _ := net.ResolveTCPAddr("tcp", node)
		conn, err := net.DialTCP("tcp", nil, nodeAddr)
		if err != nil {
			continue
		}
		defer (*conn).Close()

		request := fmt.Sprintf("*2\r\n$4\r\nPING\r\n$%d\r\n%s\r\n", len(removeAddr), removeAddr)
		_, err = conn.Write([]byte(request))
		if err != nil {
			continue
		}
	}
	if _, ok := (*cache.RouteTable)[removeAddr]; ok {
		delete((*cache.RouteTable), removeAddr)
	}
	return protocol.RequestFinish, "OK"
}
