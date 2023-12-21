// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package types

import (
	"context"

	"go.xrstf.de/rudi/pkg/coalescing"
	"go.xrstf.de/rudi/pkg/lang/ast"
)

type Document struct {
	data any
}

func NewDocument(data any) (Document, error) {
	return Document{
		data: data,
	}, nil
}

func (d *Document) Data() any {
	return d.data
}

func (d *Document) Set(wrappedData any) {
	d.data = wrappedData
}

type Context struct {
	ctx        context.Context
	document   *Document
	fixedFuncs Functions
	userFuncs  Functions
	variables  Variables
	coalescer  coalescing.Coalescer
}

func NewContext(ctx context.Context, doc Document, variables Variables, funcs Functions, coalescer coalescing.Coalescer) Context {
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
	}
}

// Coalesce is named this way to make the frequent calls read fluently
// (for example "ctx.Coalesce().ToBool(...)").
func (c Context) Coalesce() coalescing.Coalescer {
	return c.coalescer
}

func (c Context) GoContext() context.Context {
	return c.ctx
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
	return Context{
		ctx:        ctx,
		document:   c.document,
		fixedFuncs: c.fixedFuncs,
		userFuncs:  c.userFuncs,
		variables:  c.variables,
		coalescer:  c.coalescer,
	}
}

func (c Context) WithVariable(name string, val any) Context {
	return Context{
		ctx:        c.ctx,
		document:   c.document,
		fixedFuncs: c.fixedFuncs,
		userFuncs:  c.userFuncs,
		variables:  c.variables.With(name, val),
		coalescer:  c.coalescer,
	}
}

func (c Context) WithCoalescer(coalescer coalescing.Coalescer) Context {
	return Context{
		ctx:        c.ctx,
		document:   c.document,
		fixedFuncs: c.fixedFuncs,
		userFuncs:  c.userFuncs,
		variables:  c.variables,
		coalescer:  coalescer,
	}
}

func (c Context) WithRudispaceFunction(funcName string, fun Function) Context {
	return Context{
		ctx:        c.ctx,
		document:   c.document,
		fixedFuncs: c.fixedFuncs,
		userFuncs:  c.userFuncs.DeepCopy().Set(funcName, fun),
		variables:  c.variables,
		coalescer:  c.coalescer,
	}
}

type Function interface {
	Evaluate(ctx Context, args []ast.Expression) (any, error)

	// Description returns a short, one-line description of the function; markdown
	// can be used to highlight other function names, like "behaves similar
	// to `foo`, but …".
	Description() string
}

type TupleFunction func(ctx Context, args []ast.Expression) (any, error)

type basicFunc struct {
	f    TupleFunction
	desc string
}

// NewFunction creates the lowest of low level functions in Rudi and should rarely be used by
// integrators/library developers. Use the helpers in the root package to define functions
// using reflection and pattern matching instead.
func NewFunction(f TupleFunction, description string) Function {
	return basicFunc{
		f:    f,
		desc: description,
	}
}

var _ Function = basicFunc{}

func (f basicFunc) Evaluate(ctx Context, args []ast.Expression) (any, error) {
	return f.f(ctx, args)
}

func (f basicFunc) Description() string {
	return f.desc
}

type Functions map[string]Function

func NewFunctions() Functions {
	return Functions{}
}

func (f Functions) Get(name string) (Function, bool) {
	variable, exists := f[name]
	return variable, exists
}

// Set sets/replaces the function in the current set (in-place).
// The function returns the same Functions to allow fluent access.
func (f Functions) Set(name string, fun Function) Functions {
	f[name] = fun
	return f
}

// Set removes a function from the set.
// The function returns the same Functions to allow fluent access.
func (f Functions) Delete(name string) Functions {
	delete(f, name)
	return f
}

// Add adds all functions from other to the current set.
// The function returns the same Functions to allow fluent access.
func (f Functions) Add(other Functions) Functions {
	for name, fun := range other {
		f[name] = fun
	}
	return f
}

// Remove removes all functions from this set that are part of the other set,
// to enable constructs like AllFunctions.Remove(MathFunctions)
// The function returns the same Functions to allow fluent access.
func (f Functions) Remove(other Functions) Functions {
	for name := range other {
		f.Delete(name)
	}
	return f
}

func (f Functions) DeepCopy() Functions {
	result := NewFunctions()
	for key, val := range f {
		result[key] = val
	}
	return result
}

type Variables map[string]any

func NewVariables() Variables {
	return Variables{}
}

func (v Variables) Get(name string) (any, bool) {
	variable, exists := v[name]
	return variable, exists
}

// Set sets/replaces the variable value in the current set (in-place).
// The function returns the same variables to allow fluent access.
func (v Variables) Set(name string, val any) Variables {
	v[name] = val
	return v
}

// With returns a copy of the variables, with the new variable being added to it.
func (v Variables) With(name string, val any) Variables {
	return v.DeepCopy().Set(name, val)
}

func (v Variables) DeepCopy() Variables {
	result := NewVariables()
	for key, val := range v {
		result[key] = val
	}
	return result
}

func MakeShim(val any) ast.Shim {
	return ast.Shim{Value: val}
}
