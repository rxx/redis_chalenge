package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatSimpleString(t *testing.T) {
	assert.Equal(t, "+PONG\r\n", formatSimpleString("PONG"))
}

func TestFormatString(t *testing.T) {
	assert.Equal(t, "$4\r\nPONG\r\n", formatString("PONG"))
}

func TestError(t *testing.T) {
	assert.Equal(t, "-Error\r\n", formatError("Error"))
}

func TestFormatInt(t *testing.T) {
	assert.Equal(t, ":4\r\n", formatInt(4))
}

func TestFormatArray(t *testing.T) {
	expected := "*2\r\n$4\r\necho\r\n$3\r\nhey\r\n"
	actual := []string{"echo", "hey"}

	assert.Equal(t, expected, formatArray(actual))
}

func TestParseArray(t *testing.T) {
	t.SkipNow()
	actual := []byte("*2\r\n$4\r\necho\r\n$3\r\nhey\r\n")
	expected := []string{"echo", "hey"}

	assert.Equal(t, expected, parseArray(actual))
}

func TestParseString(t *testing.T) {
	t.SkipNow()
	actual := []byte("$4\r\necho\r\n")
	expected := "echo"

	assert.Equal(t, expected, parseString(actual))
}

func TestParseInt(t *testing.T) {
	actual := []byte(":4\r\n")
	expected := 4

	assert.Equal(t, expected, parseInt(actual))
}
