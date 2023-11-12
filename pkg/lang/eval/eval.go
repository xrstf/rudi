// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package eval

import (
	"go.xrstf.de/corel/pkg/lang/ast"
)

func Run(ctx Context, p *ast.Program) (interface{}, error) {
	result, err := evalProgram(ctx, p)
	if err != nil {
		return nil, err
	}

	return result, nil
}
