package main

import (
	"hw2/app"
	"log"
)

func main() {
	server, err := app.NewServer()

	if err != nil {
		log.Fatal(err)
	}
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
