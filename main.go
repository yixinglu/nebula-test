package nebula_test

import (
	"flag"
	"fmt"
	"log"

	nebula "github.com/vesoft-inc/nebula-go"
)

type NebulaConfig struct {
	NebulaTestUser     string
	NebulaTestPassword string
}

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

	client, err := nebula.NewClient(fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		panic(err)
	}
	err = client.Open()
	if err != nil {
		panic(err)
	}
	defer client.Close()

	nebulaConf := NebulaConfig{
		NebulaTestUser:     *username,
		NebulaTestPassword: *password,
	}

	err = Parse(*filename, client, &nebulaConf)
	if err != nil {
		panic(err)
	}
}
