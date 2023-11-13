// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package debug

import (
	"fmt"
	"io"

	"go.xrstf.de/otto/pkg/lang/ast"
)

func dumpIdentifier(ident *ast.Identifier, out io.Writer) error {
	return writeString(out, fmt.Sprintf("(identifier %s)", *ident))
}
