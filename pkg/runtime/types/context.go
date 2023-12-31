// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package types

import (
	"context"
	"errors"

	"go.xrstf.de/rudi/pkg/coalescing"
	"go.xrstf.de/rudi/pkg/lang/ast"
)

type Context struct {
	ctx        context.Context
	document   *Document
	fixedFuncs Functions
	userFuncs  Functions
	variables  Variables
	coalescer  coalescing.Coalescer
	runtime    Runtime
}

func NewContext(runtime Runtime, ctx context.Context, doc Document, variables Variables, funcs Functions, coalescer coalescing.Coalescer) (Context, error) {
	if runtime == nil {
		return Context{}, errors.New("no runtime provided")
	}

	if ctx == nil {
		ctx = context.Background()
	}

	if variables == nil {
		variables = NewVariables()
	}

	if funcs == nil {
		funcs = NewFunctions()
	}

	if coalescer == nil {
		coalescer = coalescing.NewStrict()
	}

	return Context{
		ctx:        ctx,
		document:   &doc,
		fixedFuncs: funcs,
		userFuncs:  NewFunctions(),
		variables:  variables,
		coalescer:  coalescer,
		runtime:    runtime,
	}, nil
}

// Coalesce is named this way to make the frequent calls read fluently
// (for example "ctx.Coalesce().ToBool(...)").
func (c Context) Coalesce() coalescing.Coalescer {
	return c.coalescer
}

func (c Context) GoContext() context.Context {
	return c.ctx
}

func (c Context) Runtime() Runtime {
	return c.runtime
}

func (c Context) GetDocument() *Document {
	return c.document
}

func (c Context) GetVariable(name string) (any, bool) {
	return c.variables.Get(name)
}

func (c Context) GetFunction(name string) (Function, bool) {
	f, ok := c.fixedFuncs.Get(name)
	if ok {
		return f, true
	}

	return c.userFuncs.Get(name)
}

func (c Context) WithGoContext(ctx context.Context) Context {
	clone := c.shallowCopy()
	clone.ctx = ctx

	return clone
}

func (c Context) WithVariable(name string, val any) Context {
	clone := c.shallowCopy()
	clone.variables = c.variables.With(name, val)

	return clone
}

func (c Context) WithVariables(vars map[string]any) Context {
	if len(vars) == 0 {
		return c
	}

	clone := c.shallowCopy()
	clone.variables = c.variables.WithMany(vars)

	return clone
}

func (c Context) WithCoalescer(coalescer coalescing.Coalescer) Context {
	clone := c.shallowCopy()
	clone.coalescer = coalescer

	return clone
}

func (c Context) WithRudispaceFunction(funcName string, fun Function) Context {
	clone := c.shallowCopy()
	clone.userFuncs = c.userFuncs.DeepCopy().Set(funcName, fun)

	return clone
}

func (c Context) shallowCopy() Context {
	return Context{
		ctx:        c.ctx,
		document:   c.document,
		fixedFuncs: c.fixedFuncs,
		userFuncs:  c.userFuncs,
		variables:  c.variables,
		coalescer:  c.coalescer,
		runtime:    c.runtime,
	}
}

func MakeShim(val any) ast.Shim {
	return ast.Shim{Value: val}
}
