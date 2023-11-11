package builtin

import (
	"errors"
	"fmt"
	"strings"

	"go.xrstf.de/corel/pkg/lang/eval/coalescing"
)

func concatFunction(args []interface{}) (interface{}, error) {
	if len(args) < 2 {
		return nil, errors.New("(concat GLUE LIST+)")
	}

	glue, err := coalescing.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("glue is not stringish: %w", err)
	}

	parts := []string{}
	for i, list := range args[1:] {
		items, ok := list.([]interface{})
		if !ok {
			part, err := coalescing.ToString(list)
			if err != nil {
				return nil, fmt.Errorf("arg %d is neither vector nor stringish: %w", i+1, err)
			}

			parts = append(parts, part)
			continue
		}

		for j, item := range items {
			part, err := coalescing.ToString(item)
			if err != nil {
				return nil, fmt.Errorf("element %d in arg %d is not stringish: %w", j, i+1, err)
			}
			parts = append(parts, part)
		}
	}

	return strings.Join(parts, glue), nil
}

func splitFunction(args []interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, errors.New("(split SEPARATOR STRING)")
	}

	sep, err := coalescing.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("separator is not stringish: %w", err)
	}

	source, err := coalescing.ToString(args[1])
	if err != nil {
		return nil, fmt.Errorf("source is not stringish: %w", err)
	}

	parts := strings.Split(source, sep)
	result := make([]interface{}, len(parts))
	for i, part := range parts {
		result[i] = part
	}

	return result, nil
}
