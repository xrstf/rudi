// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package debug

import (
	"io"

	"go.xrstf.de/otto/pkg/lang/ast"
)

func DumpNull(n *ast.Null, out io.Writer) error {
	return writeString(out, "(null)")
}
