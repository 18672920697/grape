# Grape

[![Build Status](https://travis-ci.org/Leviathan1995/grape.svg?branch=master)](https://travis-ci.org/Leviathan1995/grape)
[![Go Report Card](https://goreportcard.com/badge/github.com/leviathan1995/grape)](https://goreportcard.com/report/github.com/leviathan1995/grape)
[![Hex.pm](https://img.shields.io/hexpm/l/plug.svg)]()

## Introduction
Grape is a centralized distribution memory caching system, using consistent hashing to decide which node to store the data in. Grape support redis SET and GET command, it's easy to use the `redis-cli` to add or remove a node from cluster. Because the grape like a architecture of peer-to-peer systems, you can use redis client to connect any of a node of cluster, and the connecting node will send request to destination node.

## Installation
	go get -v github.com/leviathan1995/grape/
## Example Usage


## TODO
* LRU algorithm
* Monitor nodes survival

