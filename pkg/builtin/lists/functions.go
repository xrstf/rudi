// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package lists

import (
	"fmt"

	"go.xrstf.de/rudi/pkg/builtin/helper"
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
	loopIndexName, loopVarName, err := helper.DecodeNamingVector(namingVec)
	if err != nil {
		return nil, fmt.Errorf("argument #1: not a valid naming vector: %w", err)
	}

	var result any

	for i, item := range data {
		vars := map[string]any{
			loopVarName: item,
		}

		if loopIndexName != "" {
			vars[loopIndexName] = i
		}

		// Do not use separate contexts for each loop iteration, as the loop might build up a counter,
		// but only use the loop variables in a shallow scope, where only these two temporary variables
		// are laid over the regular scoped/global variables.
		result, err = ctx.Runtime().EvalExpression(ctx.NewShallowScope(nil, vars), expr)
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
	loopIndexName, loopVarName, err := helper.DecodeNamingVector(namingVec)
	if err != nil {
		return nil, fmt.Errorf("argument #1: not a valid naming vector: %w", err)
	}

	var (
		result any
	)

	for key, value := range data {
		vars := map[string]any{
			loopVarName: value,
		}

		if loopIndexName != "" {
			vars[loopIndexName] = key
		}

		result, err = ctx.Runtime().EvalExpression(ctx.NewShallowScope(nil, vars), expr)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

// itemHandlerFunc works for iterating over vectors as well as over objects.
type itemHandlerFunc func(ctx types.Context, _ any, value any) (any, error)

// (map VECTOR identifier)
func mapVectorAnonymousFunction(ctx types.Context, data []any, ident ast.Identifier) (any, error) {
	mapHandler := func(ctx types.Context, _ any, value any) (any, error) {
		return ctx.Runtime().CallFunction(ctx, ident, []ast.Expression{types.MakeShim(value)})
	}

	return mapVector(ctx, data, mapHandler)
}

// (map VECTOR [item] expr)
// (map VECTOR [i item] expr)
func mapVectorExpressionFunction(ctx types.Context, data []any, namingVec ast.Expression, expr ast.Expression) (any, error) {
	indexVarName, valueVarName, err := helper.DecodeNamingVector(namingVec)
	if err != nil {
		return nil, fmt.Errorf("argument #1: not a valid naming vector: %w", err)
	}

	mapHandler := func(ctx types.Context, index any, value any) (any, error) {
		vars := map[string]any{
			valueVarName: value,
		}

		if indexVarName != "" {
			vars[indexVarName] = index
		}

		return ctx.Runtime().EvalExpression(ctx.NewShallowScope(nil, vars), expr)
	}

	return mapVector(ctx, data, mapHandler)
}

func mapVector(ctx types.Context, data []any, f itemHandlerFunc) (any, error) {
	output := make([]any, len(data))

	for i, item := range data {
		result, err := f(ctx, i, item)
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

	mapHandler := func(ctx types.Context, _ any, value any) (any, error) {
		return ctx.Runtime().CallFunction(ctx, identifier, []ast.Expression{types.MakeShim(value)})
	}

	return mapObject(ctx, data, mapHandler)
}

// (map OBJECT [item] expr)
// (map OBJECT [i item] expr)
func mapObjectExpressionFunction(ctx types.Context, data map[string]any, namingVec ast.Expression, expr ast.Expression) (any, error) {
	keyVarName, valueVarName, err := helper.DecodeNamingVector(namingVec)
	if err != nil {
		return nil, fmt.Errorf("argument #1: not a valid naming vector: %w", err)
	}

	mapHandler := func(ctx types.Context, key any, value any) (any, error) {
		vars := map[string]any{
			valueVarName: value,
		}

		if keyVarName != "" {
			vars[keyVarName] = key
		}

		return ctx.Runtime().EvalExpression(ctx.NewShallowScope(nil, vars), expr)
	}

	return mapObject(ctx, data, mapHandler)
}

func mapObject(ctx types.Context, data map[string]any, f itemHandlerFunc) (any, error) {
	output := map[string]any{}

	for key, value := range data {
		result, err := f(ctx, key, value)
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

	mapHandler := func(ctx types.Context, _ any, value any) (any, error) {
		return ctx.Runtime().CallFunction(ctx, identifier, []ast.Expression{types.MakeShim(value)})
	}

	return filterVector(ctx, data, mapHandler)
}

// (filter VECTOR [item] expr)
// (filter VECTOR [i item] expr)
func filterVectorExpressionFunction(ctx types.Context, data []any, namingVec ast.Expression, expr ast.Expression) (any, error) {
	indexVarName, valueVarName, err := helper.DecodeNamingVector(namingVec)
	if err != nil {
		return nil, fmt.Errorf("argument #1: not a valid naming vector: %w", err)
	}

	mapHandler := func(ctx types.Context, index any, value any) (any, error) {
		vars := map[string]any{
			valueVarName: value,
		}

		if indexVarName != "" {
			vars[indexVarName] = index
		}

		return ctx.Runtime().EvalExpression(ctx.NewShallowScope(nil, vars), expr)
	}

	return filterVector(ctx, data, mapHandler)
}

func filterVector(ctx types.Context, data []any, f itemHandlerFunc) (any, error) {
	output := []any{}

	for i, item := range data {
		result, err := f(ctx, i, item)
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

	mapHandler := func(ctx types.Context, _ any, value any) (any, error) {
		return ctx.Runtime().CallFunction(ctx, identifier, []ast.Expression{types.MakeShim(value)})
	}

	return filterObject(ctx, data, mapHandler)
}

// (filter OBJECT [item] expr)
// (filter OBJECT [i item] expr)
func filterObjectExpressionFunction(ctx types.Context, data map[string]any, namingVec ast.Expression, expr ast.Expression) (any, error) {
	keyVarName, valueVarName, err := helper.DecodeNamingVector(namingVec)
	if err != nil {
		return nil, fmt.Errorf("argument #1: not a valid naming vector: %w", err)
	}

	mapHandler := func(ctx types.Context, key any, value any) (any, error) {
		vars := map[string]any{
			valueVarName: value,
		}

		if keyVarName != "" {
			vars[keyVarName] = key
		}

		return ctx.Runtime().EvalExpression(ctx.NewShallowScope(nil, vars), expr)
	}

	return filterObject(ctx, data, mapHandler)
}

func filterObject(ctx types.Context, data map[string]any, f itemHandlerFunc) (any, error) {
	output := map[string]any{}

	for key, value := range data {
		result, err := f(ctx, key, value)
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
