// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package ansidoc

import (
	ansi "github.com/leaanthony/go-ansi-parser"
)

// Parse parses a rendered Markdown string with ANSI sequencces into a sequence of StyledText
// structs.
func Parse(markdown string) ([]*ansi.StyledText, error) {
	return ansi.Parse(markdown)
}
