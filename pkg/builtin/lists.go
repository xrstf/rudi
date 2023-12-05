// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"fmt"
	"strings"

	"go.xrstf.de/rudi/pkg/equality"
	"go.xrstf.de/rudi/pkg/eval"
	"go.xrstf.de/rudi/pkg/eval/types"
	"go.xrstf.de/rudi/pkg/lang/ast"
)

func vectorLenFunction(vec []any) (any, error) {
	return len(vec), nil
}

func objectLenFunction(obj map[string]any) (any, error) {
	return len(obj), nil
}

func stringLenFunction(s string) (any, error) {
	return len(s), nil
}

func appendToVectorFunction(base []any, args ...any) (any, error) {
	result := []any{}
	result = append(result, base...)
	result = append(result, args...)

	return result, nil
}

func appendToStringFunction(base string, args ...string) (any, error) {
	return base + strings.Join(args, ""), nil
}

func prependToVectorFunction(base []any, args ...any) (any, error) {
	return append(args, base...), nil
}

func prependToStringFunction(base string, args ...string) (any, error) {
	return strings.Join(args, "") + base, nil
}

func reverseStringFunction(s string) (any, error) {
	// thank you https://stackoverflow.com/a/10030772
	result := []rune(s)
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return string(result), nil
}

func reverseVectorFunction(vec []any) (any, error) {
	// clone original data
	result := append([]any{}, vec...)

	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return result, nil
}

// (range VECTOR [item] expr+)
// (range VECTOR [i item] expr+)
// (range OBJECT [val] expr+)
// (range OBJECT [key val] expr+)
func rangeFunction(ctx types.Context, args []ast.Expression) (any, error) {
	// decode desired loop variable namings, as that's cheap to do
	loopIndexName, loopVarName, err := evalNamingVector(ctx, args[1])
	if err != nil {
		return nil, fmt.Errorf("argument #1: not a valid naming vector: %w", err)
	}

	// evaluate source list
	var source any

	innerCtx := ctx

	innerCtx, source, err = eval.EvalExpression(innerCtx, args[0])
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate source: %w", err)
	}

	var result any

	// list over vector elements
	if sourceVector, err := ctx.Coalesce().ToVector(source); err == nil {
		for i, item := range sourceVector {
			// do not use separate contexts for each loop iteration, as the loop might build up a counter
			innerCtx = innerCtx.WithVariable(loopVarName, item)
			if loopIndexName != "" {
				innerCtx = innerCtx.WithVariable(loopIndexName, i)
			}

			for _, expr := range args[2:] {
				innerCtx, result, err = eval.EvalExpression(innerCtx, expr)
				if err != nil {
					return nil, err
				}
			}
		}

		return result, nil
	}

	// list over object elements
	if sourceObject, err := ctx.Coalesce().ToObject(source); err == nil {
		for key, value := range sourceObject {
			// do not use separate contexts for each loop iteration, as the loop might build up a counter
			innerCtx = innerCtx.WithVariable(loopVarName, value)
			if loopIndexName != "" {
				innerCtx = innerCtx.WithVariable(loopIndexName, key)
			}

			for _, expr := range args[2:] {
				innerCtx, result, err = eval.EvalExpression(innerCtx, expr)
				if err != nil {
					return nil, err
				}
			}
		}

		return result, nil
	}

	return nil, fmt.Errorf("cannot range over %T", source)
}

// (map VECTOR identifier)
// (map VECTOR [item] expr+)
// (map VECTOR [i item] expr+)
// (map OBJECT identifier)
// (map OBJECT [item] expr+)
// (map OBJECT [i item] expr+)
func mapFunction(ctx types.Context, args []ast.Expression) (any, error) {
	// evaluate the first argument;
	// (map (map .foo +) stuff) should work, so the first argument only needs to _evaluate_
	// to a vector/object, it doesn't need to be a literal objectnode/vectornode.
	innerCtx, source, err := eval.EvalExpression(ctx, args[0])
	if err != nil {
		return nil, fmt.Errorf("argument #0: %w", err)
	}

	if err := checkIterable(ctx, source); err != nil {
		return nil, fmt.Errorf("argument #0: %w", err)
	}

	// handle plain function calls
	// (map VECTOR identifier)
	// (map OBJECT identifier)
	if len(args) == 2 {
		return anonymousMapFunction(innerCtx, source, args[1])
	}

	// all further forms are (map THING NAMING_VEC EXPR+)

	// decode desired loop variable namings
	indexVarName, valueVarName, err := evalNamingVector(innerCtx, args[1])
	if err != nil {
		return nil, fmt.Errorf("argument #1: not a valid naming vector: %w", err)
	}

	mapItem := func(ctx types.Context) (types.Context, any, error) {
		var result any

		for _, expr := range args[2:] {
			ctx, result, err = eval.EvalExpression(ctx, expr)
			if err != nil {
				return ctx, nil, err
			}
		}

		return ctx, result, nil
	}

	// list over vector elements
	if sourceVector, err := ctx.Coalesce().ToVector(source); err == nil {
		output := make([]any, len(sourceVector))

		for i, item := range sourceVector {
			// do not use separate contexts for each loop iteration, as the loop might build up a counter
			innerCtx = innerCtx.WithVariable(valueVarName, item)
			if indexVarName != "" {
				innerCtx = innerCtx.WithVariable(indexVarName, i)
			}

			var result any
			innerCtx, result, err = mapItem(innerCtx)
			if err != nil {
				return nil, err
			}

			output[i] = result
		}

		return output, nil
	}

	// list over object elements
	if sourceObject, err := ctx.Coalesce().ToObject(source); err == nil {
		output := map[string]any{}

		for key, value := range sourceObject {
			// do not use separate contexts for each loop iteration, as the loop might build up a counter
			innerCtx = innerCtx.WithVariable(valueVarName, value)
			if indexVarName != "" {
				innerCtx = innerCtx.WithVariable(indexVarName, key)
			}

			var result any
			innerCtx, result, err = mapItem(innerCtx)
			if err != nil {
				return nil, err
			}

			output[key] = result
		}

		return output, nil
	}

	return nil, fmt.Errorf("cannot map %T", source)
}

// (map VECTOR identifier)
// (map OBJECT identifier)
func anonymousMapFunction(ctx types.Context, source any, expr ast.Expression) (any, error) {
	identifier, ok := expr.(ast.Identifier)
	if !ok {
		return nil, fmt.Errorf("argument #1: expected identifier, got %T", expr)
	}

	// call the function
	innerCtx := ctx

	mapItem := func(ctx types.Context, item any) (types.Context, any, error) {
		return eval.EvalFunctionCall(ctx, identifier, []ast.Expression{types.MakeShim(item)})
	}

	if vector, err := ctx.Coalesce().ToVector(source); err == nil {
		output := make([]any, len(vector))

		for i, item := range vector {
			var (
				result any
				err    error
			)

			innerCtx, result, err = mapItem(innerCtx, item)
			if err != nil {
				return nil, err
			}

			output[i] = result
		}

		return output, nil
	}

	if object, err := ctx.Coalesce().ToObject(source); err == nil {
		output := map[string]any{}

		for key, value := range object {
			var (
				result any
				err    error
			)

			innerCtx, result, err = mapItem(innerCtx, value)
			if err != nil {
				return nil, err
			}

			output[key] = result
		}

		return output, nil
	}

	// should never happen, as this function call is already gated by a type check
	return nil, fmt.Errorf("cannot apply map to %T", source)
}

// (filter VECTOR identifier)
// (filter VECTOR [item] expr+)
// (filter VECTOR [i item] expr+)
// (filter OBJECT identifier)
// (filter OBJECT [val] expr+)
// (filter OBJECT [key val] expr+)
func filterFunction(ctx types.Context, args []ast.Expression) (any, error) {
	// evaluate the first argument;
	// (filter (filter .foo +) stuff) should work, so the first argument only needs to _evaluate_
	// to a vector/object, it doesn't need to be a literal objectnode/vectornode.
	innerCtx, source, err := eval.EvalExpression(ctx, args[0])
	if err != nil {
		return nil, fmt.Errorf("argument #0: %w", err)
	}

	if err := checkIterable(ctx, source); err != nil {
		return nil, fmt.Errorf("argument #0: %w", err)
	}

	// handle plain function calls
	// (filter VECTOR identifier)
	// (filter OBJECT identifier)
	if len(args) == 2 {
		return anonymousFilterFunction(innerCtx, source, args[1])
	}

	// all further forms are (map THING NAMING_VEC EXPR+)

	// decode desired loop variable namings
	indexVarName, valueVarName, err := evalNamingVector(innerCtx, args[1])
	if err != nil {
		return nil, fmt.Errorf("argument #1: not a valid naming vector: %w", err)
	}

	filter := func(ctx types.Context) (types.Context, bool, error) {
		var result any

		for _, expr := range args[2:] {
			ctx, result, err = eval.EvalExpression(ctx, expr)
			if err != nil {
				return ctx, false, err
			}
		}

		valid, err := ctx.Coalesce().ToBool(result)
		if err != nil {
			return ctx, false, fmt.Errorf("expression: %w", err)
		}

		return ctx, valid, nil
	}

	// list over vector elements
	if sourceVector, err := ctx.Coalesce().ToVector(source); err == nil {
		output := []any{}

		for i, item := range sourceVector {
			// do not use separate contexts for each loop iteration, as the loop might build up a counter
			innerCtx = innerCtx.WithVariable(valueVarName, item)
			if indexVarName != "" {
				innerCtx = innerCtx.WithVariable(indexVarName, i)
			}

			var result bool
			innerCtx, result, err = filter(innerCtx)
			if err != nil {
				return nil, err
			}

			if result {
				output = append(output, sourceVector[i])
			}
		}

		return output, nil
	}

	// list over object elements
	if sourceObject, err := ctx.Coalesce().ToObject(source); err == nil {
		output := map[string]any{}

		for key, value := range sourceObject {
			// do not use separate contexts for each loop iteration, as the loop might build up a counter
			innerCtx = innerCtx.WithVariable(valueVarName, value)
			if indexVarName != "" {
				innerCtx = innerCtx.WithVariable(indexVarName, key)
			}

			var result bool
			innerCtx, result, err = filter(innerCtx)
			if err != nil {
				return nil, err
			}

			if result {
				output[key] = result
			}
		}

		return output, nil
	}

	return nil, fmt.Errorf("cannot map %T", source)
}

// (filter VECTOR identifier)
// (filter OBJECT identifier)
func anonymousFilterFunction(ctx types.Context, source any, expr ast.Expression) (any, error) {
	identifier, ok := expr.(ast.Identifier)
	if !ok {
		return nil, fmt.Errorf("argument #1: expected identifier, got %T", expr)
	}

	// call the function
	innerCtx := ctx

	filterItem := func(ctx types.Context, item any) (types.Context, bool, error) {
		var (
			result any
			err    error
		)

		ctx, result, err = eval.EvalFunctionCall(ctx, identifier, []ast.Expression{types.MakeShim(item)})
		if err != nil {
			return ctx, false, err
		}

		valid, err := ctx.Coalesce().ToBool(result)
		if err != nil {
			return ctx, false, fmt.Errorf("expression: %w", err)
		}

		return ctx, valid, nil
	}

	if vector, err := ctx.Coalesce().ToVector(source); err == nil {
		output := []any{}

		for i, item := range vector {
			var (
				result bool
				err    error
			)

			innerCtx, result, err = filterItem(innerCtx, item)
			if err != nil {
				return nil, err
			}

			if result {
				output = append(output, vector[i])
			}
		}

		return output, nil
	}

	if object, err := ctx.Coalesce().ToObject(source); err == nil {
		output := map[string]any{}

		for key, value := range object {
			var (
				result bool
				err    error
			)

			innerCtx, result, err = filterItem(innerCtx, value)
			if err != nil {
				return nil, err
			}

			if result {
				output[key] = result
			}
		}

		return output, nil
	}

	// should never happen, as this function call is already gated by a type check
	return nil, fmt.Errorf("cannot apply filter to %T", source)
}

func evalNamingVector(ctx types.Context, expr ast.Expression) (indexName string, valueName string, err error) {
	namingVector, ok := expr.(ast.VectorNode)
	if !ok {
		return "", "", fmt.Errorf("expected a vector, but got %T", expr)
	}

	size := len(namingVector.Expressions)
	if size < 1 || size > 2 {
		return "", "", fmt.Errorf("expected 1 or 2 identifiers in the naming vector, got %d", size)
	}

	if size == 1 {
		varNameIdent, ok := namingVector.Expressions[0].(ast.Identifier)
		if !ok {
			return "", "", fmt.Errorf("value variable name must be an identifier, got %T", namingVector.Expressions[0])
		}

		valueName = varNameIdent.Name
	} else {
		indexIdent, ok := namingVector.Expressions[0].(ast.Identifier)
		if !ok {
			return "", "", fmt.Errorf("index variable name must be an identifier, got %T", namingVector.Expressions[0])
		}

		varNameIdent, ok := namingVector.Expressions[1].(ast.Identifier)
		if !ok {
			return "", "", fmt.Errorf("value variable name must be an identifier, got %T", namingVector.Expressions[0])
		}

		indexName = indexIdent.Name
		valueName = varNameIdent.Name

		if indexName == valueName {
			return "", "", fmt.Errorf("cannot use %s for both value and index variable", indexName)
		}
	}

	return indexName, valueName, nil
}

func stringContainsFunction(haystack string, needle string) (any, error) {
	return strings.Contains(haystack, needle), nil
}

func vectorContainsFunction(ctx types.Context, haystack []any, needle any) (any, error) {
	for _, val := range haystack {
		equal, err := equality.Equal(ctx.Coalesce(), val, needle)
		if err != nil {
			return false, err
		}
		if equal {
			return true, nil
		}
	}

	return false, nil
}
