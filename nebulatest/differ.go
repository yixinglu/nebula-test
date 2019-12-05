package nebulatest

import (
	"fmt"
	"strings"

	"github.com/vesoft-inc/nebula-go/graph"
)

type Differ interface {
	Diff(result string)
	Error() error
}

type DifferError struct {
	err error
}

func (d *DifferError) Error() error {
	return d.err
}

func NewDiffer(resp *graph.ExecutionResponse, dType string, order bool) (Differ, error) {
	switch strings.ToLower(dType) {
	case "json":
		return &JsonDiffer{
			Response: resp,
			Order:    order,
		}, nil
	case "table":
		return &TableDiffer{
			Response: resp,
		}, nil
	default:
		return nil, fmt.Errorf("Invalid differ type: %s", dType)
	}
}
