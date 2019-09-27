package nebulatest

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	nebula "github.com/vesoft-inc/nebula-go"
)

const (
	testPrefix = "=== test:"
	inPrefix   = "--- in"
	outPrefix  = "--- out"
)

func Parse(filename string, client *nebula.GraphClient, nebulaConf *NebulaConfig) error {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(b)))
	// TODO(yee): Use FSM to implement parse
	var inBuf, outBuf bytes.Buffer
	var respResult, testName string
	isInput, isOutput := false, false
	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, testPrefix) {
			if isOutput {
				diff(testName, outBuf.String(), respResult)
				outBuf.Reset()
				isOutput = false
			}

			// Reset test comment after last test output result
			prefixLen := len(testPrefix)
			if prefixLen > len(text) {
				return errors.New(fmt.Sprintf("%s length is larger than %s", testPrefix, text))
			}
			testName = strings.TrimSpace(text[prefixLen:])
		} else if strings.HasPrefix(text, inPrefix) {
			isInput = true
		} else if strings.HasPrefix(text, outPrefix) {
			isOutput = true

			if isInput {
				if respResult, err = request(inBuf.String(), client, nebulaConf); err != nil {
					return err
				}
				isInput = false
				inBuf.Reset()
			}
		} else {
			if isInput {
				inBuf.WriteString(text)
			}
			if isOutput {
				outBuf.WriteString(text)
			}
		}
	}

	if isOutput {
		diff(testName, outBuf.String(), respResult)
		outBuf.Reset()
		isOutput = false
	}

	return nil
}

func request(gql string, client *nebula.GraphClient, nebulaConf *NebulaConfig) (string, error) {
	gql = strings.TrimSpace(gql)
	authResp, err := client.Authenticate(nebulaConf.NebulaTestUser, nebulaConf.NebulaTestPassword)
	if err != nil {
		return "", err
	}
	resp, err := client.Execute(*authResp.SessionID, gql)
	if err != nil {
		return "", err
	}
	// TODO(yee): return response results
	return resp.String(), nil
}

// TODO(yee): diff output result and response result
func diff(testName, expected, real string) {
	expected = strings.TrimSpace(expected)
	real = strings.TrimSpace(real)
	if expected != real {
		log.Printf("expected: %s, real: %s", expected, real)
	} else {
		log.Printf("Test (%s) passed", testName)
	}
}
