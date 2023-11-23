// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package script

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"go.xrstf.de/rudi"
	"go.xrstf.de/rudi/cmd/rudi/types"
	"go.xrstf.de/rudi/cmd/rudi/util"
	"go.xrstf.de/rudi/pkg/debug"

	"gopkg.in/yaml.v3"
)

func Run(opts *types.Options, args []string) error {
	// determine input script to evaluate
	script := ""
	scriptName := ""

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
		if err := debug.Dump(program, os.Stdout); err != nil {
			return fmt.Errorf("failed to dump AST: %w", err)
		}

		return nil
	}

	// load all remaining args as input files
	files, err := util.LoadFiles(opts, args)
	if err != nil {
		return fmt.Errorf("failed to read inputs: %w", err)
	}

	// setup the evaluation context
	ctx, err := util.SetupRudiContext(files)
	if err != nil {
		return fmt.Errorf("failed to setup context: %w", err)
	}

	// evaluate the script
	_, evaluated, err := program.RunContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to evaluate script: %w", err)
	}

	// print the output
	var encoder interface {
		Encode(v any) error
	}

	if opts.FormatYaml {
		encoder = yaml.NewEncoder(os.Stdout)
		encoder.(*yaml.Encoder).SetIndent(2)
	} else {
		encoder = json.NewEncoder(os.Stdout)
		if opts.PrettyPrint {
			encoder.(*json.Encoder).SetIndent("", "  ")
		}
	}

	if err := encoder.Encode(evaluated); err != nil {
		return fmt.Errorf("failed to encode %v: %w", evaluated, err)
	}

	return nil
}
