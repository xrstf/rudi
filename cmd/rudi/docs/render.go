// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package docs

import (
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"

	"github.com/muesli/termenv"
)

type PainterFunc func(nodeKind Node, text string, hexColor string) (newText string, newStyle TextStyle)

func IdentityPainter(nodeKind Node, text string, hexColor string) (newText string, newStyle TextStyle) {
	return text, TextStyle{
		Foreground: intToHex(int(nodeKind)),
	}
}

func GlamourPainter(nodeKind Node, text string, hexColor string) (newText string, newStyle TextStyle) {
	var style map[Node]TextStyle

	if termenv.HasDarkBackground() {
		style = darkStyle
	} else {
		style = lightStyle
	}

	elementStyle, ok := style[nodeKind]
	if !ok {
		return IdentityPainter(nodeKind, text, hexColor)
	}

	return text, elementStyle
}

var thePattern = regexp.MustCompile(`(?s){¤(([0-9a-f]+):)?(.*?)¤}`)

func Render(tpl string, painter PainterFunc) string {
	if painter == nil {
		painter = GlamourPainter
	}

	matches := thePattern.FindAllStringSubmatch(tpl, -1)
	for _, match := range matches {
		var (
			nodeKind Node
			hexColor string
		)

		if color := match[2]; len(color) > 0 {
			hexColor = color
			nodeKind = Node(hexToInt(hexColor))
		}

		newText, newStyle := painter(nodeKind, match[3], hexColor)

		begin := newStyle.toAnsi()
		end := "\033[0m"

		replacement := fmt.Sprintf("%s%s%s", begin, newText, end)
		tpl = strings.Replace(tpl, match[0], replacement, 1)
	}

	return tpl
}

func hexToInt(h string) int {
	hexData, _ := hex.DecodeString(h)

	return int(hexData[0])<<16 + int(hexData[1])<<8 + int(hexData[2])
}

func intToHex(i int) string {
	r := byte(i >> 16)
	g := byte(i >> 8)
	b := byte(i)

	return hex.EncodeToString([]byte{r, g, b})
}
