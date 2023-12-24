// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package compare

import (
	"errors"

	"go.xrstf.de/rudi/pkg/coalescing"
	"go.xrstf.de/rudi/pkg/equality"
	"go.xrstf.de/rudi/pkg/runtime/functions"
	"go.xrstf.de/rudi/pkg/runtime/types"
)

var (
	strictCoalescer   = coalescing.NewStrict()
	pedanticCoalescer = coalescing.NewPedantic()
	humaneCoalescer   = coalescing.NewHumane()

	Functions = types.Functions{
		"eq?":        functions.NewBuilder(eqFunction).WithDescription("equality check: return true if both arguments are the same").Build(),
		"identical?": functions.NewBuilder(identicalFunction).WithDescription("like `eq?`, but always uses strict coalecsing").Build(),
		"like?":      functions.NewBuilder(likeFunction).WithDescription("like `eq?`, but always uses humane coalecsing").Build(),

		"lt?":  functions.NewBuilder(ltFunction).WithDescription("returns a < b").Build(),
		"lte?": functions.NewBuilder(lteFunction).WithDescription("returns a <= b").Build(),
		"gt?":  functions.NewBuilder(gtFunction).WithDescription("returns a > b").Build(),
		"gte?": functions.NewBuilder(gteFunction).WithDescription("returns a >= b").Build(),
	}
)

func eqFunction(ctx types.Context, left, right any) (any, error) {
	return equality.Equal(ctx.Coalesce(), left, right)
}

func likeFunction(ctx types.Context, left, right any) (any, error) {
	return equality.Equal(coalescing.NewHumane(), left, right)
}

func identicalFunction(ctx types.Context, left, right any) (any, error) {
	return equality.Equal(coalescing.NewStrict(), left, right)
}

func ltFunction(ctx types.Context, left, right any) (any, error) {
	compared, err := equality.Compare(ctx.Coalesce(), left, right)
	if err != nil {
		return nil, err
	}

	switch compared {
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

func lteFunction(ctx types.Context, left, right any) (any, error) {
	compared, err := equality.Compare(ctx.Coalesce(), left, right)
	if err != nil {
		return nil, err
	}

	switch compared {
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

func gtFunction(ctx types.Context, left, right any) (any, error) {
	compared, err := equality.Compare(ctx.Coalesce(), left, right)
	if err != nil {
		return nil, err
	}

	switch compared {
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

func gteFunction(ctx types.Context, left, right any) (any, error) {
	compared, err := equality.Compare(ctx.Coalesce(), left, right)
	if err != nil {
		return nil, err
	}

	switch compared {
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
