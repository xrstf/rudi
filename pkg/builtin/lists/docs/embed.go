// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package docs

import (
	"embed"
	_ "embed"

	rudidocs "go.xrstf.de/rudi/pkg/docs"
)

//go:embed *.md
var embeddedFS embed.FS

var Functions = rudidocs.NewFunctionProvider(&embeddedFS)
