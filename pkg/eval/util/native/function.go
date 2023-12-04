// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package native

import (
	"errors"
	"fmt"

	"go.xrstf.de/rudi/pkg/eval/types"
	"go.xrstf.de/rudi/pkg/lang/ast"
)

type function struct {
	forms       []form
	description string
}

var _ types.Function = &function{}

func NewFunction(description string, funcs ...any) types.Function {
	forms := make([]form, len(funcs))

	for i := range funcs {
		funcForm, err := newForm(funcs[i])
		if err != nil {
			panic(fmt.Sprintf("Form #%d is invalid: %v", i, err))
		}

		forms[i] = funcForm
	}

	return &function{
		forms:       forms,
		description: description,
	}
}

func (f *function) Description() string {
	return f.description
}

func (f *function) Evaluate(ctx types.Context, args []ast.Expression) (any, error) {
	cachedArgs := convertArgs(args)

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
