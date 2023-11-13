// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package debug

import (
	"fmt"
	"io"

	"go.xrstf.de/corel/pkg/lang/ast"
)

func dumpNumber(num *ast.Number, out io.Writer) error {
	return writeString(out, fmt.Sprintf("(number %s)", num))
}

// hack until PathExpression steps do use proper Expressions
func dumpInteger(i *int64, out io.Writer) error {
	return writeString(out, fmt.Sprintf("(number %d)", *i))
}
