package main

import (
	"flag"
	"log"

	"github.com/vesoft-inc/nebula-test/parser"
)

func main() {
	filename := flag.String("filename", "", "Test filename")
	flag.Parse()

	if *filename == "" {
		log.Println("Please input a test filename")
		return
	}

	err := parser.ReadFile(*filename)
	if err != nil {
		panic(err)
	}
}
