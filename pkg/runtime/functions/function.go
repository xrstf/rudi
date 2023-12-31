// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package functions

import (
	"errors"
	"fmt"

	"go.xrstf.de/rudi/pkg/coalescing"
	"go.xrstf.de/rudi/pkg/lang/ast"
	"go.xrstf.de/rudi/pkg/runtime/types"
)

type regularFunction struct {
	forms       []form
	coalescer   coalescing.Coalescer
	description string
}

var _ types.Function = &regularFunction{}

func newRegularFunction(forms []form, coalescer coalescing.Coalescer, description string) regularFunction {
	return regularFunction{
		forms:       forms,
		coalescer:   coalescer,
		description: description,
	}
}

func (b *regularFunction) Description() string {
	return b.description
}

func (b *regularFunction) Evaluate(ctx types.Context, args []ast.Expression) (any, error) {
	cachedArgs := convertArgs(args)

	if b.coalescer != nil {
		ctx = ctx.WithCoalescer(b.coalescer)
	}

	for i, form := range b.forms {
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

type BangHandlerFunc func(ctx types.Context, originalArgs []ast.Expression, value any) (types.Context, any, error)

type extendedFunction struct {
	regularFunction

	bangHandler BangHandlerFunc
}

var (
	_ types.Function    = &extendedFunction{}
	_ types.BangHandler = &extendedFunction{}
)

func (f *extendedFunction) BangHandler(ctx types.Context, originalArgs []ast.Expression, value any) (types.Context, any, error) {
	return f.bangHandler(ctx, originalArgs, value)
}
