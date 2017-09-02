// cache
package cache

import (
	"github.com/leviathan1995/grape/config"
	"github.com/leviathan1995/grape/consistent"
	"github.com/leviathan1995/grape/protocol"

	"fmt"
	"strings"
	"sync"
	"bytes"
)

type Cache struct {
	storage     *map[string]string
	config      *config.Config
	consistency *consistent.Consistent
	routeTable []string
	peerStatus *map[string]bool
	sync.Mutex
}

func NewCache(config *config.Config, consistency *consistent.Consistent, peerStatus *map[string]bool) *Cache {
	storage := make(map[string]string)

	cache := &Cache{
		storage:     &storage,
		config:      config,
		consistency: consistency,
		peerStatus: peerStatus,
	}
	for _, node := range config.RemotePeers {
		cache.routeTable = append(cache.routeTable, node)
	}
	return cache
}

// Check this key whether store in node
func (cache *Cache) CheckKey(key string) (bool, string) {
	server, _ := cache.consistency.SetKey(key)
	if server != cache.config.Address {
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

	cache.Lock()
	defer cache.Unlock()
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

	num := fmt.Sprintf("*%d\r\n", len(*cache.peerStatus)+1)
	resp.WriteString(num)

	title := fmt.Sprintf("$%d\r\n%s\r\n", len("Connect status:"), "Connect status:")
	resp.WriteString(title)
	for peer, status := range *cache.peerStatus{
		var str_status string
		if status {
			str_status = "Up"
		} else {
			str_status = "Down"
		}
		peer_status := fmt.Sprintf("$%d\r\n%s\r\n", len(peer + ": " + str_status), peer + ": " + str_status)
		resp.WriteString(peer_status)
	}
	return protocol.RequestFinish, resp.String()
}
