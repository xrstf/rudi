// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package debug

import (
	"io"

	"go.xrstf.de/otto/pkg/lang/ast"
)

func DumpStatement(stmt *ast.Statement, out io.Writer, depth int) error {
	return DumpExpression(stmt.Expression, out, depth)
}
