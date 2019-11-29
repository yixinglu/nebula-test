package nebulatest

import "github.com/vesoft-inc/nebula-go/graph"

type JsonDiffer struct {
	DifferError
	Response *graph.ExecutionResponse
	Order    bool
}

func (d *JsonDiffer) Diff(result string) {
}

func (d *JsonDiffer) PrintError(prefix string) {

}
