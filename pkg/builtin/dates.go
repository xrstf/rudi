// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"time"
)

func nowFunction(format string) (any, error) {
	return time.Now().Format(format), nil
}
