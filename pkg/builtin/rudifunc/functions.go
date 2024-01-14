// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package rudifunc

import (
	"errors"
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
// would just instantly vanish into thin air.
func funcFunction(ctx types.Context, name ast.Expression, namingVector ast.Expression, body ...ast.Expression) (any, error) {
	return nil, nil
}

// funcBangHandler is where the side effect of adding a new function to the Rudi runtime actually happens.
func funcBangHandler(ctx types.Context, originalArgs []ast.Expression) (any, error) {
	if len(originalArgs) < 3 {
		return nil, errors.New("none of the available forms matched the given expressions")
	}

	nameIdent, ok := originalArgs[0].(ast.Identifier)
	if !ok {
		return nil, fmt.Errorf("first argument must be an identifier that specifies the function name, but got %T instead", originalArgs[0])
	}

	// naming vector must be a vector consisting only of identifiers
	paramVector, ok := originalArgs[1].(ast.VectorNode)
	if !ok {
		return nil, fmt.Errorf("second argument must be vector containing the parameter names, got %T instead", originalArgs[1])
	}

	paramNames := []string{}
	for _, param := range paramVector.Expressions {
		paramIdent, ok := param.(ast.Identifier)
		if !ok {
			return nil, fmt.Errorf("parameter vector must contain only identifiers, got %T instead", param)
		}
		paramNames = append(paramNames, paramIdent.Name)
	}

	f := rudispaceFunc{
		name:   nameIdent.Name,
		params: paramNames,
		body:   originalArgs[2:],
	}

	ctx.SetRudispaceFunction(f.name, f)

	return nil, nil
}

type rudispaceFunc struct {
	name   string
	params []string
	body   []ast.Expression
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
		arg, err := ctx.Runtime().EvalExpression(ctx, args[i])
		if err != nil {
			return nil, err
		}

		funcArgs[paramName] = arg
	}

	// user-defined functions form a sub-program and all statements share the same context
	funcCtx := ctx.NewScope()
	funcCtx.SetVariables(funcArgs)

	runtime := ctx.Runtime()

	var (
		result any
		err    error
	)

	for _, expr := range f.body {
		result, err = runtime.EvalExpression(funcCtx, expr)
		if err != nil {
			return nil, err
		}
	}

	return result, err
}
