// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package rudi

import (
	"context"
	"fmt"
	"io"

	"go.xrstf.de/rudi/pkg/lang/ast"
	"go.xrstf.de/rudi/pkg/lang/parser"
	"go.xrstf.de/rudi/pkg/printer"
	"go.xrstf.de/rudi/pkg/runtime/interpreter"
)

// Program is a parsed Rudi program, ready to be run (executed). Programs are
// stateless and can be executed multiple times, even concurrently (as long
// as a different context is used per goroutine, when using RunContext).
type Program interface {
	fmt.Stringer

	// Run will evaluate the program. The given data value is used as the program's
	// document (i.e. available with bare path expressions like `.foo`). Variables
	// can be left empty if desired, but funcs must effectively always be set,
	// as programs without functions are very limited. Use NewBuiltInFunctions()
	// to get the default set of functions in Rudi.
	// When no error occurs, Run() returns both the final document value and the
	// result of the final expression. Otherwise an error is returned.
	Run(ctx context.Context, data any, variables Variables, funcs Functions, coalescer Coalescer) (document any, result any, err error)

	// RunContext is like Run(), but uses a pre-setup Context and returns the
	// bare final context instead of its document's value. The result is still
	// the result of the final expression in the program.
	RunContext(ctx Context) (result any, err error)

	// DumpSyntaxTree writes the AST to the given writer. Useful for debugging.
	// Set indent to false to prevent multiline output from being generated
	// according to a simple, conservative linebreak algorithm.
	// Note that the output looks like code, but is not executable/parseable. Use
	// DumpRudi() if you need to turn a parsed program back into Rudi code.
	DumpSyntaxTree(out io.Writer) error

	// DumpRudi writes the AST in the form of parseable Rudi code.
	DumpRudi(out io.Writer) error
}

type rudiProgram struct {
	prog *ast.Program
}

// Parse takes a program name and a script and returns a parsed Program or an
// error if parsing failed. The program can can be any string, but is often the
// filename from where the script was loaded.
func Parse(name, script string) (Program, error) {
	got, err := parser.Parse(name, []byte(script))
	if err != nil {
		return nil, ParseError{script: script, err: err}
	}

	program, ok := got.(ast.Program)
	if !ok {
		// this should never happen
		return nil, fmt.Errorf("parsed input is not an ast.Program, but %T", got)
	}

	return &rudiProgram{
		prog: &program,
	}, nil
}

// Run will evaluate the program. The given data value is used as the program's
// document (i.e. available with bare path expressions like `.foo`). Variables
// can be left empty if desired, but funcs must effectively always be set,
// as programs without functions are very limited. Use NewBuiltInFunctions()
// to get the default set of functions in Rudi.
// When no error occurs, Run() returns both the final document value and the
// result of the final expression. Otherwise an error is returned.
func (p *rudiProgram) Run(ctx context.Context, data any, variables Variables, funcs Functions, coalescer Coalescer) (document any, result any, err error) {
	doc, err := NewDocument(data)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot process %T: %w", data, err)
	}

	rudiCtx, err := NewContext(interpreter.New(), ctx, doc, variables, funcs, coalescer)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create context: %w", err)
	}

	result, err = p.RunContext(rudiCtx)
	if err != nil {
		return nil, nil, fmt.Errorf("script failed: %w", err)
	}

	// get current state of the document
	docData := rudiCtx.GetDocument().Data()

	return docData, result, nil
}

// RunContext is like Run(), but uses a pre-setup Context and returns the
// bare final context instead of its document's value. The result is still
// the result of the final expression in the program.
func (p *rudiProgram) RunContext(ctx Context) (result any, err error) {
	return ctx.Runtime().EvalProgram(ctx, p.prog)
}

// String returns the Rudi-representation of the parsed script, with comments
// removed.
func (p *rudiProgram) String() string {
	return p.prog.String()
}

// DumpSyntaxTree writes the AST to the given writer. Useful for debugging.
// Set indent to false to prevent multiline output from being generated
// according to a simple, conservative linebreak algorithm.
// Note that the output looks like code, but is not executable/parseable. Use
// DumpRudi() if you need to turn a parsed program back into Rudi code.
func (p *rudiProgram) DumpSyntaxTree(out io.Writer) error {
	return printer.NewAstPrinter(out).Program(p.prog)
}

// DumpRudi writes the AST in the form of parseable Rudi code.
func (p *rudiProgram) DumpRudi(out io.Writer) error {
	return printer.NewRudiPrinter(out).Program(p.prog)
}
