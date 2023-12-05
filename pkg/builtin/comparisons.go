// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"errors"

	"go.xrstf.de/rudi/pkg/coalescing"
	"go.xrstf.de/rudi/pkg/equality"
	"go.xrstf.de/rudi/pkg/eval/types"
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

func ltCoalescer(ctx types.Context, left, right any) (any, error) {
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

func lteCoalescer(ctx types.Context, left, right any) (any, error) {
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

func gtCoalescer(ctx types.Context, left, right any) (any, error) {
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

func gteCoalescer(ctx types.Context, left, right any) (any, error) {
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
