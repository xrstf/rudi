// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package debug

import (
	"fmt"
	"io"
)

func DumpBool(b bool, out io.Writer) error {
	return writeString(out, fmt.Sprintf("(bool %v)", b))
}
