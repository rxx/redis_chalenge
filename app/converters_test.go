package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestConverters(t *testing.T) {
	tests := []struct {
		rvalue string
		value  interface{}
		rType  string
	}{
		{
			rvalue: "+PONG\r\n",
			value:  "PONG",
			rType:  "SimpleString",
		},
		{
			rvalue: "$4\r\nPONG\r\n",
			value:  "PONG",
			rType:  "String",
		},
		{
			rvalue: "-Error\r\n",
			value:  "Error",
			rType:  "Error",
		},
		{
			rvalue: ":40\r\n",
			value:  40,
			rType:  "Int",
		},
		{
			rvalue: "*2\r\n$4\r\necho\r\n$3\r\nhey\r\n",
			value:  []string{"echo", "hey"},
			rType:  "Array",
		},
	}

	for _, tt := range tests {
		t.Run(tt.rType, func(t *testing.T) {
			if result := newRValueWithType(tt.rType, tt.value).String(); result != tt.rvalue {
				t.Errorf("String() failed: expected %s got %s", tt.rvalue, result)
			}

			parseItem := newRValueWithType(tt.rType, nil)
			// NOTE: We repeat value to check if parser will ignore any values over expected length
			repeatedValue := []byte(strings.Repeat(tt.rvalue, 2))
			if _, err := parseItem.Parse(repeatedValue); err != nil {
				t.Errorf("Parse(): %v", err)
			}

			if result := parseItem.Value(); !reflect.DeepEqual(result, tt.value) {
				t.Errorf("Parse() failed: expected %v got %v", tt.value, result)
			}
		})
	}
}
