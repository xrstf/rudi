package builtin

import (
	"errors"
	"fmt"
)

type GenericFunc func(args []interface{}) (interface{}, error)

var Functions = map[string]GenericFunc{
	"add": addFunction,
}

func addFunction(args []interface{}) (interface{}, error) {
	if len(args) < 1 {
		return nil, errors.New("(add NUM[ NUM]+)")
	}

	canInteger := true
	for _, arg := range args {
		if _, ok := arg.(int64); !ok {
			canInteger = false
			break
		}
	}

	if canInteger {
		sum := int64(0)

		for _, arg := range args {
			sum += arg.(int64)
		}

		return sum, nil
	}

	sum := float64(0)
	for i, arg := range args {
		switch val := arg.(type) {
		case int64:
			sum += float64(val)
		case float64:
			sum += val
		case nil:
			// NOP
		case bool:
			if val {
				sum += 1
			}
		default:
			return nil, fmt.Errorf("arg %d is not numeric, but %T", i, arg)
		}
	}

	return sum, nil
}
