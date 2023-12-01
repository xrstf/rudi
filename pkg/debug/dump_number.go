// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package debug

import (
	"fmt"
	"io"
)

func DumpNumber(value any, out io.Writer) error {
	return writeString(out, fmt.Sprintf("(number (%T %v))", value, value))
}
