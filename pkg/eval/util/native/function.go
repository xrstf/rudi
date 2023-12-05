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

func NewFunction(forms ...any) *Function {
	funcForms := make([]form, len(forms))

	for i := range forms {
		funcForm, err := newForm(forms[i])
		if err != nil {
			panic(fmt.Sprintf("Form #%d is invalid: %v", i, err))
		}

		funcForms[i] = funcForm
	}

	return &Function{
		forms: funcForms,
	}
}

func (f *Function) Description() string {
	return f.description
}

func (f *Function) WithDescription(s string) *Function {
	f.description = s
	return f
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
