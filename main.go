package main

import (
	"flag"
	"log"

	nebula "github.com/vesoft-inc/nebula-go"
	nt "github.com/vesoft-inc/nebula-test/nebulatest"
)

func main() {
	file := flag.String("file", "", "Test file path")
	username := flag.String("user", "user", "Nebula username")
	password := flag.String("password", "password", "Nebula password")
	address := flag.String("address", "127.0.0.1:3699", "Nebula Graph server ip address and port")
	flag.Parse()

	if *file == "" {
		log.Println("Please input a test filename")
		return
	}

	client, err := nebula.NewClient(*address)
	if err != nil {
		log.Fatal(err)
	}

	if err = client.Connect(*username, *password); err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect()

	t := nt.NewTester(client)

	if err = t.Parse(*file); err != nil {
		log.Fatal(err)
	}
}
