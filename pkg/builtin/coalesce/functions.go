// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package coalesce

import (
	"go.xrstf.de/rudi/pkg/builtin/core"
	"go.xrstf.de/rudi/pkg/coalescing"
	"go.xrstf.de/rudi/pkg/eval/functions"
	"go.xrstf.de/rudi/pkg/eval/types"
)

var (
	strictCoalescer   = coalescing.NewStrict()
	pedanticCoalescer = coalescing.NewPedantic()
	humaneCoalescer   = coalescing.NewHumane()

	Functions = types.Functions{
		"strictly":     functions.NewBuilder(core.DoFunction).WithCoalescer(strictCoalescer).WithDescription("evaluates the child expressions using strict coalescing").Build(),
		"pedantically": functions.NewBuilder(core.DoFunction).WithCoalescer(pedanticCoalescer).WithDescription("evaluates the child expressions using pedantic coalescing").Build(),
		"humanely":     functions.NewBuilder(core.DoFunction).WithCoalescer(humaneCoalescer).WithDescription("evaluates the child expressions using humane coalescing").Build(),
	}
)
