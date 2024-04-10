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
	data        []byte
	parsedIndex int
	err         error
}

func (v ParseError) Error() string {
	var errorStr strings.Builder
	errorStr.Write(v.data[0:v.parsedIndex])
	errorStr.WriteString("[error here]")
	errorStr.Write(v.data[v.parsedIndex:])

	return fmt.Sprintf("ParseError:\nexpr: %s\nerror: %v", errorStr.String(), v.err)
}

func (e ParseError) Unwrap() error {
	return e.err
}

type RValue interface {
	Value() interface{}
	String() string
	Bytes() []byte
	Parse(data []byte) (parsedIndex int, err error)
}

type SimpleStringValue struct {
	value string
}

// Example: "+PONG\r\n"
func (v SimpleStringValue) String() string {
	return fmt.Sprintf("%c%s%s", RSimpleString, v.value, CRLF)
}

func (v SimpleStringValue) Value() interface{} {
	return v.value
}

func (v SimpleStringValue) Bytes() []byte {
	return []byte(v.String())
}

func (v *SimpleStringValue) Parse(data []byte) (parsedIndex int, err error) {
	wrapError := func(err error) ParseError {
		return ParseError{data: data, parsedIndex: parsedIndex, err: err}
	}

	if data[0] != byte(RSimpleString) {
		err = wrapError(fmt.Errorf("SimpleStringValue.Parse: Missing \"%c\"", RSimpleString))
		return
	}

	str := string(data)
	parsedIndex = strings.Index(str, CRLF)

	if parsedIndex < 0 {
		err = wrapError(fmt.Errorf("SimpleStringValue.Parse: Missing \"%s\"", CRLF))
		return
	}

	v.value = str[1:parsedIndex]

	return
}

type StringValue struct {
	value string
}

// Example: "$4\r\nPONG\r\n"
func (v StringValue) String() string {
	return fmt.Sprintf("%c%d%s%s%s", RString, len(v.value), CRLF, v.value, CRLF)
}

func (v StringValue) Value() interface{} {
	return v.value
}

func (v StringValue) Bytes() []byte {
	return []byte(v.String())
}

func (v *StringValue) Parse(data []byte) (parsedIndex int, err error) {
	wrapError := func(err error) ParseError {
		return ParseError{data: data, parsedIndex: parsedIndex, err: err}
	}

	if data[0] != byte(RString) {
		err = wrapError(fmt.Errorf("StringValue.Parse: Missing \"%c\"", RString))
		return
	}

	str := string(data)
	parsedIndex = strings.Index(str, CRLF)

	if parsedIndex < 0 {
		err = wrapError(fmt.Errorf("StringValue.Parse: Missing \"%s\"", CRLF))
		return
	}

	length, err := strconv.Atoi(str[1:parsedIndex])
	if err != nil {
		err = wrapError(fmt.Errorf("StringValue.Parse: Error on reading length %w", err))
		return
	}

	parsedIndex += 2
	str = str[parsedIndex : parsedIndex+length]
	actualLen := len(str)

	if length != actualLen {
		err = wrapError(
			fmt.Errorf("StringValue.Parse: String %s has %d length but expected %d",
				str, actualLen, length))
		return
	}

	v.value = str

	return parsedIndex + length, nil
}

type ErrorValue struct {
	value string
}

func (v ErrorValue) String() string {
	return fmt.Sprintf("%c%s%s", RError, v.value, CRLF)
}

func (v ErrorValue) Value() interface{} {
	return v.value
}

func (v ErrorValue) Bytes() []byte {
	return []byte(v.String())
}

func (v *ErrorValue) Parse(data []byte) (parsedIndex int, err error) {
	wrapError := func(err error) ParseError {
		return ParseError{data: data, parsedIndex: parsedIndex, err: err}
	}

	if data[0] != byte(RError) {
		err = wrapError(fmt.Errorf("ErrorValue.Parse: Missing \"%c\"", RError))
		return
	}

	str := string(data)
	parsedIndex = strings.Index(str, CRLF)

	if parsedIndex < 0 {
		err = wrapError(fmt.Errorf("ErrorValue.Parse: Missing \"%s\"", CRLF))
		return
	}

	v.value = str[1:parsedIndex]

	return
}

type IntValue struct {
	value int
}

// Example: ":4\r\n"
func (v IntValue) String() string {
	return fmt.Sprintf("%c%d%s", RInt, v.value, CRLF)
}

func (v IntValue) Value() interface{} {
	return v.value
}

func (v IntValue) Bytes() []byte {
	return []byte(v.String())
}

func (v *IntValue) Parse(data []byte) (parsedIndex int, err error) {
	wrapError := func(err error) ParseError {
		return ParseError{data: data, parsedIndex: parsedIndex, err: err}
	}

	if data[0] != byte(RInt) {
		err = wrapError(fmt.Errorf("IntValue.Parse: Missing \"%c\"", RInt))
		return
	}

	str := string(data)
	parsedIndex = strings.Index(str, CRLF)

	if parsedIndex < 0 {
		err = wrapError(fmt.Errorf("IntValue.Parse: Missing \"%s\"", CRLF))
		return
	}

	digit, err := strconv.Atoi(str[1:parsedIndex])
	if err != nil {
		err = wrapError(fmt.Errorf("IntValue.Parse: Can't convert to due to %w", err))
		return
	}

	v.value = digit
	return
}

type ArrayValue struct {
	values []RValue
}

func (v ArrayValue) Value() interface{} {
	var result []string

	for _, item := range v.values {
		result = append(result, item.Value().(string))
	}

	return result
}

// Example: "*2\r\n$4\r\necho\r\n$3\r\nhey\r\n"
func (v ArrayValue) String() string {
	var result strings.Builder
	result.WriteRune(RArray)
	result.WriteString(strconv.Itoa(len(v.values)))
	result.WriteString(CRLF)

	for _, item := range v.values {
		result.WriteString(item.String())
	}

	return result.String()
}

func (v ArrayValue) Bytes() []byte {
	var result []byte

	for _, item := range v.values {
		result = append(result, item.Bytes()...)
	}

	return result
}

// Example:  *2\r\n$4\r\necho\r\n$3\r\nhey\r\n
// ["echo", "hey"]
func (v *ArrayValue) Parse(data []byte) (parsedIndex int, err error) {
	wrapError := func(err error) ParseError {
		return ParseError{data: data, parsedIndex: parsedIndex, err: err}
	}

	if data[0] != byte(RArray) {
		err = fmt.Errorf("ArrayValue.Parse: Missing \"%c\"", RArray)
		return
	}

	str := string(data)
	parsedIndex = strings.Index(str, CRLF)

	if parsedIndex < 0 {
		err = wrapError(fmt.Errorf("ArrayValue.Parse: Missing \"%s\" to parse array size", CRLF))
		return
	}

	size, err := strconv.Atoi(str[1:parsedIndex])
	if err != nil {
		err = wrapError(fmt.Errorf("ArrayValue.Parse: Error on reading size of an array due to %w", err))
		return
	}

	parsedIndex += 2

	var count int

	for parsedIndex < len(str) && size > 0 {
		item := NewRValue(data[parsedIndex])
		count, err = item.Parse(data[parsedIndex:])
		parsedIndex += count

		if err != nil {
			err = wrapError(err)
			return
		}

		v.values = append(v.values, item)
		parsedIndex += 2
		size--
	}

	return
}

func NewRValue(typeByte byte) RValue {
	switch typeByte {
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
		panic("NewRValue: type bite missing")
	}
}

func newRValueWithType(typeName string, value interface{}) RValue {
	switch typeName {
	case "SimpleString":
		if value == nil {
			value = ""
		}

		return &SimpleStringValue{value: value.(string)}
	case "Error":
		if value == nil {
			value = ""
		}

		return &ErrorValue{value: value.(string)}
	case "Int":
		if value == nil {
			value = 0
		}

		return &IntValue{value: value.(int)}
	case "String":
		if value == nil {
			value = ""
		}

		return &StringValue{value: value.(string)}
	case "Array":
		var values []RValue
		if value == nil {
			value = []string{}
		}

		for _, v := range value.([]string) {
			values = append(values, &StringValue{value: v})
		}

		return &ArrayValue{values: values}
	default:
		panic("newRValueWithType: Invalid type")
	}
}
