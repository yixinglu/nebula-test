package parser

import (
	"bufio"
	"io/ioutil"
	"log"
	"strings"
)

func ReadFile(filename string) error {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(bytes)))
	for scanner.Scan() {
		parse(scanner.Text())
	}

	return nil
}

// TODO(yee): Use FSM to implement parse
func parse(text string) {
	if strings.HasPrefix(text, "=== test") {
		log.Println(text)
	} else if text == "--- in" {
		log.Println("in")
	} else if text == "--- out" {
		log.Println("out")
	} else {
		log.Println(text)
	}
}
