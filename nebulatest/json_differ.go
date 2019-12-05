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
	// result = fmt.Sprintf("%q", result)
	var resp executionResponse
	if err := json.Unmarshal([]byte(result), &resp); err != nil {
		d.err = fmt.Errorf("Fail to parse JSON string, error: %s", err.Error())
	} else {
		r := resp.convertToNebulaResponse()
		if err = d.compare(r); err != nil {
			d.err = err
		} else {
			d.err = nil
		}
	}
}

func (d *JsonDiffer) compare(result *graph.ExecutionResponse) error {
	if d.Response.GetErrorCode() != result.GetErrorCode() {
		return fmt.Errorf("ErrorCode: %v vs. %v", d.Response.GetErrorCode(), result.GetErrorCode())
	}
	if result.IsSetErrorMsg() && d.Response.GetErrorMsg() != result.GetErrorMsg() {
		return fmt.Errorf("ErrorMsg: %s vs. %s", d.Response.GetErrorMsg(), result.GetErrorMsg())
	}

	if result.IsSetSpaceName() && d.Response.GetSpaceName() != result.GetSpaceName() {
		return fmt.Errorf("SpaceName: %s vs. %s", d.Response.GetSpaceName(), result.GetSpaceName())
	}

	if result.IsSetColumnNames() {
		if len(d.Response.GetColumnNames()) != len(result.GetColumnNames()) {
			return fmt.Errorf("Length of column names: %d vs. %d", len(d.Response.GetColumnNames()), len(result.GetColumnNames()))
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
				return fmt.Errorf("NotFoundColumnName: %s", string(rc))
			}
		}
	}

	if result.IsSetRows() {
		if len(d.Response.GetRows()) != len(result.GetRows()) {
			return fmt.Errorf("Number of rows: %d vs. %d", d.Response.GetRows(), result.GetRows())
		}

		if d.Order {
			for i := range d.Response.GetRows() {
				if !d.compareRowValue(d.Response.GetRows()[i], result.GetRows()[i]) {
					return fmt.Errorf("Rows: %s vs. %s", d.Response.GetRows()[i].String(), result.GetRows()[i].String())
				}
			}
		} else {
			for _, i := range d.Response.GetRows() {
				found := false
				for _, j := range result.GetRows() {
					if d.compareRowValue(i, j) {
						found = true
						break
					}
				}
				if !found {
					return fmt.Errorf("NotFoundRow: %s", i)
				}
			}
		}
	}

	return nil
}

func (d *JsonDiffer) compareRowValue(l *graph.RowValue, r *graph.RowValue) bool {
	for _, lc := range l.GetColumns() {
		found := false
		for _, rc := range r.GetColumns() {
			if d.compareColumnValue(lc, rc) {
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

func (d *JsonDiffer) compareColumnValue(l *graph.ColumnValue, r *graph.ColumnValue) bool {
	if l.IsSetBoolVal() && r.IsSetBoolVal() {
		return l.GetBoolVal() == r.GetBoolVal()
	} else if l.IsSetInteger() && r.IsSetInteger() {
		return l.GetInteger() == r.GetInteger()
	} else if l.IsSetId() && r.IsSetId() {
		return l.GetId() == r.GetId()
	} else if l.IsSetStr() && r.IsSetStr() {
		return string(l.GetStr()) == string(r.GetStr())
	} else if l.IsSetDate() && r.IsSetDate() {
		return l.GetDate().String() == r.GetDate().String()
	} else if l.IsSetDatetime() && r.IsSetDatetime() {
		return l.GetDatetime().String() == r.GetDatetime().String()
	} else if l.IsSetTimestamp() && r.IsSetTimestamp() {
		return l.GetTimestamp() == r.GetTimestamp()
	} else if l.IsSetSinglePrecision() && r.IsSetSinglePrecision() {
		return l.GetSinglePrecision() == r.GetSinglePrecision()
	} else if l.IsSetDoublePrecision() && r.IsSetDoublePrecision() {
		return l.GetDoublePrecision() == r.GetDoublePrecision()
	} else {
		return false
	}
}
