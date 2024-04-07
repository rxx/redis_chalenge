package main

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	RSimpleString = '+'
	RError        = '-'
	RInt          = ':'
	RString       = '$'
	RArray        = '*'
	CRLF          = "\r\n"
)

// Example: "*2\r\n$4\r\necho\r\n$3\r\nhey\r\n"
func formatArray(array []string) string {
	var result strings.Builder
	result.WriteRune(RArray)
	result.WriteString(strconv.Itoa(len(array)))
	result.WriteString(CRLF)

	for _, str := range array {
		result.WriteString(formatString(str))
	}

	return result.String()
}

func parseArray(data []byte) []string {
	// echo command
	// *2\r\n$4\r\necho\r\n$3\r\nhey\r\n
	_ = data
	return nil
}

// Example: "+PONG\r\n",
func formatSimpleString(str string) string {
	return fmt.Sprintf("%c%s%s", RSimpleString, str, CRLF)
}

// Example: "-Error\r\n",
func formatError(str string) string {
	return fmt.Sprintf("%c%s%s", RError, str, CRLF)
}

// Example: "$4\r\nPONG\r\n",
func formatString(str string) string {
	return fmt.Sprintf("%c%d%s%s%s", RString, len(str), CRLF, str, CRLF)
}

func parseString(data []byte) string {
	_ = data
	return ""
}

// Example: ":4\r\n"
func formatInt(digit int) string {
	return fmt.Sprintf(":%d%s", digit, CRLF)
}

func parseInt(data []byte) int {
	var index int

	if data[0] == byte(RInt) {
		index = 1
	} else {
		index = 0
	}
	digit, err := strconv.Atoi(strings.TrimRight(string(data[index:]), CRLF))
	if err != nil {
		fmt.Println("Error on converting int", err)
		return -1
	}

	return digit
}
