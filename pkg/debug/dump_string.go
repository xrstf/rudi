// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package debug

import (
	"fmt"
	"io"
)

func DumpString(str string, out io.Writer) error {
	return writeString(out, fmt.Sprintf("(string %q)", str))
}
