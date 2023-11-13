// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package debug

import (
	"fmt"
	"io"

	"go.xrstf.de/corel/pkg/lang/ast"
)

func dumpString(str *ast.String, out io.Writer) error {
	return writeString(out, fmt.Sprintf("(string %s)", str))
}
