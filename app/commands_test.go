package main

import "testing"

func TestExecuteGetCommand(t *testing.T) {
	tests := []struct {
		req []string
		res RValue
	}{
		{
			req: []string{},
			res: &ErrorValue{value: "Required key"},
		},
		{
			req: []string{"invalid"},
			res: &StringValue{value: "nil"},
		},
	}

	for _, tt := range tests {
		res := executeGetCommand(tt.req)

		if res.String() != tt.res.String() {
			t.Errorf("Expected %q, got %q", tt.res, res)
		}
	}
}

func TestExecuteSetCommand(t *testing.T) {
	tests := []struct {
		req []string
		res RValue
	}{
		{
			req: []string{},
			res: &ErrorValue{value: "Required key and value"},
		},
		{
			req: []string{"set_unique_key"},
			res: &ErrorValue{value: "Required key and value"},
		},
		{
			req: []string{"set_unique_key", "value"},
			res: &SimpleStringValue{value: "OK"},
		},
		{
			req: []string{"set_unique_key", "value", "px", "100"},
			res: &SimpleStringValue{value: "OK"},
		},
	}

	for _, tt := range tests {
		res := executeSetCommand(tt.req)

		if res.String() != tt.res.String() {
			t.Errorf("Expected %q, got %q", tt.res, res)
		}
	}
}
