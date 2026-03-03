package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"time"
)

const (
	protocol   = "tcp"
	addrClient = "localhost:8080"
	respOk     = "OK\n"
	timeout    = time.Second * 5
)

func main() {
	conn, err := net.Dial(protocol, addrClient)
	if err != nil {
		log.Printf("Error connecting to server: %v", err)
		os.Exit(1)
	}
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(timeout))

	r := bufio.NewReader(conn)
	line, err := r.ReadString('\n')

	if err != nil {
		log.Printf("Error reading from server: %v", err)
		os.Exit(1)
	}
	if line != respOk {
		log.Printf("Unexpected response: got %q, want %q", line, respOk)
		os.Exit(1)
	}
}
