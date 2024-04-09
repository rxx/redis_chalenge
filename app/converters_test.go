package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatSimpleString(t *testing.T) {
	assert.Equal(t, "+PONG\r\n", SimpleStringValue{value: "PONG"}.String())
}

func TestFormatString(t *testing.T) {
	assert.Equal(t, "$4\r\nPONG\r\n", StringValue{value: "PONG"}.String())
}

func TestError(t *testing.T) {
	assert.Equal(t, "-Error\r\n", ErrorValue{value: "Error"}.String())
}

func TestFormatInt(t *testing.T) {
	assert.Equal(t, ":4\r\n", IntValue{value: 4}.String())
}

func formatArray(arr []string) string {
	_ = arr
	return ""
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

func TestParseString(t *testing.T) {
	actual := []byte("$4\r\npiNg\r\n")
	expected := "ping"

	stringValue := &StringValue{}
	_, err := stringValue.Parse(actual)

	assert.NoError(t, err)
	assert.Equal(t, expected, stringValue.Value())
}

func TestParseInt(t *testing.T) {
	actual := []byte(":4\r\n")
	expected := 4

	value := &IntValue{}
	_, err := value.Parse(actual)

	assert.NoError(t, err)
	assert.Equal(t, expected, value.Value())
}
