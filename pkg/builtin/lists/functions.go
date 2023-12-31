// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package lists

import (
	"fmt"

	"go.xrstf.de/rudi/pkg/lang/ast"
	"go.xrstf.de/rudi/pkg/runtime/functions"
	"go.xrstf.de/rudi/pkg/runtime/types"
)

var (
	Functions = types.Functions{
		// these ones are defined in the strings module, but also work with lists;
		// since strings is a builtin module and it's unexpected for anyone to configure
		// Rudi without string functions, these functions are not mirrored here.
		// "len":       …
		// "append":    …
		// "prepend":   …
		// "reverse":   …
		// "contains?": …

		"range": functions.
			NewBuilder(
				rangeVectorFunction,
				rangeObjectFunction,
			).
			WithDescription("allows to iterate (loop) over a vector or object").
			Build(),

		"map": functions.
			NewBuilder(
				mapVectorExpressionFunction,
				mapObjectExpressionFunction,
				mapVectorAnonymousFunction,
				mapObjectAnonymousFunction,
			).
			WithDescription("applies an expression to every element in a vector or object").
			Build(),

		"filter": functions.
			NewBuilder(
				filterVectorExpressionFunction,
				filterObjectExpressionFunction,
				filterVectorAnonymousFunction,
				filterObjectAnonymousFunction,
			).
			WithDescription("returns a copy of a given vector/object with only those elements remaining that satisfy a condition").
			Build(),
	}
)

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
		vars := map[string]any{
			loopVarName: item,
		}

		if loopIndexName != "" {
			vars[loopIndexName] = i
		}

		// do not use separate contexts for each loop iteration, as the loop might build up a counter
		loopCtx = loopCtx.WithVariables(vars)

		loopCtx, result, err = ctx.Runtime().EvalExpression(loopCtx, expr)
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
		vars := map[string]any{
			loopVarName: value,
		}

		if loopIndexName != "" {
			vars[loopIndexName] = key
		}

		// do not use separate contexts for each loop iteration, as the loop might build up a counter
		loopCtx = loopCtx.WithVariables(vars)

		loopCtx, result, err = ctx.Runtime().EvalExpression(loopCtx, expr)
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
		return ctx.Runtime().CallFunction(ctx, identifier, []ast.Expression{types.MakeShim(value)})
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
		vars := map[string]any{
			valueVarName: value,
		}

		if indexVarName != "" {
			vars[indexVarName] = index
		}

		ctx = ctx.WithVariables(vars)

		return ctx.Runtime().EvalExpression(ctx, expr)
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
		return ctx.Runtime().CallFunction(ctx, identifier, []ast.Expression{types.MakeShim(value)})
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
		vars := map[string]any{
			valueVarName: value,
		}

		if keyVarName != "" {
			vars[keyVarName] = key
		}

		ctx = ctx.WithVariables(vars)

		return ctx.Runtime().EvalExpression(ctx, expr)
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
		return ctx.Runtime().CallFunction(ctx, identifier, []ast.Expression{types.MakeShim(value)})
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
		vars := map[string]any{
			valueVarName: value,
		}

		if indexVarName != "" {
			vars[indexVarName] = index
		}

		ctx = ctx.WithVariables(vars)

		return ctx.Runtime().EvalExpression(ctx, expr)
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
		return ctx.Runtime().CallFunction(ctx, identifier, []ast.Expression{types.MakeShim(value)})
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
		vars := map[string]any{
			valueVarName: value,
		}

		if keyVarName != "" {
			vars[keyVarName] = key
		}

		ctx = ctx.WithVariables(vars)

		return ctx.Runtime().EvalExpression(ctx, expr)
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
