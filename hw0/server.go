package main

import (
	"io"
	"log"
	"net"
)

const (
	network    = "tcp"
	addrServer = ":8080"
	responseOk = "OK\n"
)

func main() {
	server, err := net.Listen(network, addrServer)
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()
	for {
		conn, err := server.Accept()
		if err != nil {
			continue
		}
		go func(c net.Conn) {
			_, _ = io.WriteString(c, responseOk)
			defer c.Close()
		}(conn)
	}
}
