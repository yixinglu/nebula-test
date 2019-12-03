package nebulatest

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/vesoft-inc/nebula-go/graph"
)

func TestJsonUnmarshalNebulaExecutionResponseBase64(t *testing.T) {
	jsonStr := `{"error_code": 0, "error_msg": "", "column_names": ["%s", "%s"], "space_name": "nba", "rows": [{"columns": [{"str": "%s"}, {"integer": 36}]}]}`
	jsonStr = fmt.Sprintf(jsonStr, base64.StdEncoding.EncodeToString([]byte("Teammate")), base64.StdEncoding.EncodeToString([]byte("Age")), base64.StdEncoding.EncodeToString([]byte("Tony Parker")))
	t.Log(jsonStr)
	var resp graph.ExecutionResponse
	if err := json.Unmarshal([]byte(jsonStr), &resp); err != nil {
		t.Fatal(err)
	} else {
		t.Log("Succeed to unmarshal json string to nebula execution response")
	}
}

func TestJsonUnmarshalNebulaExecutionResponse(t *testing.T) {
	// jsonStr := `{"error_code": 0, "error_msg": "", "column_names": ["Teammate", "Age"], "space_name": "nba", "rows": [{"columns": [{"str": "Tony Parker"}, {"integer": 36}]}]}`
	jsonStr := `{"error_code": 0, "error_msg": "", "column_names": [], "rows": [{"columns": [{"integer": 36}]}]}`
	var resp graph.ExecutionResponse
	if err := json.Unmarshal([]byte(jsonStr), &resp); err != nil {
		t.Fatal(err)
	} else {
		t.Log("Succeed to unmarshal json string to nebula execution response")
	}
}

func TestJsonUnmarshalNebulaExecutionResponse2(t *testing.T) {
	jsonStr := `{"error_code": 0, "error_msg": "", "column_names": ["Teammate", "Age"], "space_name": "nba", "rows": [{"columns": [{"str": "Tony Parker"}, {"integer": 36}]}]}`
	var resp executionResponse
	if err := json.Unmarshal([]byte(jsonStr), &resp); err != nil {
		t.Fatal(err)
	} else {
		t.Log("Succeed to unmarshal json string to nebula execution response")
	}
}

type Json struct {
	Name string `json:"name"`
}

func TestJsonUnmarshalGeneralStruct(t *testing.T) {
	jsonStr := `{"name": "hello"}`
	var j Json
	if err := json.Unmarshal([]byte(jsonStr), &j); err != nil {
		t.Fatal(err)
	}
}
