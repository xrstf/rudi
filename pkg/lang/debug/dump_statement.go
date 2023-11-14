// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package debug

import (
	"io"

	"go.xrstf.de/otto/pkg/lang/ast"
)

func dumpStatement(stmt *ast.Statement, out io.Writer, depth int) error {
	return dumpTuple(&stmt.Tuple, out, depth)
}
