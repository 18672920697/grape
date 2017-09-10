package main

import (
	"github.com/leviathan1995/grape/cache"
	"github.com/leviathan1995/grape/config"
	"github.com/leviathan1995/grape/consistent"
	"github.com/leviathan1995/grape/server"
	"github.com/leviathan1995/grape/logger"

	"flag"
	"os"
)

func printConf(conf config.Config) {
	logger.Info.Printf("address: %s", conf.Address)
	logger.Info.Printf("heartbeat interval: %d", conf.HeartbeatInterval)
	for _, addr := range conf.RemotePeers {
		logger.Info.Printf("the peer-server: %s", addr)
	}
}

func main() {
	logger.Init(os.Stdout, os.Stdout, os.Stderr)

	logger.Info.Printf("Grape is starting")

	var conf string
	logger.Info.Printf("================Config================")
	flag.StringVar(&conf, "c", "config.yaml", "config file")
	flag.Parse()
	logger.Info.Printf("use %s as the config file", conf)
	config := config.LoadConfig(conf)
	printConf(*config)
	logger.Info.Printf("======================================")

	consistency := consistent.New()

	// Add node to consistency process
	consistency.AddNode(config.Address)
	for _, peer := range config.RemotePeers {
		consistency.AddNode(peer)
	}

	cache := cache.NewCache(config, consistency)

	server.StartServer(config, cache)
}
