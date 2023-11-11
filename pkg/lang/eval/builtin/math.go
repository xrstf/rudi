package builtin

import (
	"errors"
	"fmt"

	"go.xrstf.de/corel/pkg/lang/eval/coalescing"
)

func sumFunction(args []interface{}) (interface{}, error) {
	if len(args) < 1 {
		return nil, errors.New("(+ NUM+)")
	}

	sum := float64(0)
	for i, arg := range args {
		val, err := coalescing.ToFloat64(arg)
		if err != nil {
			return nil, fmt.Errorf("arg %d is not numeric: %w", i, err)
		}
		sum += val
	}

	return sum, nil
}

func minusFunction(args []interface{}) (interface{}, error) {
	if len(args) < 2 {
		return nil, errors.New("(- BASE NUM+)")
	}

	difference, err := coalescing.ToFloat64(args[0])
	if err != nil {
		return nil, fmt.Errorf("arg 0 is not numeric: %w", err)
	}

	for i, arg := range args[1:] {
		val, err := coalescing.ToFloat64(arg)
		if err != nil {
			return nil, fmt.Errorf("arg %d is not numeric: %w", i+1, err)
		}
		difference -= val
	}

	return difference, nil
}

func multiplyFunction(args []interface{}) (interface{}, error) {
	if len(args) < 2 {
		return nil, errors.New("(* NUM+)")
	}

	product := float64(1)
	for i, arg := range args {
		factor, err := coalescing.ToFloat64(arg)
		if err != nil {
			return nil, fmt.Errorf("arg %d is not numeric: %w", i, err)
		}
		product *= factor
	}

	return product, nil
}

func divideFunction(args []interface{}) (interface{}, error) {
	if len(args) < 2 {
		return nil, errors.New("(/ NUM+)")
	}

	result := float64(0)
	for i, arg := range args {
		divisor, err := coalescing.ToFloat64(arg)
		if err != nil {
			return nil, fmt.Errorf("arg %d is not numeric: %w", i, err)
		}
		if divisor == 0 {
			return nil, errors.New("division by zero")
		}
		result /= divisor
	}

	return result, nil
}
