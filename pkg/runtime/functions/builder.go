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

// NewBuilder returns a new Rudi function builder by combining multiple "forms"
// into one Rudi function. Each form is a regular Go function that can have
// arbitrary arguments, but must return (any, error).
// When called from within a Rudi program, the first compatible form is evaluated,
// so care must be taken to not accidentally shadow a form with another.
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
