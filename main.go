package main

import (
	"github.com/leviathan1995/grape/config"
	"github.com/leviathan1995/grape/server"
	"github.com/leviathan1995/grape/cache"

)

func main() {

	config := config.LoadConfig()

	cache := cache.NewCache(config)
	network.StartServer(config, cache)
}