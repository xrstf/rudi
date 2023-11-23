// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package eval

import (
	"errors"
	"fmt"

	"go.xrstf.de/rudi/pkg/eval/types"
	"go.xrstf.de/rudi/pkg/lang/ast"
)

func EvalSymbol(ctx types.Context, sym ast.Symbol) (types.Context, any, error) {
	// pre-evaluate the path expression
	pathExpr := ast.EvaluatedPathExpression{}

	if sym.PathExpression != nil {
		evaluated, err := EvalPathExpression(ctx, sym.PathExpression)
		if err != nil {
			return ctx, nil, fmt.Errorf("invalid path expression: %w", err)
		}
		pathExpr = *evaluated
	}

	return EvalSymbolWithEvaluatedPath(ctx, sym, pathExpr)
}

func EvalSymbolWithEvaluatedPath(ctx types.Context, sym ast.Symbol, path ast.EvaluatedPathExpression) (types.Context, any, error) {
	rootDoc := ctx.GetDocument()
	rootValue := rootDoc.Get()

	// sanity check
	if sym.Variable == nil && sym.PathExpression == nil {
		return ctx, nil, errors.New("invalid symbol")
	}

	// . always returns the root document
	if sym.IsDot() {
		wrapped, err := types.WrapNative(rootValue)
		if err != nil {
			return ctx, nil, err
		}

		return ctx, wrapped, nil
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

	deeper, err := TraverseEvaluatedPathExpression(ctx, rootValue, path)
	if err != nil {
		return ctx, nil, fmt.Errorf("cannot evaluate %s: %w", sym.String(), err)
	}

	return ctx, deeper, nil
}
