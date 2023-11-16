// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package debug

import (
	"fmt"
	"io"

	"go.xrstf.de/otto/pkg/lang/ast"
)

func DumpString(str *ast.String, out io.Writer) error {
	return writeString(out, fmt.Sprintf("(string %s)", str))
}
