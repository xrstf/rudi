// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"time"

	"go.xrstf.de/rudi/pkg/eval/types"
)

func nowFunction(ctx types.Context, args []any) (any, error) {
	formatString, err := ctx.Coalesce().ToString(args[0])
	if err != nil {
		return nil, err
	}

	formatted := time.Now().Format(string(formatString))

	return formatted, nil
}
