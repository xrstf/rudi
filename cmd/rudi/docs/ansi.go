// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package docs

import (
	"encoding/hex"
	"strconv"
	"strings"
)

// TextEffect is a type representing the ansi text effects. This type and its
// values match https://github.com/leaanthony/go-ansi-parser/ types.
type TextEffect int

const (
	Bold TextEffect = 1 << iota
	Faint
	Italic
	Blinking
	Inversed
	Invisible
	Underlined
	Strikethrough
	Bright
)

type TextStyle struct {
	Foreground string
	Background string
	Effects    TextEffect
}

func (s TextStyle) toAnsi() string {
	params := strings.Join(s.toParams(), ";")
	return "\033[0;" + params + "m"
}

func (s TextStyle) toParams() []string {
	var params []string
	if s.Effects&Bold == Bold {
		params = append(params, "1")
	}
	if s.Effects&Faint == Faint {
		params = append(params, "2")
	}
	if s.Effects&Italic == Italic {
		params = append(params, "3")
	}
	if s.Effects&Underlined == Underlined {
		params = append(params, "4")
	}
	if s.Effects&Blinking == Blinking {
		params = append(params, "5")
	}
	if s.Effects&Inversed == Inversed {
		params = append(params, "7")
	}
	if s.Effects&Invisible == Invisible {
		params = append(params, "8")
	}
	if s.Effects&Strikethrough == Strikethrough {
		params = append(params, "9")
	}

	if s.Foreground != "" {
		params = append(params, colorToAnsi("38", s.Foreground)...)
	}

	if s.Background != "" {
		params = append(params, colorToAnsi("48", s.Background)...)
	}

	return params
}

func colorToAnsi(kind string, color string) []string {
	// assume color is a hex string
	if len(color) == 6 {
		b, err := hex.DecodeString(color)
		if err == nil {
			return []string{
				kind,
				"2",
				strconv.FormatInt(int64(b[0]), 10),
				strconv.FormatInt(int64(b[1]), 10),
				strconv.FormatInt(int64(b[2]), 10),
			}
		}
	}

	// if not, treat it as 256 color index instead
	return []string{kind, "5", color}
}
