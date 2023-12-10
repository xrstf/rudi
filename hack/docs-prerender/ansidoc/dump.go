// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package ansidoc

import (
	"strings"

	ansi "github.com/leaanthony/go-ansi-parser"
)

// Dump is useful for debugging the generated ANSI escape sequences.
func Dump(tokens []*ansi.StyledText, rendered string) string {
	var out strings.Builder
	for _, token := range tokens {
		out.WriteString(dumpToken(token, rendered))
	}

	return out.String()
}

func dumpToken(token *ansi.StyledText, rendered string) string {
	var out strings.Builder

	out.WriteString("<node")

	if col := token.FgCol; col != nil {
		out.WriteString(" fg:")
		out.WriteString(col.Hex[1:])
	}

	if col := token.BgCol; col != nil {
		out.WriteString(" bg:")
		out.WriteString(col.Hex[1:])
	}

	out.WriteString(">")

	text := rendered[token.Offset : token.Offset+token.Len]
	cleaned, _ := ansi.Cleanse(text)

	out.WriteString(cleaned)
	out.WriteString("</node>")

	return out.String()
}
