// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"errors"
	"fmt"

	"go.xrstf.de/corel/pkg/lang/eval/coalescing"
	"go.xrstf.de/corel/pkg/lang/eval/types"
)

func int64sOnly(vals []any) bool {
	for _, val := range vals {
		if !coalescing.Int64Compatible(val) {
			return false
		}
	}

	return true
}

func sumFunction(ctx types.Context, args []Argument) (any, error) {
	if size := len(args); size < 2 {
		return nil, fmt.Errorf("expected 2+ arguments, got %d", size)
	}

	values, err := evalArgs(ctx, args, 0)
	if err != nil {
		return nil, err
	}

	if int64sOnly(values) {
		sum := int64(0)
		for _, arg := range values {
			val, _ := coalescing.ToInt64(arg)
			sum += val
		}

		return sum, nil
	}

	sum := float64(0)
	for i, arg := range values {
		val, err := coalescing.ToFloat64(arg)
		if err != nil {
			return nil, fmt.Errorf("arg %d is not numeric: %w", i, err)
		}
		sum += val
	}

	return sum, nil
}

func minusFunction(ctx types.Context, args []Argument) (any, error) {
	if size := len(args); size < 2 {
		return nil, fmt.Errorf("expected 2+ arguments, got %d", size)
	}

	values, err := evalArgs(ctx, args, 0)
	if err != nil {
		return nil, err
	}

	if int64sOnly(values) {
		difference, _ := coalescing.ToInt64(values[0])

		for _, arg := range values[1:] {
			val, _ := coalescing.ToInt64(arg)
			difference -= val
		}

		return difference, nil
	}

	difference, err := coalescing.ToFloat64(values[0])
	if err != nil {
		return nil, fmt.Errorf("argument #0 is not numeric: %w", err)
	}

	for i, arg := range values[1:] {
		val, err := coalescing.ToFloat64(arg)
		if err != nil {
			return nil, fmt.Errorf("argument #%d is not numeric: %w", i+1, err)
		}
		difference -= val
	}

	return difference, nil
}

func multiplyFunction(ctx types.Context, args []Argument) (any, error) {
	if size := len(args); size < 2 {
		return nil, fmt.Errorf("expected 2+ arguments, got %d", size)
	}

	values, err := evalArgs(ctx, args, 0)
	if err != nil {
		return nil, err
	}

	if int64sOnly(values) {
		product := int64(1)

		for _, arg := range values {
			factor, _ := coalescing.ToInt64(arg)
			product *= factor
		}

		return product, nil
	}

	product := float64(1)
	for i, arg := range values {
		factor, err := coalescing.ToFloat64(arg)
		if err != nil {
			return nil, fmt.Errorf("arg %d is not numeric: %w", i, err)
		}
		product *= factor
	}

	return product, nil
}

func divideFunction(ctx types.Context, args []Argument) (any, error) {
	if size := len(args); size < 2 {
		return nil, fmt.Errorf("expected 2+ arguments, got %d", size)
	}

	values, err := evalArgs(ctx, args, 0)
	if err != nil {
		return nil, err
	}

	result, err := coalescing.ToFloat64(values[0])
	if err != nil {
		return nil, fmt.Errorf("argument #0 is not numeric: %w", err)
	}

	for i, arg := range args[1:] {
		divisor, err := coalescing.ToFloat64(arg)
		if err != nil {
			return nil, fmt.Errorf("argument #%d is not numeric: %w", i+1, err)
		}
		if divisor == 0 {
			return nil, errors.New("division by zero")
		}
		result /= divisor
	}

	return result, nil
}
