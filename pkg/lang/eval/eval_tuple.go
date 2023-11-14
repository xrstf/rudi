// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package eval

import (
	"errors"
	"fmt"

	"go.xrstf.de/otto/pkg/lang/ast"
	"go.xrstf.de/otto/pkg/lang/eval/types"
)

func EvalTuple(ctx types.Context, tup ast.Tuple) (types.Context, any, error) {
	if len(tup.Expressions) == 0 {
		return ctx, nil, errors.New("invalid tuple: tuple cannot be empty")
	}

	funcExpr, ok := tup.Expressions[0].(ast.Identifier)
	if !ok {
		return ctx, nil, errors.New("invalid tuple: first expression must be an identifier")
	}

	funcName := string(funcExpr)

	function, ok := ctx.GetFunction(funcName)
	if !ok {
		return ctx, nil, fmt.Errorf("unknown function %s", funcName)
	}

	// call the function
	newContext, result, err := function(ctx, tup.Expressions[1:])
	if err != nil {
		return ctx, nil, fmt.Errorf("%s: %w", funcName, err)
	}

	return newContext, result, nil
}
