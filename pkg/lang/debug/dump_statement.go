// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package debug

import (
	"io"

	"go.xrstf.de/corel/pkg/lang/ast"
)

func dumpStatement(stmt *ast.Statement, out io.Writer, depth int) error {
	if err := dumpTuple(&stmt.Tuple, out, depth); err != nil {
		return err
	}

	out.Write([]byte("\n"))

	return nil
}
