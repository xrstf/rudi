// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"errors"

	"go.xrstf.de/rudi/pkg/coalescing"
	"go.xrstf.de/rudi/pkg/equality"
	"go.xrstf.de/rudi/pkg/eval/types"
)

func eqFunction(ctx types.Context, args []any) (any, error) {
	return equality.Equal(ctx.Coalesce(), args[0], args[1])
}

func likeFunction(ctx types.Context, args []any) (any, error) {
	return equality.Equal(coalescing.NewHumane(), args[0], args[1])
}

func identicalFunction(ctx types.Context, args []any) (any, error) {
	return equality.Equal(coalescing.NewStrict(), args[0], args[1])
}

func ltCoalescer(ctx types.Context, args []any) (any, error) {
	compared, err := equality.Compare(ctx.Coalesce(), args[0], args[1])
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

func lteCoalescer(ctx types.Context, args []any) (any, error) {
	compared, err := equality.Compare(ctx.Coalesce(), args[0], args[1])
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

func gtCoalescer(ctx types.Context, args []any) (any, error) {
	compared, err := equality.Compare(ctx.Coalesce(), args[0], args[1])
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

func gteCoalescer(ctx types.Context, args []any) (any, error) {
	compared, err := equality.Compare(ctx.Coalesce(), args[0], args[1])
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
