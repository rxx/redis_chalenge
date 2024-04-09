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

type ParseError struct {
	Expr        string
	ParsedIndex int
	Msg         string
}

func (e ParseError) Error() string {
	errorStr := e.Expr[0:e.ParsedIndex] + "[error here]" + e.Expr[e.ParsedIndex:]
	return fmt.Sprintf("ParseError: %s - expr: %q", e.Msg, errorStr)
}

type RValue interface {
	String() string
	Bytes() []byte
	Parse(data []byte) (parsedIndex int, err error)
}

type SimpleStringValue struct {
	Value string
}

// Example: "+PONG\r\n"
func (v SimpleStringValue) String() string {
	return fmt.Sprintf("%c%s%s", RSimpleString, v.Value, CRLF)
}

func (v SimpleStringValue) Bytes() []byte {
	return []byte(v.String())
}

func (v *SimpleStringValue) Parse(data []byte) (parsedIndex int, err error) {
	if data[0] != byte(RSimpleString) {
		err = fmt.Errorf("Wrong data type. Missing \"%c\" to parse simple string", RSimpleString)
		return 0, err
	}

	str := string(data)
	parsedIndex = strings.Index(str, CRLF)

	if parsedIndex < 0 {
		err = fmt.Errorf("Wrong data type. Missing \"%s\" to parse simple string", CRLF)
		return 0, err
	}

	v.Value = str[1:parsedIndex]

	return
}

type StringValue struct {
	Value string
}

// Example: "$4\r\nPONG\r\n"
func (v StringValue) String() string {
	return fmt.Sprintf("%c%d%s%s%s", RString, len(v.Value), CRLF, v.Value, CRLF)
}

func (v StringValue) Bytes() []byte {
	return []byte(v.String())
}

func (v *StringValue) Parse(data []byte) (parsedIndex int, err error) {
	if data[0] != byte(RString) {
		err = fmt.Errorf("Wrong data type. Missing \"%c\" to parse string", RString)
		return 0, err
	}

	str := string(data)
	parsedIndex = strings.Index(str, CRLF)

	if parsedIndex < 0 {
		err = fmt.Errorf("Wrong data type. Missing \"%s\" to parse string length", CRLF)
		return 0, err
	}

	length, err := strconv.Atoi(str[1:parsedIndex])
	if err != nil {
		err = fmt.Errorf("Error on reading String length", err.Error())
		return 0, err
	}

	parsedIndex += 2
	str = str[parsedIndex:length]
	actualLen := len(str)
	if length != actualLen {
		err = fmt.Errorf("String %s has %d length but expected %d", str, actualLen, length)
		return 0, err
	}

	v.Value = str
	return
}

type ErrorValue struct {
	Value string
}

func (v ErrorValue) String() string {
	return fmt.Sprintf("%c%s%s", RError, v.Value, CRLF)
}

func (v ErrorValue) Bytes() []byte {
	return []byte(v.String())
}

func (v *ErrorValue) Parse(data []byte) (parsedIndex int, err error) {
	if data[0] != byte(RError) {
		err = fmt.Errorf("Wrong data type. Missing \"%c\" to parse error", RError)
		return 0, err
	}

	str := string(data)
	parsedIndex = strings.Index(str, CRLF)

	if parsedIndex < 0 {
		err = fmt.Errorf("Wrong data type. Missing \"%s\" to parse error", CRLF)
		return 0, err
	}

	v.Value = str[1:parsedIndex]
	return
}

type IntValue struct {
	Value int
}

// Example: ":4\r\n"
func (v IntValue) String() string {
	return fmt.Sprintf(":%d%s", v.Value, CRLF)
}

func (v IntValue) Bytes() []byte {
	return []byte(v.String())
}

func (v *IntValue) Parse(data []byte) (parsedIndex int, err error) {
	if data[0] != byte(RInt) {
		err = fmt.Errorf("Wrong data type. Missing \"%c\" to parse integer", RInt)
		return 0, err
	}

	str := string(data)
	parsedIndex = strings.Index(str, CRLF)

	if parsedIndex < 0 {
		err = fmt.Errorf("Wrong data type. Missing \"%s\" to parse integer", CRLF)
		return 0, err
	}

	digit, err := strconv.Atoi(str[1:parsedIndex])
	if err != nil {
		err = fmt.Errorf("Error on converting int", err.Error())
		return 0, err
	}

	v.Value = digit
	return
}

type ArrayValue struct {
	Value []RValue
}

// Example: "*2\r\n$4\r\necho\r\n$3\r\nhey\r\n"
func (v ArrayValue) String() string {
	var result strings.Builder
	result.WriteRune(RArray)
	result.WriteString(strconv.Itoa(len(v.Value)))
	result.WriteString(CRLF)

	for _, item := range v.Value {
		result.WriteString(item.String())
	}

	return result.String()
}

func (v ArrayValue) Bytes() []byte {
	var result []byte

	for _, item := range v.Value {
		result = append(result, item.Bytes()...)
	}

	return result
}

// Example:  *2\r\n$4\r\necho\r\n$3\r\nhey\r\n
// ["echo", "hey"]
func (v *ArrayValue) Parse(data []byte) (parsedIndex int, err error) {
	if data[0] != byte(RArray) {
		err = fmt.Errorf("Wrong data type. Missing \"%c\" to parse array", RArray)
		return 0, err
	}

	str := string(data)
	parsedIndex = strings.Index(str, CRLF)

	if parsedIndex < 0 {
		err = fmt.Errorf("Wrong data type. Missing \"%s\" to parse array size", CRLF)
		return 0, err
	}

	size, err := strconv.Atoi(str[1:parsedIndex])
	if err != nil {
		err = fmt.Errorf("Error on reading size of an array", err.Error())
		return 0, err
	}

	v.Value = make([]RValue, size)
	parsedIndex += 2

	for parsedIndex < len(str) {
		item := NewRValue(data[parsedIndex])
		count, err := item.Parse(data[parsedIndex:])
		if err != nil {
			return parsedIndex + count, err
		}

		v.Value = append(v.Value, item)
		parsedIndex += count + 2
	}

	return
}

func NewRValue(rType byte) RValue {
	switch rType {
	case RSimpleString:
		return &SimpleStringValue{}
	case RError:
		return &ErrorValue{}
	case RInt:
		return &IntValue{}
	case RString:
		return &StringValue{}
	case RArray:
		return &ArrayValue{}
	default:
		panic("Invalid type")
	}
}
