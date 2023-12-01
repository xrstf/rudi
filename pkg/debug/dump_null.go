// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package debug

import (
	"io"
)

func DumpNull(out io.Writer) error {
	return writeString(out, "(null)")
}
