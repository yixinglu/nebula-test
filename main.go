package main

import (
	"flag"
	"fmt"
	"log"

	nebula "github.com/vesoft-inc/nebula-go"
	nt "github.com/vesoft-inc/nebula-test/nebulatest"
)

func main() {
	filename := flag.String("filename", "", "Test filename")
	username := flag.String("username", "user", "Nebula username")
	password := flag.String("password", "password", "Nebula password")
	address := flag.String("address", "127.0.0.1", "Nebula Graph server ip address")
	port := flag.Int64("port", 3699, "Nebula Graph server ip port")
	flag.Parse()

	if *filename == "" {
		log.Println("Please input a test filename")
		return
	}

	client, err := nebula.NewClient(fmt.Sprintf("%s:%d", *address, *port))
	if err != nil {
		log.Fatal(err)
	}

	if err = client.Connect(*username, *password); err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect()

	if err = nt.Parse(*filename, client); err != nil {
		log.Fatal(err)
	}
}