// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package types

import (
	"errors"

	"github.com/spf13/pflag"
)

// These constants must be lowercased because the validation function normalizes the given user
// input to lowercase.

type Encoding string

func (f Encoding) String() string {
	return string(f)
}

const (
	JsonEncoding Encoding = "json"
	YamlEncoding Encoding = "yaml"
	TomlEncoding Encoding = "toml"
)

type Coalescer string

func (c Coalescer) String() string {
	return string(c)
}

const (
	StrictCoalescer   Coalescer = "strict"
	PedanticCoalescer Coalescer = "pedantic"
	HumaneCoalescer   Coalescer = "humane"
)

type Options struct {
	ShowHelp                 bool
	Interactive              bool
	ScriptFile               string
	StdinFormat              Encoding
	OutputFormat             Encoding
	PrintAst                 bool
	ShowVersion              bool
	Coalescing               Coalescer
	EnableRudispaceFunctions bool
}

func NewDefaultOptions() Options {
	return Options{
		Coalescing:  StrictCoalescer,
		StdinFormat: YamlEncoding,
	}
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.SortFlags = false

	stdinFormatFlag := newEnumFlag(&o.StdinFormat, JsonEncoding, YamlEncoding, TomlEncoding)
	outputFormatFlag := newEnumFlag(&o.OutputFormat, JsonEncoding, YamlEncoding, TomlEncoding)
	coalescingFlag := newEnumFlag(&o.Coalescing, StrictCoalescer, PedanticCoalescer, HumaneCoalescer)

	fs.BoolVarP(&o.Interactive, "interactive", "i", o.Interactive, "Start an interactive REPL to run expressions.")
	fs.StringVarP(&o.ScriptFile, "script", "s", o.ScriptFile, "Load Rudi script from file instead of first argument (only in non-interactive mode).")
	stdinFormatFlag.Add(fs, "stdin-format", "f", "What data format is used for data provided on stdin")
	outputFormatFlag.Add(fs, "output-format", "o", "What data format to use for outputting data (if not given, unformatted JSON is used)")
	fs.BoolVar(&o.EnableRudispaceFunctions, "enable-funcs", o.EnableRudispaceFunctions, "Enable the func! function to allow defining new functions in Rudi code.")
	coalescingFlag.Add(fs, "coalesce", "c", "Type conversion handling")
	fs.BoolVarP(&o.ShowHelp, "help", "h", o.ShowHelp, "Show help and documentation.")
	fs.BoolVarP(&o.ShowVersion, "version", "V", o.ShowVersion, "Show version and exit.")
	fs.BoolVarP(&o.PrintAst, "debug-ast", "", o.PrintAst, "Output syntax tree of the parsed script in non-interactive mode.")
}

func (o *Options) Validate() error {
	if o.Interactive && o.ScriptFile != "" {
		return errors.New("cannot combine --interactive with --script")
	}

	if o.Interactive && o.PrintAst {
		return errors.New("cannot combine --interactive with --debug-ast")
	}

	return nil
}
