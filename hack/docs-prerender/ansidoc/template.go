// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package ansidoc

import (
	"strings"

	ansi "github.com/leaanthony/go-ansi-parser"
)

// Templatify turns a (usually Optimize()d) list of styled tokens into a custom homegrown template
// syntax, which is later easy to postprocess when actually rendering the documents. The code that
// parses the templates and injects the proper colors lives in the cmd/rudi module, but the
// templatifier needs to live here to keep the go-ansi-parser dependency out of cmd/rudi.
func Templatify(tokens []*ansi.StyledText, rendered string) string {
	var out strings.Builder
	for _, token := range tokens {
		out.WriteString(templatifyToken(token, rendered))
	}

	lines := strings.Split(out.String(), "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " ")
	}

	return strings.Join(lines, "\n")
}

func templatifyToken(token *ansi.StyledText, rendered string) string {
	var out strings.Builder

	out.WriteString("{¤")

	if fg := token.FgCol; fg != nil {
		out.WriteString(fg.Hex[1:])
		out.WriteString(":")
	}

	text := rendered[token.Offset : token.Offset+token.Len]
	cleaned, _ := ansi.Cleanse(text)

	out.WriteString(cleaned)
	out.WriteString("¤}")

	return out.String()
}
