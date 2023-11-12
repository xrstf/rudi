// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"errors"
	"fmt"

	"go.xrstf.de/corel/pkg/lang/ast"
	"go.xrstf.de/corel/pkg/lang/eval/coalescing"
	"go.xrstf.de/corel/pkg/lang/eval/types"
)

func eqFunction(ctx types.Context, args []Argument) (any, error) {
	if len(args) != 2 {
		return nil, errors.New("(eq LEFT RIGHT)")
	}

	_, leftData, err := args[0].Eval(ctx)
	if err != nil {
		return nil, fmt.Errorf("argument #0: %w", err)
	}

	_, rightData, err := args[1].Eval(ctx)
	if err != nil {
		return nil, fmt.Errorf("argument #1: %w", err)
	}

	switch leftAsserted := leftData.(type) {
	case ast.String:
		rightAsserted, err := coalescing.ToString(rightData)
		if err != nil {
			return nil, fmt.Errorf("cannot compare %T with %T", leftData, rightData)
		}

		return ast.Bool(string(leftAsserted) == rightAsserted), nil
	}

	return ast.Bool(false), fmt.Errorf("do not know how to compare %T with anything", leftData)
}
