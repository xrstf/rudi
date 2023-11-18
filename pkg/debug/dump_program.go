// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package debug

import (
	"fmt"
	"io"

	"go.xrstf.de/otto/pkg/lang/ast"
)

func DumpProgram(p *ast.Program, out io.Writer, depth int) error {
	for _, stmt := range p.Statements {
		if err := DumpStatement(&stmt, out, depth); err != nil {
			return fmt.Errorf("failed to dump statement %s: %w", stmt.String(), err)
		}

		separator := "\n"
		if depth == DoNotIndent {
			separator = " "
		}

		if err := writeString(out, separator); err != nil {
			return err
		}
	}

	return nil
}
