// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package rudifunc

import (
	"fmt"

	"go.xrstf.de/rudi/pkg/lang/ast"
	"go.xrstf.de/rudi/pkg/runtime/functions"
	"go.xrstf.de/rudi/pkg/runtime/types"
)

var (
	Functions = types.Functions{
		"func": functions.NewBuilder(funcFunction).WithBangHandler(funcBangHandler).WithDescription("defines a new function").Build(),
	}
)

// funcFunction should never be called without the bang modifier, as without it, the created function
// just instantly vanishes into thin air.
func funcFunction(ctx types.Context, name ast.Expression, namingVector ast.Expression, body ast.Expression) (any, error) {
	nameIdent, ok := name.(ast.Identifier)
	if !ok {
		return nil, fmt.Errorf("first argument must be an identifier that specifies the function name, but got %T instead", name)
	}

	// naming vector must be a vector consisting only of identifiers
	paramVector, ok := namingVector.(ast.VectorNode)
	if !ok {
		return nil, fmt.Errorf("second argument must be vector containing the parameter names, got %T instead", namingVector)
	}

	paramNames := []string{}
	for _, param := range paramVector.Expressions {
		paramIdent, ok := param.(ast.Identifier)
		if !ok {
			return nil, fmt.Errorf("parameter vector must contain only identifiers, got %T instead", param)
		}
		paramNames = append(paramNames, paramIdent.Name)
	}

	return rudispaceFunc{
		name:   nameIdent.Name,
		params: paramNames,
		body:   body,
	}, nil
}

// funcBangHandler is where the side effect of adding a new function to the Rudi runtime actually happens.
func funcBangHandler(ctx types.Context, originalArgs []ast.Expression, value any) (types.Context, any, error) {
	intermediate, ok := value.(rudispaceFunc)
	if !ok {
		panic("This should never happen: func! bang handler called with non-intermediate function.")
	}

	return ctx.WithRudispaceFunction(intermediate.name, intermediate), nil, nil
}

type rudispaceFunc struct {
	name   string
	params []string
	body   ast.Expression
}

var _ types.Function = rudispaceFunc{}

func (rudispaceFunc) Description() string {
	return "" // no docs required/useful for rudispace functions
}

func (f rudispaceFunc) Evaluate(ctx types.Context, args []ast.Expression) (any, error) {
	if len(args) != len(f.params) {
		return nil, fmt.Errorf("expected %d argument(s), got %d", len(f.params), len(args))
	}

	funcArgs := map[string]any{}
	for i, paramName := range f.params {
		_, arg, err := ctx.Runtime().EvalExpression(ctx, args[i])
		if err != nil {
			return nil, err
		}

		funcArgs[paramName] = arg
	}

	_, result, err := ctx.Runtime().EvalExpression(ctx.WithVariables(funcArgs), f.body)

	return result, err
}
