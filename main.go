package main

import (
	"github.com/leviathan1995/grape/config"
	"github.com/leviathan1995/grape/server"
	"github.com/leviathan1995/grape/cache"
	"github.com/leviathan1995/grape/consistent"
)

func main() {

	config := config.LoadConfig()

	consistency := consistent.New()

	consistency.AddNode(config.Address)
	for _, peer := range config.RemotePeers {
		consistency.AddNode(peer)
	}

	cache := cache.NewCache(config, consistency)

	server.StartServer(config, cache)
}