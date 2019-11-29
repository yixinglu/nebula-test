package nebulatest

import (
	"encoding/json"
	"fmt"

	"github.com/vesoft-inc/nebula-go/graph"
)

type JsonDiffer struct {
	DifferError
	Response *graph.ExecutionResponse
	Order    bool
}

func (d *JsonDiffer) Diff(result string) {
	var resp graph.ExecutionResponse
	if err := json.Unmarshal([]byte(result), &resp); err != nil {
		d.err = err
	} else {
		if d.Order {
			if !d.compareInOrder(&resp) {
				d.err = fmt.Errorf("Not equal")
			} else {
				d.err = nil
			}
		} else {
			// TODO
			d.err = nil
		}
	}
}

func (d *JsonDiffer) compareInOrder(result *graph.ExecutionResponse) bool {
	if d.Response.GetErrorCode() != result.GetErrorCode() {
		return false
	}
	// if d.Response.GetErrorMsg() != result.GetErrorMsg() {
	// 	return false
	// }

	if d.Response.GetSpaceName() != result.GetSpaceName() {
		return false
	}

	if len(d.Response.GetColumnNames()) != len(result.GetColumnNames()) {
		return false
	}
	for _, rc := range d.Response.GetColumnNames() {
		found := false
		for _, ec := range result.GetColumnNames() {
			if string(rc) == string(ec) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	if len(d.Response.GetRows()) != len(result.GetRows()) {
		return false
	}
	for i := range d.Response.GetRows() {
		if !d.compareRowValue(d.Response.GetRows()[i], result.GetRows()[i]) {
			return false
		}
	}

	return true
}

func (d *JsonDiffer) compareRowValue(l *graph.RowValue, r *graph.RowValue) bool {
	for _, lc := range l.GetColumns() {
		found := false
		for _, rc := range r.GetColumns() {
			if rc.String() == lc.String() {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
