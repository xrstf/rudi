// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package otto

import (
	"strings"

	"go.xrstf.de/otto/pkg/lang/parser"
)

type ParseError struct {
	script string
	err    error
}

var _ error = ParseError{}

func (p ParseError) Error() string {
	return p.err.Error()
}

func (p ParseError) Snippet() string {
	if el, ok := p.err.(parser.ErrorLister); ok {
		var buffer strings.Builder

		for _, e := range el.Errors() {
			parserErr, ok := e.(parser.ParserError)
			if !ok {
				return ""
			}

			_, col, off := parserErr.Pos()
			line := extractLine(p.script, off)
			if col >= len(line) {
				col = len(line) - 1
			} else {
				if col > 0 {
					col--
				}
			}
			if col < 0 {
				col = 0
			}
			pos := col
			for _, chr := range line[:col] {
				if chr == '\t' {
					pos += 7
				}
			}

			buffer.WriteString(line + "\n")
			buffer.WriteString(strings.Repeat(" ", pos) + "^")
		}

		return buffer.String()
	}

	return ""
}

func extractLine(input string, initPos int) string {
	if initPos < 0 {
		initPos = 0
	}
	if initPos >= len(input) && len(input) > 0 {
		initPos = len(input) - 1
	}
	startPos := initPos
	endPos := initPos
	for ; startPos > 0; startPos-- {
		if input[startPos] == '\n' {
			if startPos != initPos {
				startPos++
				break
			}
		}
	}
	for ; endPos < len(input); endPos++ {
		if input[endPos] == '\n' {
			if endPos == initPos {
				endPos++
			}
			break
		}
	}
	return input[startPos:endPos]
}
