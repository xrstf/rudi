// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"go.xrstf.de/rudi/pkg/builtin/coalesce"
	"go.xrstf.de/rudi/pkg/builtin/compare"
	"go.xrstf.de/rudi/pkg/builtin/core"
	"go.xrstf.de/rudi/pkg/builtin/datetime"
	"go.xrstf.de/rudi/pkg/builtin/encoding"
	"go.xrstf.de/rudi/pkg/builtin/hashing"
	"go.xrstf.de/rudi/pkg/builtin/lists"
	"go.xrstf.de/rudi/pkg/builtin/logic"
	"go.xrstf.de/rudi/pkg/builtin/math"
	"go.xrstf.de/rudi/pkg/builtin/rudifunc"
	"go.xrstf.de/rudi/pkg/builtin/strings"
	"go.xrstf.de/rudi/pkg/builtin/types"
	evaltypes "go.xrstf.de/rudi/pkg/runtime/types"
)

var (
	SafeFunctions = evaltypes.Functions{}.
			Add(core.Functions).
			Add(logic.Functions).
			Add(coalesce.Functions).
			Add(compare.Functions).
			Add(math.Functions).
			Add(strings.Functions).
			Add(lists.Functions).
			Add(hashing.Functions).
			Add(encoding.Functions).
			Add(datetime.Functions).
			Add(types.Functions)

	RudifuncFunctions = rudifunc.Functions

	UnsafeFunctions = evaltypes.Functions{}.
			Add(RudifuncFunctions)
)
