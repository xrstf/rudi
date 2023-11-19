// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"errors"
	"fmt"

	"go.xrstf.de/otto/pkg/eval"
	"go.xrstf.de/otto/pkg/eval/coalescing"
	"go.xrstf.de/otto/pkg/eval/types"
	"go.xrstf.de/otto/pkg/lang/ast"
)

// (if COND:Expr YES:Expr NO:Expr?)
func ifFunction(ctx types.Context, args []ast.Expression) (any, error) {
	if size := len(args); size < 2 || size > 3 {
		return nil, fmt.Errorf("expected 2 or 3 arguments, got %d", size)
	}

	_, condition, err := eval.EvalExpression(ctx, args[0])
	if err != nil {
		return nil, fmt.Errorf("condition: %w", err)
	}

	success, ok := condition.(ast.Bool)
	if !ok {
		return nil, fmt.Errorf("condition is not bool, but %T", condition)
	}

	if success {
		// discard context changes from the true path
		_, result, err := eval.EvalExpression(ctx, args[1])
		return result, err
	}

	// optional else part
	if len(args) > 2 {
		// discard context changes from the false path
		_, result, err := eval.EvalExpression(ctx, args[2])
		return result, err
	}

	return ast.Null{}, nil
}

// (do STEP:Expr+)
func doFunction(ctx types.Context, args []ast.Expression) (any, error) {
	if size := len(args); size < 1 {
		return nil, fmt.Errorf("expected 1+ arguments, got %d", size)
	}

	tupleCtx := ctx

	var (
		result any
		err    error
	)

	// do not use evalArgs(), as we want to inherit the context between expressions
	for i, arg := range args {
		tupleCtx, result, err = eval.EvalExpression(tupleCtx, arg)
		if err != nil {
			return nil, fmt.Errorf("argument #%d: %w", i, err)
		}
	}

	return result, nil
}

// (has? SYM:SymbolWithPathExpression)
func hasFunction(ctx types.Context, args []ast.Expression) (any, error) {
	if size := len(args); size != 1 {
		return nil, fmt.Errorf("expected 1 argument, got %d", size)
	}

	symbol, ok := args[0].(ast.Symbol)
	if !ok {
		return nil, fmt.Errorf("argument #0 is not a symbol, but %T", args[0])
	}

	if symbol.PathExpression == nil {
		return nil, errors.New("argument #0 has no path expression")
	}

	if symbol.Variable != nil {
		varName := string(*symbol.Variable)

		if _, ok = ctx.GetVariable(varName); !ok {
			return nil, fmt.Errorf("unknown variable %s", varName)
		}
	}

	// do a syntax check by pre-computing the path
	evaluatedPath, err := eval.EvalPathExpression(ctx, symbol.PathExpression)
	if err != nil {
		return nil, fmt.Errorf("argument #0: invalid path expression: %w", err)
	}

	_, value, err := eval.EvalSymbolWithEvaluatedPath(ctx, symbol, *evaluatedPath)
	if err != nil {
		return false, nil
	}

	return ast.Bool(value != nil), nil
}

// (default TEST:Expression FALLBACK:any)
func defaultFunction(ctx types.Context, args []ast.Expression) (any, error) {
	if size := len(args); size != 2 {
		return nil, fmt.Errorf("expected 2 arguments, got %d", size)
	}

	_, result, err := eval.EvalExpression(ctx, args[0])
	if err != nil {
		return nil, fmt.Errorf("argument #0: %w", err)
	}

	isEmpty, err := coalescing.IsEmpty(result)
	if err != nil {
		return nil, fmt.Errorf("argument #0: %w", err)
	}

	if !isEmpty {
		return result, nil
	}

	_, result, err = eval.EvalExpression(ctx, args[1])
	if err != nil {
		return nil, fmt.Errorf("argument #1: %w", err)
	}

	return result, nil
}

// (try TEST:Expression FALLBACK:any?)
func tryFunction(ctx types.Context, args []ast.Expression) (any, error) {
	if size := len(args); size < 1 || size > 2 {
		return nil, fmt.Errorf("expected 1 or 2 arguments, got %d", size)
	}

	_, result, err := eval.EvalExpression(ctx, args[0])
	if err != nil {
		if len(args) == 1 {
			return nil, nil
		}

		_, result, err = eval.EvalExpression(ctx, args[1])
		if err != nil {
			return nil, fmt.Errorf("argument #1: %w", err)
		}
	}

	return result, nil
}

// (set VAR:Variable VALUE:any)
// (set EXPR:PathExpression VALUE:any)
func setFunction(ctx types.Context, args []ast.Expression) (types.Context, any, error) {
	if size := len(args); size != 2 {
		return ctx, nil, fmt.Errorf("expected 2 arguments, got %d", size)
	}

	symbol, ok := args[0].(ast.Symbol)
	if !ok {
		return ctx, nil, fmt.Errorf("argument #0 is not a symbol, but %T", args[0])
	}

	// catch symbols that are technically invalid
	if symbol.Variable == nil && symbol.PathExpression == nil {
		return ctx, nil, fmt.Errorf("argument #0: must be path expression or variable, got %s", symbol.ExpressionName())
	}

	// discard any context changes within the newValue expression
	_, newValue, err := eval.EvalExpression(ctx, args[1])
	if err != nil {
		return ctx, nil, fmt.Errorf("argument #1: %w", err)
	}

	// pre-evaluate the path
	var pathExpr *ast.EvaluatedPathExpression
	if p := symbol.PathExpression; p != nil {
		pathExpr, err = eval.EvalPathExpression(ctx, p)
		if err != nil {
			return ctx, nil, fmt.Errorf("argument #1: invalid path expression: %w", err)
		}
	}

	// get the current value
	var currentValue any

	if symbol.Variable != nil {
		varName := string(*symbol.Variable)

		// a non-existing variable is fine, this is how you define new variables in the first place
		currentValue, _ = ctx.GetVariable(varName)
	} else {
		doc := ctx.GetDocument()
		currentValue = doc.Get()
	}

	// if there is a path expression, merge in the new value
	updatedValue := newValue
	if pathExpr != nil {
		updatedValue, err = setValueAtPath(currentValue, pathExpr.Steps, newValue)
		if err != nil {
			return ctx, nil, fmt.Errorf("cannot set value in %T at %s: %w", currentValue, pathExpr, err)
		}
	}

	// set a variable, which will result in a new context
	if symbol.Variable != nil {
		varName := string(*symbol.Variable)

		// make the variable's value the return value, so `(def $foo 12)` = 12
		return ctx.WithVariable(varName, updatedValue), newValue, nil
	}

	// update the global document
	// (the document Go struct stays the same, so this does not result in a new context)
	doc := ctx.GetDocument()
	doc.Set(updatedValue)

	return ctx, newValue, nil
}

func setValueAtPath(dest any, steps []ast.EvaluatedPathStep, newValue any) (any, error) {
	if len(steps) == 0 {
		return newValue, nil
	}

	target, err := types.UnwrapType(dest)
	if err != nil {
		return nil, fmt.Errorf("cannot descend into %T", dest)
	}

	thisStep := steps[0]
	remainingSteps := steps[1:]

	// [index]...
	if iv := thisStep.IntegerValue; iv != nil {
		index := int(*iv)
		if index < 0 {
			return nil, fmt.Errorf("index %d out of bounds", index)
		}

		if vector, ok := target.([]any); ok {
			if index >= len(vector) {
				return nil, fmt.Errorf("index %d out of bounds", index)
			}

			existingValue := vector[index]

			updatedValue, err := setValueAtPath(existingValue, remainingSteps, newValue)
			if err != nil {
				return nil, err
			}

			vector[index] = updatedValue

			return vector, nil
		}

		return nil, fmt.Errorf("cannot descend with [%d] into %T", index, target)
	}

	// .key
	if sv := thisStep.StringValue; sv != nil {
		key := *sv

		if object, ok := target.(map[string]any); ok {
			// getting the empty value for non-existing keys is fine
			existingValue, _ := object[key]

			updatedValue, err := setValueAtPath(existingValue, remainingSteps, newValue)
			if err != nil {
				return nil, err
			}

			object[key] = updatedValue

			return object, nil
		}

		// nulls can be turned into objects
		if target == nil {
			updatedValue, err := setValueAtPath(nil, remainingSteps, newValue)
			if err != nil {
				return nil, err
			}

			return map[string]any{
				key: updatedValue,
			}, nil
		}

		return nil, fmt.Errorf("cannot descend with [%s] into %T", key, target)
	}

	return nil, errors.New("invalid path step: neither key nor index")
}

// (empty? VALUE:any)
func isEmptyFunction(ctx types.Context, args []ast.Expression) (any, error) {
	if size := len(args); size != 1 {
		return nil, fmt.Errorf("expected 1 or 2 arguments, got %d", size)
	}

	_, result, err := eval.EvalExpression(ctx, args[0])
	if err != nil {
		return nil, err
	}

	switch asserted := result.(type) {
	case ast.Bool:
		return bool(asserted) == false, nil
	case ast.Number:
		return asserted.ToFloat() == 0, nil
	case ast.String:
		return len(string(asserted)) == 0, nil
	case ast.Null:
		return true, nil
	case ast.Vector:
		return len(asserted.Data) == 0, nil
	case ast.Object:
		return len(asserted.Data) == 0, nil
	default:
		return nil, fmt.Errorf("unexpected argument %v (%T)", result, result)
	}
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

			for _, expr := range args[2:] {
				innerCtx, result, err = eval.EvalExpression(innerCtx, expr)
				if err != nil {
					return nil, err
				}
			}

			// TODO: this updates the source vector, but it should be idempotent: return a map copy instead
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

			for _, expr := range args[2:] {
				innerCtx, result, err = eval.EvalExpression(innerCtx, expr)
				if err != nil {
					return nil, err
				}
			}

			// TODO: this updates the source vector, but it should be idempotent: return a map copy instead
			output.Data[key] = result
		}

		return output, nil
	}

	return nil, fmt.Errorf("cannot map %T", source)
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

	if vector, ok := source.(ast.Vector); ok {
		output := ast.Vector{
			Data: make([]any, len(vector.Data)),
		}

		for i, item := range vector.Data {
			wrapped, err := types.WrapNative(item)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", funcName, err)
			}

			var result any
			innerCtx, result, err = function(innerCtx, []ast.Expression{wrapped})
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
			wrapped, err := types.WrapNative(value)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", funcName, err)
			}

			var result any
			innerCtx, result, err = function(innerCtx, []ast.Expression{wrapped})
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
