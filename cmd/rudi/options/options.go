// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package options

import (
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"go.xrstf.de/rudi/cmd/rudi/encoding"
	"go.xrstf.de/rudi/cmd/rudi/types"

	"github.com/spf13/pflag"
)

type Options struct {
	ShowHelp                 bool
	Interactive              bool
	ScriptFile               string
	LibraryFiles             []string
	StdinFormat              types.Encoding
	OutputFormat             types.Encoding
	PrintAst                 bool
	ShowVersion              bool
	Coalescing               types.Coalescing
	EnableRudispaceFunctions bool
	ExtraVariables           map[string]any
	extraVariableFlags       []string
}

func NewDefaultOptions() Options {
	return Options{
		Coalescing:     types.StrictCoalescing,
		StdinFormat:    types.YamlEncoding,
		OutputFormat:   types.JsonEncoding,
		ExtraVariables: map[string]any{},
	}
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.SortFlags = false

	stdinFormatFlag := newEnumFlag(&o.StdinFormat, types.InputEncodings...)
	outputFormatFlag := newEnumFlag(&o.OutputFormat, types.OutputEncodings...)
	coalescingFlag := newEnumFlag(&o.Coalescing, types.AllCoalescings...)

	fs.BoolVarP(&o.Interactive, "interactive", "i", o.Interactive, "Start an interactive REPL to run expressions.")
	fs.StringVarP(&o.ScriptFile, "script", "s", o.ScriptFile, "Load Rudi script from file instead of first argument (only in non-interactive mode).")
	fs.StringArrayVarP(&o.LibraryFiles, "library", "l", o.LibraryFiles, "Load additional Rudi file(s) to be be evaluated before the script (can be given multiple times).")
	fs.StringArrayVar(&o.extraVariableFlags, "var", o.extraVariableFlags, "Define additional global variables (can be given multiple times).")
	stdinFormatFlag.Add(fs, "stdin-format", "f", "What data format is used for data provided on stdin")
	outputFormatFlag.Add(fs, "output-format", "o", "What data format to use for outputting data")
	fs.BoolVar(&o.EnableRudispaceFunctions, "enable-funcs", o.EnableRudispaceFunctions, "Enable the func! function to allow defining new functions in Rudi code.")
	coalescingFlag.Add(fs, "coalesce", "c", "Type conversion handling")
	fs.BoolVarP(&o.ShowHelp, "help", "h", o.ShowHelp, "Show help and documentation.")
	fs.BoolVarP(&o.ShowVersion, "version", "V", o.ShowVersion, "Show version and exit.")
	fs.BoolVarP(&o.PrintAst, "debug-ast", "", o.PrintAst, "Output syntax tree of the parsed script in non-interactive mode.")
}

func (o *Options) Validate() error {
	if o.Interactive && o.PrintAst {
		return errors.New("cannot combine --interactive with --debug-ast")
	}

	if err := o.parseExtraVariables(); err != nil {
		return fmt.Errorf("invalid --var flags: %w", err)
	}

	if err := o.validateLibraryFiles(); err != nil {
		return fmt.Errorf("invalid --library flags: %w", err)
	}

	return nil
}

var extraVariableFlagFormat = regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9_]*)=([a-z]+):([a-z]+):(.+)$`)

func (o *Options) parseExtraVariables() error {
	for i, flagValue := range o.extraVariableFlags {
		varName, value, err := o.parseExtraVariable(flagValue)
		if err != nil {
			return fmt.Errorf("--var flag %d: %w", i, err)
		}

		o.ExtraVariables[varName] = value
	}

	return nil
}

func (o *Options) parseExtraVariable(flagValue string) (string, any, error) {
	flagValue = strings.TrimSpace(flagValue)

	match := extraVariableFlagFormat.FindStringSubmatch(flagValue)
	if match == nil {
		return "", nil, errors.New("must be in the form of \"varname=encoding:source:data\"")
	}

	varName := match[1]
	enc := types.Encoding(match[2])
	source := types.VariableSource(match[3])
	data := match[4]

	// validate the given parameters for this variable

	if _, exists := o.ExtraVariables[varName]; exists {
		return "", nil, fmt.Errorf("variable $%s is defined multiple times", varName)
	}

	if !enc.IsValid() {
		return "", nil, fmt.Errorf("invalid encoding %q, must be one of %v", enc, types.AllEncodings)
	}

	if !source.IsValid() {
		return "", nil, fmt.Errorf("invalid source type %q, must be one of %v", source, types.AllVariableSources)
	}

	// resolve the variable source

	var input io.Reader

	switch source {
	case types.StringVariableSource:
		input = strings.NewReader(data)
	case types.EnvironmentVariableSource:
		input = strings.NewReader(os.Getenv(data))
	case types.FileVariableSource:
		f, err := os.Open(data)
		if err != nil {
			return "", nil, fmt.Errorf("failed to open %q: %w", data, err)
		}
		defer f.Close()

		input = f
	default:
		// This should never happen.
		return "", nil, fmt.Errorf("unknown source type %q", source)
	}

	// parse the data as requested

	varData, err := encoding.Decode(input, enc)
	if err != nil {
		return "", nil, fmt.Errorf("failed to decode data: %w", err)
	}

	return varName, varData, nil
}

func (o *Options) validateLibraryFiles() error {
	for _, file := range o.LibraryFiles {
		info, err := os.Stat(file)
		if err != nil {
			return fmt.Errorf("invalid library %q: %w", file, err)
		}
		if info.IsDir() {
			return fmt.Errorf("invalid library %q: is directory", file)
		}
	}

	return nil
}
