// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"errors"
	"fmt"

	"go.xrstf.de/otto/pkg/eval"
	"go.xrstf.de/otto/pkg/eval/types"
	"go.xrstf.de/otto/pkg/lang/ast"
)

func lenFunction(ctx types.Context, args []ast.Expression) (any, error) {
	if size := len(args); size != 1 {
		return nil, fmt.Errorf("expected 1 argument, got %d", size)
	}

	_, list, err := eval.EvalExpression(ctx, args[0])
	if err != nil {
		return nil, err
	}

	if str, ok := list.(ast.String); ok {
		return ast.Number{Value: int64(len(str))}, nil
	}

	if vector, ok := list.(ast.Vector); ok {
		return ast.Number{Value: int64(len(vector.Data))}, nil
	}

	if obj, ok := list.(ast.Object); ok {
		return ast.Number{Value: int64(len(obj.Data))}, nil
	}

	return nil, errors.New("argument is neither a string, vector nor object")
}

func appendFunction(ctx types.Context, args []ast.Expression) (any, error) {
	if size := len(args); size < 2 {
		return nil, fmt.Errorf("expected 2+ arguments, got %d", size)
	}

	_, list, err := eval.EvalExpression(ctx, args[0])
	if err != nil {
		return nil, err
	}

	vector, ok := list.(ast.Vector)
	if !ok {
		return nil, fmt.Errorf("argument #0 is not a vector, but %T", list)
	}

	evaluated, err := evalArgs(ctx, args, 1)
	if err != nil {
		return nil, err
	}

	result := vector.Clone()
	result.Data = append(result.Data, evaluated...)

	return result, nil
}

func prependFunction(ctx types.Context, args []ast.Expression) (any, error) {
	if size := len(args); size < 2 {
		return nil, fmt.Errorf("expected 2+ arguments, got %d", size)
	}

	_, list, err := eval.EvalExpression(ctx, args[0])
	if err != nil {
		return nil, err
	}

	vector, ok := list.(ast.Vector)
	if !ok {
		return nil, fmt.Errorf("argument #0 is not a vector, but %T", list)
	}

	evaluated, err := evalArgs(ctx, args, 1)
	if err != nil {
		return nil, err
	}

	wrapped, err := types.WrapNative(evaluated)
	if err != nil {
		panic("failed to wrap a []any, this should never happen")
	}

	evaluatedVector, ok := wrapped.(ast.Vector)
	if !ok {
		return nil, fmt.Errorf("argument #0 is not a vector, but %T", list)
	}

	evaluatedVector.Data = append(evaluatedVector.Data, vector.Data...)

	return evaluatedVector, nil
}

func reverseFunction(ctx types.Context, args []ast.Expression) (any, error) {
	if size := len(args); size != 1 {
		return nil, fmt.Errorf("expected 1 argument, got %d", size)
	}

	_, list, err := eval.EvalExpression(ctx, args[0])
	if err != nil {
		return nil, err
	}

	strVal, ok := list.(ast.String)
	if ok {
		// thank you https://stackoverflow.com/a/10030772
		result := []rune(strVal)
		for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
			result[i], result[j] = result[j], result[i]
		}

		return ast.String(result), nil
	}

	vector, ok := list.(ast.Vector)
	if ok {
		result := vector.Clone()
		for i, j := 0, len(result.Data)-1; i < j; i, j = i+1, j-1 {
			result.Data[i], result.Data[j] = result.Data[j], result.Data[i]
		}

		return result, nil
	}

	return nil, fmt.Errorf("argument is neither a vector nor a string, but %T", list)
}

// (range VECTOR [item] expr+)
// (range VECTOR [i item] expr+)
// (range OBJECT [val] expr+)
// (range OBJECT [key val] expr+)
func rangeFunction(ctx types.Context, args []ast.Expression) (any, error) {
	if size := len(args); size < 3 {
		return nil, fmt.Errorf("expected 3+ arguments, got %d", size)
	}

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
	if sourceVector, ok := source.(ast.Vector); ok {
		for i, item := range sourceVector.Data {
			// do not use separate contexts for each loop iteration, as the loop might build up a counter
			innerCtx = innerCtx.WithVariable(loopVarName, item)
			if loopIndexName != "" {
				innerCtx = innerCtx.WithVariable(loopIndexName, ast.Number{Value: int64(i)})
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
	if sourceObject, ok := source.(ast.Object); ok {
		for key, value := range sourceObject.Data {
			// do not use separate contexts for each loop iteration, as the loop might build up a counter
			innerCtx = innerCtx.WithVariable(loopVarName, value)
			if loopIndexName != "" {
				innerCtx = innerCtx.WithVariable(loopIndexName, ast.String(key))
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
	if size := len(args); size < 2 {
		return nil, fmt.Errorf("expected 2+ arguments, got %d", size)
	}

	// evaluate the first argument;
	// (map (map .foo +) stuff) should work, so the first argument only needs to _evaluate_
	// to a vector/object, it doesn't need to be a literal objectnode/vectornode.
	innerCtx, source, err := eval.EvalExpression(ctx, args[0])
	if err != nil {
		return nil, fmt.Errorf("argument #0: %w", err)
	}

	var lit ast.Literal
	if vec, ok := source.(ast.Vector); ok {
		lit = vec
	}
	if obj, ok := source.(ast.Object); ok {
		lit = obj
	}

	if lit == nil {
		return nil, fmt.Errorf("argument #0: expected Vector or Object, got %T", source)
	}

	// handle plain function calls
	// (map VECTOR identifier)
	// (map OBJECT identifier)
	if len(args) == 2 {
		return anonymousMapFunction(innerCtx, lit, args[1])
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
	if sourceVector, ok := source.(ast.Vector); ok {
		output := ast.Vector{
			Data: make([]any, len(sourceVector.Data)),
		}

		for i, item := range sourceVector.Data {
			// do not use separate contexts for each loop iteration, as the loop might build up a counter
			innerCtx = innerCtx.WithVariable(valueVarName, item)
			if indexVarName != "" {
				innerCtx = innerCtx.WithVariable(indexVarName, ast.Number{Value: int64(i)})
			}

			var result any
			innerCtx, result, err = mapItem(innerCtx)
			if err != nil {
				return nil, err
			}

			output.Data[i] = result
		}

		return output, nil
	}

	// list over object elements
	if sourceObject, ok := source.(ast.Object); ok {
		output := ast.Object{
			Data: map[string]any{},
		}

		for key, value := range sourceObject.Data {
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

			output.Data[key] = result
		}

		return output, nil
	}

	return nil, fmt.Errorf("cannot map %T", source)
}

// (map VECTOR identifier)
// (map OBJECT identifier)
func anonymousMapFunction(ctx types.Context, source ast.Literal, expr ast.Expression) (any, error) {
	identifier, ok := expr.(ast.Identifier)
	if !ok {
		return nil, fmt.Errorf("argument #1: expected identifier, got %T", expr)
	}

	funcName := string(identifier)

	function, ok := ctx.GetFunction(funcName)
	if !ok {
		return nil, fmt.Errorf("unknown function %s", funcName)
	}

	// call the function
	innerCtx := ctx

	mapItem := func(ctx types.Context, item any) (types.Context, any, error) {
		wrapped, err := types.WrapNative(item)
		if err != nil {
			return ctx, nil, err
		}

		return function(ctx, []ast.Expression{wrapped})
	}

	if vector, ok := source.(ast.Vector); ok {
		output := ast.Vector{
			Data: make([]any, len(vector.Data)),
		}

		for i, item := range vector.Data {
			var (
				result any
				err    error
			)

			innerCtx, result, err = mapItem(innerCtx, item)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", funcName, err)
			}

			output.Data[i] = result
		}

		return output, nil
	}

	if object, ok := source.(ast.Object); ok {
		output := ast.Object{
			Data: map[string]any{},
		}

		for key, value := range object.Data {
			var (
				result any
				err    error
			)

			innerCtx, result, err = mapItem(innerCtx, value)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", funcName, err)
			}

			output.Data[key] = result
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
	if size := len(args); size < 2 {
		return nil, fmt.Errorf("expected 2+ arguments, got %d", size)
	}

	// evaluate the first argument;
	// (filter (filter .foo +) stuff) should work, so the first argument only needs to _evaluate_
	// to a vector/object, it doesn't need to be a literal objectnode/vectornode.
	innerCtx, source, err := eval.EvalExpression(ctx, args[0])
	if err != nil {
		return nil, fmt.Errorf("argument #0: %w", err)
	}

	var lit ast.Literal
	if vec, ok := source.(ast.Vector); ok {
		lit = vec
	}
	if obj, ok := source.(ast.Object); ok {
		lit = obj
	}

	if lit == nil {
		return nil, fmt.Errorf("argument #0: expected Vector or Object, got %T", source)
	}

	// handle plain function calls
	// (filter VECTOR identifier)
	// (filter OBJECT identifier)
	if len(args) == 2 {
		return anonymousFilterFunction(innerCtx, lit, args[1])
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

		native, err := types.UnwrapType(result)
		if err != nil {
			return ctx, false, err
		}

		valid, ok := native.(bool)
		if !ok {
			return ctx, false, fmt.Errorf("expression did not return Bool but %T", result)
		}

		return ctx, valid, nil
	}

	// list over vector elements
	if sourceVector, ok := source.(ast.Vector); ok {
		output := ast.Vector{
			Data: []any{},
		}

		for i, item := range sourceVector.Data {
			// do not use separate contexts for each loop iteration, as the loop might build up a counter
			innerCtx = innerCtx.WithVariable(valueVarName, item)
			if indexVarName != "" {
				innerCtx = innerCtx.WithVariable(indexVarName, ast.Number{Value: int64(i)})
			}

			var result bool
			innerCtx, result, err = filter(innerCtx)
			if err != nil {
				return nil, err
			}

			if result {
				output.Data = append(output.Data, sourceVector.Data[i])
			}
		}

		return output, nil
	}

	// list over object elements
	if sourceObject, ok := source.(ast.Object); ok {
		output := ast.Object{
			Data: map[string]any{},
		}

		for key, value := range sourceObject.Data {
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
				output.Data[key] = result
			}
		}

		return output, nil
	}

	return nil, fmt.Errorf("cannot map %T", source)
}

// (filter VECTOR identifier)
// (filter OBJECT identifier)
func anonymousFilterFunction(ctx types.Context, source ast.Literal, expr ast.Expression) (any, error) {
	identifier, ok := expr.(ast.Identifier)
	if !ok {
		return nil, fmt.Errorf("argument #1: expected identifier, got %T", expr)
	}

	funcName := string(identifier)

	function, ok := ctx.GetFunction(funcName)
	if !ok {
		return nil, fmt.Errorf("unknown function %s", funcName)
	}

	// call the function
	innerCtx := ctx

	filterItem := func(ctx types.Context, item any) (types.Context, bool, error) {
		wrapped, err := types.WrapNative(item)
		if err != nil {
			return ctx, false, err
		}

		var result any
		ctx, result, err = function(ctx, []ast.Expression{wrapped})

		native, err := types.UnwrapType(result)
		if err != nil {
			return ctx, false, err
		}

		valid, ok := native.(bool)
		if !ok {
			return ctx, false, fmt.Errorf("expression did not return Bool but %T", result)
		}

		return ctx, valid, nil
	}

	if vector, ok := source.(ast.Vector); ok {
		output := ast.Vector{
			Data: []any{},
		}

		for i, item := range vector.Data {
			var (
				result bool
				err    error
			)

			innerCtx, result, err = filterItem(innerCtx, item)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", funcName, err)
			}

			if result {
				output.Data = append(output.Data, vector.Data[i])
			}
		}

		return output, nil
	}

	if object, ok := source.(ast.Object); ok {
		output := ast.Object{
			Data: map[string]any{},
		}

		for key, value := range object.Data {
			var (
				result bool
				err    error
			)

			innerCtx, result, err = filterItem(innerCtx, value)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", funcName, err)
			}

			if result {
				output.Data[key] = result
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

		valueName = string(varNameIdent)
	} else {
		indexIdent, ok := namingVector.Expressions[0].(ast.Identifier)
		if !ok {
			return "", "", fmt.Errorf("index variable name must be an identifier, got %T", namingVector.Expressions[0])
		}

		varNameIdent, ok := namingVector.Expressions[1].(ast.Identifier)
		if !ok {
			return "", "", fmt.Errorf("value variable name must be an identifier, got %T", namingVector.Expressions[0])
		}

		indexName = string(indexIdent)
		valueName = string(varNameIdent)

		if indexName == valueName {
			return "", "", fmt.Errorf("cannot use %s for both value and index variable", indexName)
		}
	}

	return indexName, valueName, nil
}
