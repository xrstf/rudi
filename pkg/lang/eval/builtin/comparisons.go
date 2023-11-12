// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"errors"
	"fmt"

	"go.xrstf.de/corel/pkg/lang/eval/coalescing"
)

func eqFunction(args []interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, errors.New("(eq EXPRESSION EXPRESSION)")
	}

	left := args[0]
	right := args[1]

	switch leftAsserted := left.(type) {
	case string:
		rightAsserted, err := coalescing.ToString(right)
		if err != nil {
			return nil, fmt.Errorf("cannot compare %T with %T", left, right)
		}

		return rightAsserted == leftAsserted, nil
	}

	return false, fmt.Errorf("do not know how to compare %T with anything", left)
}
