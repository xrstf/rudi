// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package functions

import (
	"fmt"

	"go.xrstf.de/rudi/pkg/coalescing"
	"go.xrstf.de/rudi/pkg/runtime/types"
)

type Builder struct {
	forms       []form
	coalescer   coalescing.Coalescer
	bangHandler BangHandlerFunc
	description string
}

func NewBuilder(forms ...any) *Builder {
	b := &Builder{
		forms: []form{},
	}

	for _, f := range forms {
		b.AddForm(f)
	}

	return b
}

func (b *Builder) AddForm(form any) *Builder {
	funcForm, err := newForm(form)
	if err != nil {
		panic(fmt.Sprintf("invalid form: %v", err))
	}

	b.forms = append(b.forms, funcForm)
	return b
}

func (b *Builder) WithDescription(s string) *Builder {
	b.description = s
	return b
}

func (b *Builder) WithCoalescer(c coalescing.Coalescer) *Builder {
	b.coalescer = c
	return b
}

func (b *Builder) WithBangHandler(h BangHandlerFunc) *Builder {
	b.bangHandler = h
	return b
}

func (b *Builder) Build() types.Function {
	f := newRegularFunction(b.forms, b.coalescer, b.description)

	if b.bangHandler != nil {
		return &extendedFunction{
			regularFunction: f,
			bangHandler:     b.bangHandler,
		}
	}

	return &f
}
