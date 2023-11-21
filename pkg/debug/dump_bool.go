// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package debug

import (
	"fmt"
	"io"

	"go.xrstf.de/rudi/pkg/lang/ast"
)

func DumpBool(b *ast.Bool, out io.Writer) error {
	return writeString(out, fmt.Sprintf("(bool %s)", b))
}
