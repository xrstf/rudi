// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"go.xrstf.de/otto/pkg/lang/eval"
	"go.xrstf.de/otto/pkg/lang/eval/builtin"
	"go.xrstf.de/otto/pkg/lang/eval/types"

	"gopkg.in/yaml.v3"
)

func runScript(opts *options, args []string) error {
	// determine input script to evaluate
	script := ""
	if opts.scriptFile != "" {
		content, err := os.ReadFile(opts.scriptFile)
		if err != nil {
			return fmt.Errorf("failed to read script from %s: %w", opts.scriptFile, err)
		}

		script = strings.TrimSpace(string(content))
	} else {
		if len(args) == 0 {
			return errors.New("no script provided either via argument or --script")
		}

		// consume one arg for the script
		script = args[0]
		args = args[1:]
	}

	// parse the script
	program, err := parseScript(script)
	if err != nil {
		return fmt.Errorf("invalid script: %w", err)
	}

	// load all remaining args as input files
	files, err := loadFiles(opts, args)
	if err != nil {
		return fmt.Errorf("failed to read inputs: %w", err)
	}

	// setup the evaluation context
	document, err := types.NewDocument(files[0])
	if err != nil {
		return fmt.Errorf("cannot use %s as document: %w", args[0], err)
	}

	vars := eval.NewVariables().
		Set("files", files)

	ctx := eval.NewContext(document, builtin.Functions, vars)

	// evaluate the script
	_, evaluated, err := eval.Run(ctx, program)
	if err != nil {
		return fmt.Errorf("failed to evaluate script: %w", err)
	}

	// print the output
	var encoder interface {
		Encode(v any) error
	}

	if opts.formatYaml {
		encoder = yaml.NewEncoder(os.Stdout)
		encoder.(*yaml.Encoder).SetIndent(2)
	} else {
		encoder = json.NewEncoder(os.Stdout)
		if opts.prettyPrint {
			encoder.(*json.Encoder).SetIndent("", "  ")
		}
	}

	if err := encoder.Encode(evaluated); err != nil {
		return fmt.Errorf("failed to encode %v: %w", evaluated, err)
	}

	return nil
}
