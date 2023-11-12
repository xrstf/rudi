package builtin

import (
	"errors"
	"fmt"

	"go.xrstf.de/corel/pkg/lang/eval/coalescing"
)

func andFunction(args []interface{}) (interface{}, error) {
	if len(args) < 1 {
		return nil, errors.New("(and CONDITION+)")
	}

	result := true
	for i, item := range args {
		part, err := coalescing.ToBool(item)
		if err != nil {
			return nil, fmt.Errorf("arg %d is nor boolish: %w", i, err)
		}

		result = result && part
	}

	return result, nil
}

func orFunction(args []interface{}) (interface{}, error) {
	if len(args) < 1 {
		return nil, errors.New("(or CONDITION+)")
	}

	result := false
	for i, item := range args {
		part, err := coalescing.ToBool(item)
		if err != nil {
			return nil, fmt.Errorf("arg %d is nor boolish: %w", i, err)
		}

		result = result || part
	}

	return result, nil
}
