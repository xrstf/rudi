// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package types

import (
	"errors"

	"github.com/spf13/pflag"
)

type Options struct {
	ShowHelp    bool
	Interactive bool
	ScriptFile  string
	PrettyPrint bool
	FormatYaml  bool
	PrintAst    bool
	ShowVersion bool
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.BoolVarP(&o.ShowHelp, "help", "h", o.ShowHelp, "Show help and documentation.")
	fs.BoolVarP(&o.Interactive, "interactive", "i", o.Interactive, "Start an interactive REPL to run expressions.")
	fs.StringVarP(&o.ScriptFile, "script", "s", o.ScriptFile, "Load Otto script from file instead of first argument (only in non-interactive mode).")
	fs.BoolVarP(&o.PrettyPrint, "pretty", "p", o.PrettyPrint, "Output pretty-printed JSON.")
	fs.BoolVarP(&o.FormatYaml, "yaml", "y", o.FormatYaml, "Output pretty-printed YAML instead of JSON.")
	fs.BoolVarP(&o.PrintAst, "debug-ast", "", o.PrintAst, "Output syntax tree of the parsed script in non-interactive mode.")
	fs.BoolVarP(&o.ShowVersion, "version", "V", o.ShowVersion, "Show version and exit.")
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
