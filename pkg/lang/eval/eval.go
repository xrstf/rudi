// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package eval

import (
	"go.xrstf.de/corel/pkg/lang/ast"
	"go.xrstf.de/corel/pkg/lang/eval/types"
)

func Run(ctx types.Context, p *ast.Program) (interface{}, error) {
	result, err := evalProgram(ctx, p)
	if err != nil {
		return nil, err
	}

	return result, nil
}
