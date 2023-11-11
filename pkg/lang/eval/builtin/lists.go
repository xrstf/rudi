package builtin

import (
	"errors"
)

func lenFunction(args []interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, errors.New("(len LIST)")
	}

	list, ok := args[0].([]interface{})
	if !ok {
		return nil, errors.New("argument is not a vector")
	}

	return len(list), nil
}
