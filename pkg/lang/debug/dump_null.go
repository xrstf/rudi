// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package debug

import (
	"io"

	"go.xrstf.de/corel/pkg/lang/ast"
)

func dumpNull(n *ast.Null, out io.Writer) error {
	return writeString(out, "(null)")
}
