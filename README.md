# Grape

[![Build Status](https://travis-ci.org/Leviathan1995/grape.svg?branch=master)](https://travis-ci.org/Leviathan1995/grape)
[![Go Report Card](https://goreportcard.com/badge/github.com/leviathan1995/grape)](https://goreportcard.com/report/github.com/leviathan1995/grape)
[![Hex.pm](https://img.shields.io/hexpm/l/plug.svg)]()

## Introduction
Grape is a decentralized distribution memory caching system, using consistent hashing to decide which node to store the data in. Grape support redis SET and GET command, it's easy to use the `redis-cli` to add or remove a node from cluster. Because the grape like a architecture of peer-to-peer systems, you can use redis client to connect any of a node of cluster, and the connecting node will send request to destination node.

## Installation
	go get -v github.com/leviathan1995/grape/
## Getting Started
### Creating and using a cluster
First start a double-member cluster, each node set the peer-server as `RemotePeers` in config file
		
	./grape -c tests/node1.yaml
	./grape -c tests/node2.yaml

Next, connect to any of a node and store a key-value and retrieve the stored key
	
	./redis-cli -p 9221
	127.0.0.1:9221> SET key value
	OK
	127.0.0.1:9221> GET key
	"value"

### Adding a new node
if you want to add a node3 to the exist cluster, need to configurate any of a node in cluster to `RemotePeers` in config file
	
	./grape -c test/node3.yaml

and it will be automatically join the cluster

### Removing a node
Connect any of a node in cluster, use `REMOVE IP PORT`
	
	./redis-cli -p 9221
	127.0.0.1:9221> REMOVE 127.0.0.1 9001
	OK
	
## TODO
* TTL
* LRU algorithm
* Monitor nodes survival

