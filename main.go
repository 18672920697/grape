package main

import (
	"github.com/leviathan1995/grape/cache"
	"github.com/leviathan1995/grape/config"
	"github.com/leviathan1995/grape/consistent"
	"github.com/leviathan1995/grape/logger"
	"github.com/leviathan1995/grape/server"

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

	var confFile string
	logger.Info.Printf("================Config================")
	flag.StringVar(&confFile, "c", "config.yaml", "config file")
	flag.Parse()
	logger.Info.Printf("use %s as the config file", confFile)
	conf := config.LoadConfig(confFile)
	printConf(*conf)
	logger.Info.Printf("======================================")

	consistency := consistent.New()
	// Add node to consistency process
	consistency.AddNode(conf.Address)
	for _, peer := range conf.RemotePeers {
		consistency.AddNode(peer)
	}

	cacheStore := cache.NewCache(conf, consistency)

	server.StartServer(conf, cacheStore)
}
