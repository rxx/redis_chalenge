package main

import (
	"fmt"
	"time"
)

func executeSetCommand(commands []string) RValue {
	const (
		reqArgCount  = 2
		reqArgWithPx = 4
		keyIndex     = 0
		valIndex     = 1
		durKey       = "px"
		durIndex     = 2
		durVal       = 3
	)

	if len(commands) < reqArgCount {
		return &ErrorValue{value: "Required key and value"}
	}

	key := commands[keyIndex]
	value := commands[valIndex]

	if len(commands) == reqArgWithPx && commands[durIndex] == durKey {
		duration, err := time.ParseDuration(commands[durVal] + "ms")
		if err != nil {
			return &ErrorValue{value: fmt.Errorf("can't parse px duration due to %w", err).Error()}
		}

		store.Set(key, value, duration)
		return &SimpleStringValue{value: "OK"}
	}

	store.Set(key, value, 0)
	return &SimpleStringValue{value: "OK"}
}

func executeGetCommand(commands []string) RValue {
	if len(commands) < 1 {
		return &ErrorValue{value: "Required key"}
	}

	key := commands[0]

	value, ok := store.Get(key)
	if !ok {
		return &StringValue{value: "nil"}
	}

	return &StringValue{value: value}
}

func executeInfoCommand(commands []string) RValue {
	var value string

	if len(commands) > 0 && commands[0] == "replication" {
		value = ReplicationInfoToString()
	}
	return &StringValue{value: value}
}
