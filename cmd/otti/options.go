// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package main

import (
	"errors"

	"github.com/spf13/pflag"
)

type options struct {
	interactive bool
	scriptFile  string
	prettyPrint bool
	formatYaml  bool
	printAst    bool
	version     bool
}

func (o *options) AddFlags(fs *pflag.FlagSet) {
	fs.BoolVarP(&o.interactive, "interactive", "i", o.interactive, "Start an interactive REPL to run expressions.")
	fs.StringVarP(&o.scriptFile, "script", "s", o.scriptFile, "Load Otto script from file instead of first argument (only in non-interactive mode).")
	fs.BoolVarP(&o.prettyPrint, "pretty", "p", o.prettyPrint, "Output pretty-printed JSON.")
	fs.BoolVarP(&o.formatYaml, "yaml", "y", o.formatYaml, "Output pretty-printed YAML instead of JSON.")
	fs.BoolVarP(&o.printAst, "debug-ast", "", o.printAst, "Output syntax tree of the parsed script in non-interactive mode.")
	fs.BoolVarP(&o.version, "version", "V", o.version, "Show version and exit.")
}

func (o *options) Validate() error {
	if o.interactive && o.scriptFile != "" {
		return errors.New("cannot combine --interactive with --script")
	}

	if o.interactive && o.printAst {
		return errors.New("cannot combine --interactive with --debug-ast")
	}

	return nil
}
