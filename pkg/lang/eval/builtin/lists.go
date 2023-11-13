// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"errors"
	"fmt"

	"go.xrstf.de/otto/pkg/lang/ast"
	"go.xrstf.de/otto/pkg/lang/eval/types"
)

func lenFunction(ctx types.Context, args []Argument) (any, error) {
	if size := len(args); size != 1 {
		return nil, fmt.Errorf("expected 1 argument, got %d", size)
	}

	_, list, err := args[0].Eval(ctx)
	if err != nil {
		return nil, err
	}

	vector, ok := list.(ast.Vector)
	if !ok {
		return nil, errors.New("argument is not a vector")
	}

	return ast.Number{Value: len(vector.Data)}, nil
}
