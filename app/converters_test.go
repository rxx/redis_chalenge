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
			rvalue: "$0\r\n\r\n",
			value:  "",
			rType:  "String",
		},
		{
			rvalue: "$-1\r\n",
			value:  "nil",
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
		// {
		// 	rvalue: "*2\r\n$4\r\necho\r\n$3\r\nhey\r\n",
		// 	value:  []string{"echo", "hey"},
		// 	rType:  "Array",
		// },
		{
			rvalue: "*0\r\n",
			value:  []string{},
			rType:  "Array",
		},
	}

	for _, tt := range tests {
		t.Run(tt.rType, func(t *testing.T) {
			if result := newRValueWithType(tt.rType, tt.value).String(); result != tt.rvalue {
				t.Errorf("String() failed: expected %s got %s", tt.rvalue, result)
			}

			parseItem := newRValueWithType(tt.rType, nil)
			// We repeat value to check if parser will ignore any values over expected length
			repeatedValue := []byte(strings.Repeat(tt.rvalue, 2))
			if _, err := parseItem.Parse(repeatedValue); err != nil {
				t.Errorf("Parse(): %v", err)
			}
			result := parseItem.Value()

			var blankArray bool
			if reflect.TypeOf(tt.value) == reflect.TypeOf([]string{}) && len(tt.value.([]string)) == 0 {
				blankArray = true
			}

			if blankArray && len(result.([]string)) > 0 {
				t.Errorf("Parse() failed: expected blank array, got %v", result)
			}

			if !reflect.DeepEqual(result, tt.value) && !blankArray {
				t.Errorf("Parse() failed: expected %v got %v", tt.value, result)
			}
		})
	}
}
