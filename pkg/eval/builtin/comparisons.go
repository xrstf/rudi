// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"errors"
	"fmt"

	"go.xrstf.de/rudi/pkg/coalescing"
	"go.xrstf.de/rudi/pkg/equality"
	"go.xrstf.de/rudi/pkg/eval"
	"go.xrstf.de/rudi/pkg/eval/types"
	"go.xrstf.de/rudi/pkg/lang/ast"
)

func makeEqualityFunc(coalescerGetter func(ctx types.Context) coalescing.Coalescer, desc string) types.Function {
	return types.BasicFunction(func(ctx types.Context, args []ast.Expression) (any, error) {
		if size := len(args); size != 2 {
			return nil, fmt.Errorf("expected exactly 2 arguments, got %d", size)
		}

		_, leftData, err := eval.EvalExpression(ctx, args[0])
		if err != nil {
			return nil, fmt.Errorf("argument #0: %w", err)
		}

		_, rightData, err := eval.EvalExpression(ctx, args[1])
		if err != nil {
			return nil, fmt.Errorf("argument #1: %w", err)
		}

		return equality.Equal(coalescerGetter(ctx), leftData, rightData)
	}, desc)
}

type comparisonCoalescer func(result int) (bool, error)

func makeComparatorFunc(cc comparisonCoalescer, desc string) types.Function {
	return types.BasicFunction(func(ctx types.Context, args []ast.Expression) (any, error) {
		if size := len(args); size != 2 {
			return nil, fmt.Errorf("expected 2 argument(s), got %d", size)
		}

		_, left, err := eval.EvalExpression(ctx, args[0])
		if err != nil {
			return nil, fmt.Errorf("argument #0: %w", err)
		}

		_, right, err := eval.EvalExpression(ctx, args[1])
		if err != nil {
			return nil, fmt.Errorf("argument #1: %w", err)
		}

		compared, err := equality.Compare(ctx.Coalesce(), left, right)
		if err != nil {
			return nil, err
		}

		return cc(compared)
	}, desc)
}

func ltCoalescer(result int) (bool, error) {
	switch result {
	case equality.IsEqual:
		return false, nil
	case equality.IsSmaller:
		return true, nil
	case equality.IsGreater:
		return false, nil
	case equality.Unorderable:
		return false, errors.New("cannot order the given arguments")
	default:
		panic("Unexpected comparison result.")
	}
}

func lteCoalescer(result int) (bool, error) {
	switch result {
	case equality.IsEqual:
		return true, nil
	case equality.IsSmaller:
		return true, nil
	case equality.IsGreater:
		return false, nil
	case equality.Unorderable:
		return false, errors.New("cannot order the given arguments")
	default:
		panic("Unexpected comparison result.")
	}
}

func gtCoalescer(result int) (bool, error) {
	switch result {
	case equality.IsEqual:
		return false, nil
	case equality.IsSmaller:
		return false, nil
	case equality.IsGreater:
		return true, nil
	case equality.Unorderable:
		return false, errors.New("cannot order the given arguments")
	default:
		panic("Unexpected comparison result.")
	}
}

func gteCoalescer(result int) (bool, error) {
	switch result {
	case equality.IsEqual:
		return true, nil
	case equality.IsSmaller:
		return false, nil
	case equality.IsGreater:
		return true, nil
	case equality.Unorderable:
		return false, errors.New("cannot order the given arguments")
	default:
		panic("Unexpected comparison result.")
	}
}
