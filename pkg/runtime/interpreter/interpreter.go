// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package interpreter

import (
	"go.xrstf.de/rudi/pkg/runtime/types"
)

type interpreter struct{}

var _ types.Runtime = &interpreter{}

func New() types.Runtime {
	return &interpreter{}
}
