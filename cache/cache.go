package cache

import (
	"github.com/leviathan1995/grape/config"
	"github.com/leviathan1995/grape/protocol"
	"sync"
	"strings"
	"fmt"
)

type Cache struct {
	storage *map[string]string
	sync.Mutex
}
func NewCache(c *config.Config) *Cache {
	storage := make(map[string]string)

	cache := &Cache{
		storage: &storage,
	}
	return cache
}

func (cache *Cache) HandleCommand(data protocol.CommandData) (protocol.Status, string){
	switch strings.ToUpper(data.Args[0]) {
	case "COMMAND":
		{
			return protocol.ProtocolNotSupport, ""
		}
	case "SET":
		{
			return cache.HandleSet(data.Args)
		}
	case "GET":
		{
			return cache.HandleGet(data.Args)
		}
	default:
		{
			return protocol.ProtocolNotSupport , ""
		}
	}
}

func (cache *Cache) HandleSet(args []string) (protocol.Status, string) {
	key := args[1]
	value := args[2]

	cache.Lock()
	defer cache.Unlock()
	(*cache.storage)[key] = value

	resp := fmt.Sprintf("+OK\r\n")
	return protocol.RequestFinish, resp
}

func (cache *Cache) HandleGet(args []string) (protocol.Status, string) {
	key := args[1]

	if value, ok := (*cache.storage)[key]; ok {
		resp := fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)
		return protocol.RequestFinish, resp
	} else {
		return protocol.RequestNotFound, ""
	}
}