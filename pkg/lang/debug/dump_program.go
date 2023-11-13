// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package debug

import (
	"fmt"
	"io"

	"go.xrstf.de/corel/pkg/lang/ast"
)

func dumpProgram(p *ast.Program, out io.Writer, depth int) error {
	for _, stmt := range p.Statements {
		if err := dumpStatement(&stmt, out, depth); err != nil {
			return fmt.Errorf("failed to dump statement %s: %w", stmt.String(), err)
		}
	}

	return nil
}
