// main
package main

import (
	"github.com/leviathan1995/grape/cache"
	"github.com/leviathan1995/grape/config"
	"github.com/leviathan1995/grape/consistent"
	"github.com/leviathan1995/grape/server"

	"flag"
	"log"
)

func printConf(conf config.Config) {
	log.Printf("address: %s", conf.Address)
	log.Printf("heart beat interval: %d", conf.HeartbeatInterval)
	for _, addr := range conf.RemotePeers {
		log.Printf("the peer-server: %s", addr)
	}
}

func main() {
	log.Printf("Grape is starting")

	var conf string
	flag.StringVar(&conf, "c", "config.yaml", "config file")
	flag.Parse()
	log.Printf("use %s as the config file", conf)
	config := config.LoadConfig(conf)

	printConf(*config)
	consistency := consistent.New()

	consistency.AddNode(config.Address)
	for _, peer := range config.RemotePeers {
		consistency.AddNode(peer)
	}

	cache := cache.NewCache(config, consistency, &server.PeerStatus)

	server.StartServer(config, cache)
}
