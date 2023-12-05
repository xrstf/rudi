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

// (range VECTOR [item] expr)
// (range VECTOR [i item] expr)
func rangeVectorFunction(ctx types.Context, data []any, namingVec ast.Expression, expr ast.Expression) (any, error) {
	// decode desired loop variable namings
	loopIndexName, loopVarName, err := decodeNamingVector(namingVec)
	if err != nil {
		return nil, fmt.Errorf("argument #1: not a valid naming vector: %w", err)
	}

	var (
		result  any
		loopCtx = ctx
	)

	for i, item := range data {
		// do not use separate contexts for each loop iteration, as the loop might build up a counter
		loopCtx = loopCtx.WithVariable(loopVarName, item)
		if loopIndexName != "" {
			loopCtx = loopCtx.WithVariable(loopIndexName, i)
		}

		loopCtx, result, err = eval.EvalExpression(loopCtx, expr)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

// (range OBJECT [val] expr)
// (range OBJECT [key val] expr)
func rangeObjectFunction(ctx types.Context, data map[string]any, namingVec ast.Expression, expr ast.Expression) (any, error) {
	// decode desired loop variable namings
	loopIndexName, loopVarName, err := decodeNamingVector(namingVec)
	if err != nil {
		return nil, fmt.Errorf("argument #1: not a valid naming vector: %w", err)
	}

	var (
		result  any
		loopCtx = ctx
	)

	for key, value := range data {
		// do not use separate contexts for each loop iteration, as the loop might build up a counter
		loopCtx = loopCtx.WithVariable(loopVarName, value)
		if loopIndexName != "" {
			loopCtx = loopCtx.WithVariable(loopIndexName, key)
		}

		loopCtx, result, err = eval.EvalExpression(loopCtx, expr)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

// itemHandlerFunc works for iterating over vectors as well as over objects.
type itemHandlerFunc func(ctx types.Context, _ any, value any) (types.Context, any, error)

// (map VECTOR identifier)
func mapVectorAnonymousFunction(ctx types.Context, data []any, ident ast.Expression) (any, error) {
	// type check
	identifier, ok := ident.(ast.Identifier)
	if !ok {
		return nil, fmt.Errorf("argument #1: expected identifier, got %T", ident)
	}

	mapHandler := func(ctx types.Context, _ any, value any) (types.Context, any, error) {
		return eval.EvalFunctionCall(ctx, identifier, []ast.Expression{types.MakeShim(value)})
	}

	return mapVector(ctx, data, mapHandler)
}

// (map VECTOR [item] expr)
// (map VECTOR [i item] expr)
func mapVectorExpressionFunction(ctx types.Context, data []any, namingVec ast.Expression, expr ast.Expression) (any, error) {
	indexVarName, valueVarName, err := decodeNamingVector(namingVec)
	if err != nil {
		return nil, fmt.Errorf("argument #1: not a valid naming vector: %w", err)
	}

	mapHandler := func(ctx types.Context, index any, value any) (types.Context, any, error) {
		ctx = ctx.WithVariable(valueVarName, value)
		if indexVarName != "" {
			ctx = ctx.WithVariable(indexVarName, index)
		}

		return eval.EvalExpression(ctx, expr)
	}

	return mapVector(ctx, data, mapHandler)
}

func mapVector(ctx types.Context, data []any, f itemHandlerFunc) (any, error) {
	output := make([]any, len(data))
	loopCtx := ctx

	for i, item := range data {
		var (
			result any
			err    error
		)

		// do not use separate contexts for each loop iteration, as the loop might build up a counter
		loopCtx, result, err = f(loopCtx, i, item)
		if err != nil {
			return nil, err
		}

		output[i] = result
	}

	return output, nil
}

// (map OBJECT identifier)
func mapObjectAnonymousFunction(ctx types.Context, data map[string]any, ident ast.Expression) (any, error) {
	// type check
	identifier, ok := ident.(ast.Identifier)
	if !ok {
		return nil, fmt.Errorf("argument #1: expected identifier, got %T", ident)
	}

	mapHandler := func(ctx types.Context, _ any, value any) (types.Context, any, error) {
		return eval.EvalFunctionCall(ctx, identifier, []ast.Expression{types.MakeShim(value)})
	}

	return mapObject(ctx, data, mapHandler)
}

// (map OBJECT [item] expr)
// (map OBJECT [i item] expr)
func mapObjectExpressionFunction(ctx types.Context, data map[string]any, namingVec ast.Expression, expr ast.Expression) (any, error) {
	keyVarName, valueVarName, err := decodeNamingVector(namingVec)
	if err != nil {
		return nil, fmt.Errorf("argument #1: not a valid naming vector: %w", err)
	}

	mapHandler := func(ctx types.Context, key any, value any) (types.Context, any, error) {
		ctx = ctx.WithVariable(valueVarName, value)
		if keyVarName != "" {
			ctx = ctx.WithVariable(keyVarName, key)
		}

		return eval.EvalExpression(ctx, expr)
	}

	return mapObject(ctx, data, mapHandler)
}

func mapObject(ctx types.Context, data map[string]any, f itemHandlerFunc) (any, error) {
	output := map[string]any{}
	loopCtx := ctx

	for key, value := range data {
		var (
			result any
			err    error
		)

		// do not use separate contexts for each loop iteration, as the loop might build up a counter
		loopCtx, result, err = f(loopCtx, key, value)
		if err != nil {
			return nil, err
		}

		output[key] = result
	}

	return output, nil
}

// (filter VECTOR identifier)
func filterVectorAnonymousFunction(ctx types.Context, data []any, ident ast.Expression) (any, error) {
	// type check
	identifier, ok := ident.(ast.Identifier)
	if !ok {
		return nil, fmt.Errorf("argument #1: expected identifier, got %T", ident)
	}

	mapHandler := func(ctx types.Context, _ any, value any) (types.Context, any, error) {
		return eval.EvalFunctionCall(ctx, identifier, []ast.Expression{types.MakeShim(value)})
	}

	return filterVector(ctx, data, mapHandler)
}

// (filter VECTOR [item] expr)
// (filter VECTOR [i item] expr)
func filterVectorExpressionFunction(ctx types.Context, data []any, namingVec ast.Expression, expr ast.Expression) (any, error) {
	indexVarName, valueVarName, err := decodeNamingVector(namingVec)
	if err != nil {
		return nil, fmt.Errorf("argument #1: not a valid naming vector: %w", err)
	}

	mapHandler := func(ctx types.Context, index any, value any) (types.Context, any, error) {
		ctx = ctx.WithVariable(valueVarName, value)
		if indexVarName != "" {
			ctx = ctx.WithVariable(indexVarName, index)
		}

		return eval.EvalExpression(ctx, expr)
	}

	return filterVector(ctx, data, mapHandler)
}

func filterVector(ctx types.Context, data []any, f itemHandlerFunc) (any, error) {
	output := []any{}
	loopCtx := ctx

	for i, item := range data {
		var (
			result any
			err    error
		)

		loopCtx, result, err = f(loopCtx, i, item)
		if err != nil {
			return nil, err
		}

		valid, err := ctx.Coalesce().ToBool(result)
		if err != nil {
			return nil, fmt.Errorf("expression: %w", err)
		}

		if valid {
			output = append(output, data[i])
		}
	}

	return output, nil
}

// (filter OBJECT identifier)
func filterObjectAnonymousFunction(ctx types.Context, data map[string]any, ident ast.Expression) (any, error) {
	// type check
	identifier, ok := ident.(ast.Identifier)
	if !ok {
		return nil, fmt.Errorf("argument #1: expected identifier, got %T", ident)
	}

	mapHandler := func(ctx types.Context, _ any, value any) (types.Context, any, error) {
		return eval.EvalFunctionCall(ctx, identifier, []ast.Expression{types.MakeShim(value)})
	}

	return filterObject(ctx, data, mapHandler)
}

// (filter OBJECT [item] expr)
// (filter OBJECT [i item] expr)
func filterObjectExpressionFunction(ctx types.Context, data map[string]any, namingVec ast.Expression, expr ast.Expression) (any, error) {
	keyVarName, valueVarName, err := decodeNamingVector(namingVec)
	if err != nil {
		return nil, fmt.Errorf("argument #1: not a valid naming vector: %w", err)
	}

	mapHandler := func(ctx types.Context, key any, value any) (types.Context, any, error) {
		ctx = ctx.WithVariable(valueVarName, value)
		if keyVarName != "" {
			ctx = ctx.WithVariable(keyVarName, key)
		}

		return eval.EvalExpression(ctx, expr)
	}

	return filterObject(ctx, data, mapHandler)
}

func filterObject(ctx types.Context, data map[string]any, f itemHandlerFunc) (any, error) {
	output := map[string]any{}
	loopCtx := ctx

	for key, value := range data {
		var (
			result any
			err    error
		)

		loopCtx, result, err = f(loopCtx, key, value)
		if err != nil {
			return nil, err
		}

		valid, err := ctx.Coalesce().ToBool(result)
		if err != nil {
			return nil, fmt.Errorf("expression: %w", err)
		}

		if valid {
			output[key] = data[key]
		}
	}

	return output, nil
}

func decodeNamingVector(expr ast.Expression) (indexName string, valueName string, err error) {
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
