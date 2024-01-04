// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package script

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"go.xrstf.de/rudi"
	"go.xrstf.de/rudi/cmd/rudi/encoding"
	"go.xrstf.de/rudi/cmd/rudi/options"
	"go.xrstf.de/rudi/cmd/rudi/util"
)

func Run(handler *util.SignalHandler, opts *options.Options, args []string) error {
	// determine input script to evaluate
	var (
		script     string
		scriptName string
	)

	if opts.ScriptFile != "" {
		content, err := os.ReadFile(opts.ScriptFile)
		if err != nil {
			return fmt.Errorf("failed to read script from %s: %w", opts.ScriptFile, err)
		}

		script = strings.TrimSpace(string(content))
		scriptName = opts.ScriptFile
	} else {
		if len(args) == 0 {
			return errors.New("no script provided either via argument or --script")
		}

		// consume one arg for the script
		script = args[0]
		args = args[1:]
		scriptName = "(stdin)"
	}

	// parse the script
	program, err := rudi.Parse(scriptName, script)
	if err != nil {
		return fmt.Errorf("invalid script: %w", err)
	}

	// show AST and quit if desired
	if opts.PrintAst {
		if err := program.DumpSyntaxTree(os.Stdout); err != nil {
			return fmt.Errorf("failed to dump AST: %w", err)
		}

		return nil
	}

	// load all remaining args as input fileContents
	fileContents, err := util.LoadFiles(opts, args)
	if err != nil {
		return fmt.Errorf("failed to read inputs: %w", err)
	}

	// setup the evaluation context
	rudiCtx, err := util.SetupRudiContext(opts, args, fileContents)
	if err != nil {
		return fmt.Errorf("failed to setup context: %w", err)
	}

	// allow to interrupt the script
	subCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	handler.SetCancelFn(cancel)

	// evaluate the script
	evaluated, err := program.RunContext(rudiCtx.WithGoContext(subCtx))
	if err != nil {
		return fmt.Errorf("failed to evaluate script: %w", err)
	}

	// print the output
	if err := encoding.Encode(evaluated, opts.OutputFormat, os.Stdout); err != nil {
		return fmt.Errorf("failed to encode data: %w", err)
	}

	return nil
}
