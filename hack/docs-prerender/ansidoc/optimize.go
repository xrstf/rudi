// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package ansidoc

import (
	ansi "github.com/leaanthony/go-ansi-parser"
)

type style struct {
	Foreground string
	Background string
	Style      ansi.TextStyle
}

// Optimize removes redundant style declarations and merges adjunct nodes together.
func Optimize(tokens []*ansi.StyledText) []*ansi.StyledText {
	output := []*ansi.StyledText{}
	var tip *ansi.StyledText

	for i, token := range tokens {
		// first item
		if tip == nil {
			output = append(output, tokens[i])
			tip = output[len(output)-1]
			continue
		}

		var newStyle *style
		if colorChange(tip.FgCol, token.FgCol) {
			fg := ""
			if token.FgCol != nil {
				fg = token.FgCol.Hex
			}

			newStyle = &style{
				Foreground: fg,
			}
		}

		if colorChange(tip.BgCol, token.BgCol) {
			if newStyle == nil {
				newStyle = &style{}
			}

			bg := ""
			if token.BgCol != nil {
				bg = token.BgCol.Hex
			}

			newStyle.Background = bg
		}

		if token.Style != tip.Style {
			if newStyle == nil {
				newStyle = &style{}
			}

			newStyle.Style = token.Style
		}

		// If there is no change, extend the tip token,
		// otherwise produce a new token with the changed style.
		if newStyle == nil {
			tip.Len += token.Len
		} else {
			newToken := &ansi.StyledText{
				Label:      token.Label,
				Style:      newStyle.Style,
				ColourMode: token.ColourMode,
				Offset:     token.Offset,
				Len:        token.Len,
			}

			if newStyle.Foreground != "" {
				newToken.FgCol = &ansi.Col{
					Hex: newStyle.Foreground,
				}
			}

			if newStyle.Background != "" {
				newToken.BgCol = &ansi.Col{
					Hex: newStyle.Background,
				}
			}

			output = append(output, newToken)
			tip = output[len(output)-1]
		}
	}

	return output
}

func colorChange(old, new *ansi.Col) bool {
	if new == nil && old == nil {
		return false
	}

	return new == nil || old == nil || old.Hex != new.Hex
}
