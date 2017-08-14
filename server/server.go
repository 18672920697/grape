package server

import (
	"github.com/leviathan1995/grape/config"
	"github.com/leviathan1995/grape/cache"
	"github.com/leviathan1995/grape/protocol"

	"net"
	"fmt"
	"log"
	"bufio"
	//"time"
	"io"
)

func StartServer(config *config.Config, cache *cache.Cache) {
	listen, err := net.Listen("tcp", fmt.Sprintf("%s", config.Address))
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
			log.Println("Connect to", conn.RemoteAddr())
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
			if err == io.EOF {
				log.Printf("Close connection")
				(*conn).Close()
				return
			}
		}

		command, _:= protocol.Parser(string(request))
		status, resp := cache.HandleCommand(command)
		switch status {
		case protocol.RequestFinish:
			(*conn).Write([]byte(resp))
		case protocol.ProtocolNotSupport:
			(*conn).Write([]byte("-Protocol not support\r\n"))
		case protocol.ProtocolOtherNode:
			// TODO
			resp = resendRequest(string(request), resp)
			(*conn).Write([]byte(resp))
		}
	}
}

func resendRequest(request, addr string) string {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Printf("ResolveTCPAddr failed:", err.Error())
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err !=nil {
		log.Printf("Dial failed:", err.Error())
	}
	_, err = conn.Write([]byte(request))
	if err != nil {
		log.Printf("Write to the peer-server failed:", err.Error())
	}
	reply := make([]byte, 1024)
	_, err = conn.Read(reply)
	if err != nil {
		log.Printf("Read from the peer-server failed:", err.Error())
	}
	return string(reply)
}