package nebulatest

import (
	"log"

	"github.com/vesoft-inc/nebula-go/graph"
)

type columnValue struct {
	BoolVal         *bool            `thrift:"bool_val,1" db:"bool_val" json:"bool_val,omitempty"`
	Integer         *int64           `thrift:"integer,2" db:"integer" json:"integer,omitempty"`
	Id              *graph.IdType    `thrift:"id,3" db:"id" json:"id,omitempty"`
	SinglePrecision *float32         `thrift:"single_precision,4" db:"single_precision" json:"single_precision,omitempty"`
	DoublePrecision *float64         `thrift:"double_precision,5" db:"double_precision" json:"double_precision,omitempty"`
	Str             *string          `thrift:"str,6" db:"str" json:"str,omitempty"`
	Timestamp       *graph.Timestamp `thrift:"timestamp,7" db:"timestamp" json:"timestamp,omitempty"`
	Year            *graph.Year      `thrift:"year,8" db:"year" json:"year,omitempty"`
	Month           *graph.YearMonth `thrift:"month,9" db:"month" json:"month,omitempty"`
	Date            *graph.Date      `thrift:"date,10" db:"date" json:"date,omitempty"`
	Datetime        *graph.DateTime  `thrift:"datetime,11" db:"datetime" json:"datetime,omitempty"`
}

type rowValue struct {
	Columns []*columnValue `thrift:"columns,1" db:"columns" json:"columns"`
}

type executionResponse struct {
	ErrorCode   graph.ErrorCode `thrift:"error_code,1,required" db:"error_code" json:"error_code"`
	LatencyInUs int32           `thrift:"latency_in_us,2,required" db:"latency_in_us" json:"latency_in_us"`
	ErrorMsg    *string         `thrift:"error_msg,3" db:"error_msg" json:"error_msg,omitempty"`
	ColumnNames []string        `thrift:"column_names,4" db:"column_names" json:"column_names,omitempty"`
	Rows        []*rowValue     `thrift:"rows,5" db:"rows" json:"rows,omitempty"`
	SpaceName   *string         `thrift:"space_name,6" db:"space_name" json:"space_name,omitempty"`
}

func (e *executionResponse) convertToNebulaResponse() *graph.ExecutionResponse {
	resp := graph.ExecutionResponse{
		ErrorCode:   e.ErrorCode,
		LatencyInUs: e.LatencyInUs,
		ErrorMsg:    e.ErrorMsg,
		SpaceName:   e.SpaceName,
	}
	if len(e.ColumnNames) > 0 {
		resp.ColumnNames = make([][]byte, len(e.ColumnNames))
		for i := range e.ColumnNames {
			resp.ColumnNames[i] = []byte(e.ColumnNames[i])
		}
	}

	if len(e.Rows) > 0 {
		resp.Rows = make([]*graph.RowValue, len(e.Rows))
		for i := range e.Rows {
			var row graph.RowValue
			row.Columns = make([]*graph.ColumnValue, len(e.Rows[i].Columns))
			for j := range e.Rows[i].Columns {
				c := e.Rows[i].Columns[j]
				var column graph.ColumnValue
				if c.Str != nil {
					column.Str = []byte(*c.Str)
				} else if c.BoolVal != nil {
					column.BoolVal = c.BoolVal
				} else if c.Integer != nil {
					column.Integer = c.Integer
				} else if c.Id != nil {
					column.Id = c.Id
				} else if c.SinglePrecision != nil {
					column.SinglePrecision = c.SinglePrecision
				} else if c.DoublePrecision != nil {
					column.DoublePrecision = c.DoublePrecision
				} else if c.Timestamp != nil {
					column.Timestamp = c.Timestamp
				} else if c.Year != nil {
					column.Year = c.Year
				} else if c.Month != nil {
					column.Month = c.Month
				} else if c.Date != nil {
					column.Date = c.Date
				} else if c.Datetime != nil {
					column.Datetime = c.Datetime
				} else {
					log.Println("Error column value")
				}
				row.Columns[j] = &column
			}
			resp.Rows[i] = &row
		}
	}
	return &resp
}
