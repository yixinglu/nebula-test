package nebulatest

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	nebula "github.com/vesoft-inc/nebula-go"
	"github.com/vesoft-inc/nebula-go/graph"
)

const (
	testPrefix = "=== test"
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
	var testName string
	var response *graph.ExecutionResponse
	var differ Differ
	var wait time.Duration
	isInput, isOutput := false, false
	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, testPrefix) {
			if isOutput {
				differ.Diff(outBuf.String())
				differ.PrintError(testName)
				outBuf.Reset()
				isOutput = false
			}

			// Reset test comment after last test output result
			prefixLen := len(testPrefix)
			if prefixLen > len(text) {
				return fmt.Errorf("%s length is larger than %s", testPrefix, text)
			}
			testName = strings.TrimLeft(strings.TrimSpace(text[prefixLen:]), ": ")
		} else if strings.HasPrefix(text, inPrefix) {
			isInput = true
			w := strings.TrimLeft(strings.TrimSpace(text[len(inPrefix):]), ": ")
			if wait, err = tester.parseInputWait(w); err != nil {
				return err
			}
		} else if strings.HasPrefix(text, outPrefix) {
			isOutput = true

			if isInput {
				time.Sleep(wait)
				if response, err = tester.request(inBuf.String()); err != nil {
					return err
				}
				if d, err := tester.newDiffer(text, response); err != nil {
					return err
				} else {
					differ = d
				}
				isInput = false
				inBuf.Reset()
			}
		} else {
			if isInput {
				if !strings.HasPrefix(text, "--") && !strings.HasPrefix(text, "#") && !strings.HasPrefix(text, "//") {
					// text = fmt.Sprintf("%q", text)
					text = strings.TrimRight(text, "\\ \"")
					text = strings.TrimLeft(text, "\"")
					inBuf.WriteString(text)
				}
			}
			if isOutput {
				if outBuf.Len() > 0 {
					outBuf.WriteString("\n")
				}
				outBuf.WriteString(text)
			}
		}
	}

	if isOutput {
		differ.Diff(outBuf.String())
		differ.PrintError(testName)
		outBuf.Reset()
	}

	return nil
}

func (tester *Tester) request(gql string) (*graph.ExecutionResponse, error) {
	gql = strings.TrimSpace(gql)
	resp, err := tester.client.Execute(gql)
	if err != nil {
		return nil, err
	}

	if resp.GetErrorCode() != graph.ErrorCode_SUCCEEDED {
		return nil, fmt.Errorf("ErrorCode: %v, ErrorMsg: %s", resp.GetErrorCode(), resp.GetErrorMsg())
	}

	return resp, nil
}

func (tester *Tester) newDiffer(outText string, response *graph.ExecutionResponse) (Differ, error) {
	dType, order := "table", false
	index := strings.Index(outText, ",")
	if index >= 0 {
		index = strings.Index(outText, ":")
		dType, order = tester.getOptions(outText[index+1:])
	}
	if differ, err := NewDiffer(response, dType, order); err != nil {
		return nil, err
	} else {
		return differ, nil
	}
}

func (t *Tester) getOptions(config string) (dType string, order bool) {
	options := strings.Split(config, ",")
	dType = "table"
	order = false
	for _, op := range options {
		if index := strings.Index(op, "="); index < 0 {
			continue
		}
		kv := strings.Split(op, "=")
		key := strings.Trim(strings.ToLower(kv[0]), " ")
		value := strings.Trim(strings.ToLower(kv[1]), " ")
		switch key {
		case "type":
			dType = value
		case "order":
			if b, err := strconv.ParseBool(value); err != nil {
				log.Printf("Invalid order type: %s", kv[1])
			} else {
				order = b
			}
		default:
			log.Fatalf("Unvalid key: %s", key)
		}
	}
	return dType, order
}

func (t *Tester) parseInputWait(s string) (time.Duration, error) {
	if len(s) == 0 {
		return time.ParseDuration("0s")
	}
	kv := strings.Split(s, "=")
	if len(kv) != 2 || strings.ToLower(kv[0]) != "wait" {
		log.Println("Invalid option format, like wait=10s")
		return time.ParseDuration("0s")
	}
	return time.ParseDuration(strings.TrimSpace(kv[1]))
}
