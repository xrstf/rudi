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
	ctx             context.Context
	document        *Document
	fixedFuncs      Functions
	userFuncs       Functions
	globalVariables Variables
	scopeVariables  Variables
	tempVariables   Variables
	coalescer       coalescing.Coalescer
	runtime         Runtime
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
		ctx:             ctx,
		document:        &doc,
		fixedFuncs:      funcs,
		userFuncs:       NewFunctions(),
		globalVariables: variables,
		scopeVariables:  NewVariables(),
		coalescer:       coalescer,
		runtime:         runtime,
	}, nil
}

func (c Context) NewScope() Context {
	clone := c.shallowCopy()
	clone.scopeVariables = NewVariables()
	clone.tempVariables = nil

	return clone
}

func (c Context) NewShallowScope(extraVars Variables) Context {
	clone := c.shallowCopy()
	clone.tempVariables = extraVars

	return clone
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
	value, ok := c.tempVariables.Get(name)
	if ok {
		return value, true
	}

	value, ok = c.scopeVariables.Get(name)
	if ok {
		return value, true
	}

	return c.globalVariables.Get(name)
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

func (c Context) WithCoalescer(coalescer coalescing.Coalescer) Context {
	clone := c.shallowCopy()
	clone.coalescer = coalescer

	return clone
}

func (c Context) SetVariable(name string, val any) {
	var vars Variables

	if _, ok := c.tempVariables.Get(name); ok {
		vars = c.tempVariables
	} else if _, ok := c.globalVariables.Get(name); ok {
		vars = c.globalVariables
	} else {
		vars = c.scopeVariables
	}

	vars.Set(name, val)
}

func (c Context) SetVariables(vars Variables) {
	for key, value := range vars {
		c.SetVariable(key, value)
	}
}

func (c *Context) SetRudispaceFunction(funcName string, fun Function) {
	c.userFuncs = c.userFuncs.Set(funcName, fun)
}

func (c Context) shallowCopy() Context {
	return Context{
		ctx:             c.ctx,
		document:        c.document,
		fixedFuncs:      c.fixedFuncs,
		userFuncs:       c.userFuncs,
		globalVariables: c.globalVariables,
		scopeVariables:  c.scopeVariables,
		tempVariables:   c.tempVariables,
		coalescer:       c.coalescer,
		runtime:         c.runtime,
	}
}

func MakeShim(val any) ast.Shim {
	return ast.Shim{Value: val}
}
