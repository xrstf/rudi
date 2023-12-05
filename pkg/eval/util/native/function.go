// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package native

import (
	"errors"
	"fmt"

	"go.xrstf.de/rudi/pkg/coalescing"
	"go.xrstf.de/rudi/pkg/eval/types"
	"go.xrstf.de/rudi/pkg/lang/ast"
)

type Function struct {
	forms       []form
	coalescer   coalescing.Coalescer
	description string
}

var _ types.Function = &Function{}

func NewFunction(description string, funcs ...any) *Function {
	forms := make([]form, len(funcs))

	for i := range funcs {
		funcForm, err := newForm(funcs[i])
		if err != nil {
			panic(fmt.Sprintf("Form #%d is invalid: %v", i, err))
		}

		forms[i] = funcForm
	}

	return &Function{
		forms:       forms,
		description: description,
	}
}

func (f *Function) Description() string {
	return f.description
}

func (f *Function) WithCoalescer(c coalescing.Coalescer) *Function {
	f.coalescer = c
	return f
}

func (f *Function) Evaluate(ctx types.Context, args []ast.Expression) (any, error) {
	cachedArgs := convertArgs(args)

	if f.coalescer != nil {
		ctx = ctx.WithCoalescer(f.coalescer)
	}

	for i, form := range f.forms {
		matched, err := form.Match(ctx, cachedArgs)
		if err != nil {
			return nil, fmt.Errorf("form#%d: %w", i, err)
		}

		if matched {
			return form.Call(ctx)
		}
	}

	return nil, errors.New("none of the available forms matched the given expressions")
}
