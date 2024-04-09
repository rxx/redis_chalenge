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
	actual := []byte("*2\r\n$4\r\necho\r\n$3\r\nhey\r\n")
	expected := []string{"echo", "hey"}

	assert.Equal(t, expected, parseArray(actual))
}

func TestFormatEmptyArray(t *testing.T) {
	expected := "*0\r\n"
	actual := []string{}

	assert.Equal(t, expected, formatArray(actual))
}

func TestFormatNilArray(t *testing.T) {
	expected := "*-1\r\n"
	actual := nil

	assert.Equal(t, expected, formatArray(actual))
}

func TestParseBlankArray(t *testing.T) {
	actual := []byte("*0\r\n")
	expected := []string{}

	assert.Equal(t, expected, parseArray(actual))
}

func TestParseNilArray(t *testing.T) {
	actual := []byte("*-1\r\n")
	expected := nil

	assert.Equal(t, expected, parseArray(actual))
}

func TestParseArrayOfInts(t *testing.T) {
	actual := []byte("*3\r\n:1\r\n:2\r\n:3\r\n")
	expected := []int{1, 2, 3}

	assert.Equal(t, expected, formatArray(actual))
}

func TestParseMixedArray(t *testing.T) {
	actual := []byte("*5\r\n:1\r\n:2\r\n:3\r\n:4\r\n$5\r\nhello\r\n")
	expected := []int{1, 2, 3, 4, "hello"}

	assert.Equal(t, expected, formatArray(actual))
}

func TestParseMixedArray(t *testing.T) {
	actual := []byte("*5\r\n:1\r\n:2\r\n:3\r\n:4\r\n$5\r\nhello\r\n")
	expected := []int{1, 2, 3, 4, "hello"}

	assert.Equal(t, expected, formatArray(actual))
}

func TestParseMixedNestedArray(t *testing.T) {
	actual := []byte("*2\r\n*3\r\n:1\r\n:2\r\n:3\r\n*2\r\n+Hello\r\n$5\r\nWorld\r\n")
	expected := []int{[]int{1, 2, 3}, []string{"Hello", "World"}}

	assert.Equal(t, expected, formatArray(actual))
}

func TestParseString(t *testing.T) {
	actual := []byte("$4\r\necho\r\n")
	expected := "echo"

	assert.Equal(t, expected, parseString(actual))
}

func TestFormatEmptyString(t *testing.T) {
	assert.Equal(t, "$0\r\n\r\n", formatString(""))
}

func TestFormatNilString(t *testing.T) {
	assert.Equal(t, "$-1\r\n", formatString(nil))
}

func TestParseEmptyString(t *testing.T) {
	actual := []byte("$0\r\n\r\n")
	expected := ""

	assert.Equal(t, expected, parseString(actual))
}

func TestParseNilString(t *testing.T) {
	actual := []byte("$-1\r\n")
	expected := nil

	assert.Equal(t, expected, parseString(actual))
}

func TestParseInt(t *testing.T) {
	actual := []byte(":4\r\n")
	expected := 4

	assert.Equal(t, expected, parseInt(actual))
}
