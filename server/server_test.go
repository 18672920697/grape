package server

import (
	"flag"
	"github.com/leviathan1995/grape/cache"
	"github.com/leviathan1995/grape/config"
	"github.com/leviathan1995/grape/consistent"
	"testing"
)

func TestForwardRequest(t *testing.T) {

	var conf string
	flag.StringVar(&conf, "c", "server_test.yaml", "config file")
	flag.Parse()
	config := config.LoadConfig(conf)

	consistency := consistent.New()

	// Add node to consistency process
	consistency.AddNode(config.Address)
	for _, peer := range config.RemotePeers {
		consistency.AddNode(peer)
	}

	cache := cache.NewCache(config, consistency)

	req := "*3\r\n$3\r\nset\r\n$3\r\nkey\r\n$3\r\nval\r\n"
	for i := 0; i < 10000; i++ {
		reply := ForwardRequest(req, "127.0.0.1:6002", cache)
		if reply != "+OK\r\n" {
			t.Errorf(reply)
		}
	}
}
