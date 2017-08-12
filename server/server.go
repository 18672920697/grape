package server

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
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", config.Ip, config.Port))
	if err != nil {
		panic(err)
	}
	defer listen.Close()

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
}

func handleConnection(conn *net.Conn, cache *cache.Cache) {
	request := make([]byte, 1024)
	defer (*conn).Close()

	reader := bufio.NewReader(*conn)

	for {
		_, err := reader.Read(request)
		if err != nil {
			log.Printf("Read data from %s failed", (*conn).RemoteAddr())
		}

		command, _:= protocol.Parser(string(request))
		status, resp := cache.HandleCommand(command)
		switch status {
		case protocol.RequestFinish:
			(*conn).Write([]byte(resp))
		case protocol.ProtocolNotSupport:
			(*conn).Write([]byte("-Protocol not support\r\n"))
		}
	}
}