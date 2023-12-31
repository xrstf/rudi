// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package printer

import (
	"io"
)

type basePrinter struct {
	out io.Writer
}

func newBasePrinter(out io.Writer) basePrinter {
	return basePrinter{
		out: out,
	}
}

func (p *basePrinter) write(str string) error {
	_, err := p.out.Write([]byte(str))
	return err
}
