package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleString(t *testing.T) {
	tests := []struct {
		rvalue      string
		value       string
		parsedIndex int
	}{
		{
			rvalue:      "+PONG\r\n",
			value:       "PONG",
			parsedIndex: 5,
		},
	}

	for _, tt := range tests {
		result := SimpleStringValue{value: tt.value}.String()

		if result != tt.rvalue {
			t.Errorf("String() failed: expected %s got %s", tt.rvalue, result)
		}

		parseItem := SimpleStringValue{}
		repeatedValue := []byte(strings.Repeat(tt.rvalue, 2))
		read, err := parseItem.Parse(repeatedValue)
		if err != nil {
			t.Errorf("Parse(): %v", err)
		}

		if read != tt.parsedIndex {
			t.Errorf("Parse() failed: expected to read %d but read %d", tt.parsedIndex, read)
		}

		result = parseItem.Value().(string)

		if result != tt.value {
			t.Errorf("Parse() failed: expected %s got %s", tt.value, result)
		}
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		rvalue      string
		value       string
		parsedIndex int
	}{
		{
			rvalue:      "$4\r\nPONG\r\n",
			value:       "PONG",
			parsedIndex: 8,
		},
	}

	for _, tt := range tests {
		result := StringValue{value: tt.value}.String()

		if result != tt.rvalue {
			t.Errorf("String() failed: expected %s got %s", tt.rvalue, result)
		}

		parseItem := StringValue{}
		repeatedValue := []byte(strings.Repeat(tt.rvalue, 2))
		read, err := parseItem.Parse(repeatedValue)
		if err != nil {
			t.Errorf("Parse(): %v", err)
		}

		if read != tt.parsedIndex {
			t.Errorf("Parse() failed: expected to read %d but read %d", tt.parsedIndex, read)
		}

		result = parseItem.Value().(string)

		if result != tt.value {
			t.Errorf("Parse() failed: expected %s got %s", tt.value, result)
		}
	}
}

func TestError(t *testing.T) {
	t.SkipNow()
	assert.Equal(t, "-Error\r\n", ErrorValue{value: "Error"}.String())
}

func TestFormatInt(t *testing.T) {
	t.SkipNow()
	assert.Equal(t, ":4\r\n", IntValue{value: 4}.String())
}

func TestFormatArray(t *testing.T) {
	t.SkipNow()
	// expected := "*2\r\n$4\r\necho\r\n$3\r\nhey\r\n"
	// actual := []string{"echo", "hey"}

	// assert.Equal(t, expected, ArrayValue{values: actual})
}

func TestParseArray(t *testing.T) {
	t.SkipNow()
	// actual := []byte("*2\r\n$4\r\necho\r\n$3\r\nhey\r\n")
	// expected := []string{"echo", "hey"}

	// assert.Equal(t, expected, parseArray(actual))
}

func TestParseInt(t *testing.T) {
	t.SkipNow()
	actual := []byte(":4\r\n")
	expected := 4

	value := &IntValue{}
	_, err := value.Parse(actual)

	assert.NoError(t, err)
	assert.Equal(t, expected, value.Value())
}
