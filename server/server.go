package network

import (
	"github.com/leviathan1995/grape/config"
	"github.com/leviathan1995/grape/cache"
	"github.com/leviathan1995/grape/protocol"

	"net"
	"fmt"
	"log"
	"bufio"
)

func StartServer(config *config.Config, cache *cache.Cache) {
	go func() {
		listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", config.Ip, config.Port))
		if err != nil {
			panic(err)
		}
		defer  listen.Close()

		for {
			select {
			default:
				conn, err := listen.Accept()
				if err != nil {
					log.Println(err)
					continue
				}
				log.Println("Connect to ", conn.RemoteAddr())
				go handleConnection(&conn, cache)
			}
		}
	}()
}

func handleConnection(conn *net.Conn, cache *cache.Cache) {
	defer (*conn).Close()

	reader := bufio.NewReader(*conn)

	for {
		request, _, err := reader.ReadLine()
		if err != nil {
			// TODO
		}

		command, _:= protocol.Parser(string(request))
		cache.HandleCommand(command)

	}
}