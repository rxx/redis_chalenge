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

// Example:  *2\r\n$4\r\necho\r\n$3\r\nhey\r\n
// ["echo", "hey"]
func parseArray(data []byte) []string {
	if data[0] != byte(RArray) {
		fmt.Printf("Expected %c to parse array", RArray)
		return nil
	}

	currentIndex := 1
	for currentIndex < len(data) && data[currentIndex+1] != byte('\r') {
		currentIndex++
	}

	if currentIndex == len(data) {
		fmt.Println("Invalid array, can not find size")
		return nil
	}

	size, err := strconv.Atoi(string(data[1:currentIndex]))
	if err != nil {
		fmt.Println("Error on reading size of an array", err.Error())
		return nil
	}

	result := make([]string, size)
	currentIndex = currentIndex + 3

	for currentIndex < len(data) {
		switch data[currentIndex] {
		case RInt:
		}
	}

	return result
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
	if data[0] != byte(RString) {
		fmt.Printf("Expected %c to parse bulk string", RString)
		return ""
	}

	parts := strings.Split(string(data[1:]), CRLF)
	length, err := strconv.Atoi(parts[0])
	if err != nil {
		fmt.Println("Error on reading String length", err.Error())
		return ""
	}

	str := parts[1]
	if length != len(str) {
		fmt.Printf("String %s has %d length but expected %d", str, len(str), length)
		return ""
	}

	return str
}

// Example: ":4\r\n"
func formatInt(digit int) string {
	return fmt.Sprintf(":%d%s", digit, CRLF)
}

func parseInt(data []byte) int {
	if data[0] != byte(RInt) {
		fmt.Printf("Expected %c to parse integer", RInt)
		return -1
	}

	digit, err := strconv.Atoi(strings.TrimRight(string(data[1:]), CRLF))
	if err != nil {
		fmt.Println("Error on converting int", err)
		return -1
	}

	return digit
}
