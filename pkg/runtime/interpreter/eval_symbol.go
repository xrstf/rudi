// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package interpreter

import (
	"errors"
	"fmt"

	"go.xrstf.de/rudi/pkg/lang/ast"
	"go.xrstf.de/rudi/pkg/runtime/pathexpr"
	"go.xrstf.de/rudi/pkg/runtime/types"
)

func (*interpreter) EvalSymbol(ctx types.Context, sym ast.Symbol) (types.Context, any, error) {
	rootValue := ctx.GetDocument().Data()

	// sanity check
	if sym.Variable == nil && sym.PathExpression == nil {
		return ctx, nil, errors.New("invalid symbol")
	}

	// . always returns the root document
	if sym.IsDot() {
		return ctx, rootValue, nil
	}

	// if this symbol is a variable, replace the root value with the variable's value
	if sym.Variable != nil {
		var ok bool

		varName := string(*sym.Variable)

		rootValue, ok = ctx.GetVariable(varName)
		if !ok {
			return ctx, nil, fmt.Errorf("unknown variable %s", varName)
		}
	}

	deeper, err := pathexpr.Apply(ctx, rootValue, sym.PathExpression)
	if err != nil {
		return ctx, nil, fmt.Errorf("cannot evaluate %s: %w", sym.String(), err)
	}

	return ctx, deeper, nil
}
