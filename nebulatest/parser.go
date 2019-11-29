package nebulatest

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	nebula "github.com/vesoft-inc/nebula-go"
	"github.com/vesoft-inc/nebula-go/graph"
)

const (
	testPrefix = "=== test:"
	inPrefix   = "--- in"
	outPrefix  = "--- out"
)

type Tester struct {
	client *nebula.GraphClient
}

func NewTester(client *nebula.GraphClient) *Tester {
	return &Tester{
		client: client,
	}
}

func (tester *Tester) Parse(filename string) error {
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
				return fmt.Errorf("%s length is larger than %s", testPrefix, text)
			}
			testName = strings.TrimSpace(text[prefixLen:])
		} else if strings.HasPrefix(text, inPrefix) {
			isInput = true
		} else if strings.HasPrefix(text, outPrefix) {
			isOutput = true

			if isInput {
				if respResult, err = tester.request(inBuf.String()); err != nil {
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
	}

	return nil
}

func (tester *Tester) request(gql string) (string, error) {
	gql = strings.TrimSpace(gql)
	resp, err := tester.client.Execute(gql)
	if err != nil {
		return "", err
	}

	if resp.GetErrorCode() != graph.ErrorCode_SUCCEEDED {
		return "", fmt.Errorf("ErrorCode: %v, ErrorMsg: %s", resp.GetErrorCode(), resp.GetErrorMsg())
	}

	return PrintResult(resp), nil
}

func diff(testName, expected, real string) {
	expected = strings.TrimSpace(expected)
	real = strings.TrimSpace(real)
	if expected != real {
		log.Fatalf("expected:\n%s, real:\n%s", expected, real)
	} else {
		log.Printf("Test (%s) passed", testName)
	}
}
