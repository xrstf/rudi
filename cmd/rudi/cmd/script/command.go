// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package script

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"go.xrstf.de/rudi"
	"go.xrstf.de/rudi/cmd/rudi/options"
	"go.xrstf.de/rudi/cmd/rudi/types"
	"go.xrstf.de/rudi/cmd/rudi/util"
	"go.xrstf.de/rudi/pkg/printer"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
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
		renderer := printer.AST{}
		if err := renderer.WriteMultiline(program, os.Stdout); err != nil {
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
	_, evaluated, err := program.RunContext(rudiCtx.WithGoContext(subCtx))
	if err != nil {
		return fmt.Errorf("failed to evaluate script: %w", err)
	}

	// print the output
	var encoder interface {
		Encode(v any) error
	}

	switch opts.OutputFormat {
	case types.JsonEncoding:
		encoder = json.NewEncoder(os.Stdout)
		encoder.(*json.Encoder).SetIndent("", "  ")
	case types.YamlEncoding:
		encoder = yaml.NewEncoder(os.Stdout)
		encoder.(*yaml.Encoder).SetIndent(2)
	case types.TomlEncoding:
		encoder = toml.NewEncoder(os.Stdout)
		encoder.(*toml.Encoder).Indent = "  "
	default:
		encoder = &rawEncoder{}
	}

	if err := encoder.Encode(evaluated); err != nil {
		return fmt.Errorf("failed to encode %v: %w", evaluated, err)
	}

	return nil
}

type rawEncoder struct{}

func (e *rawEncoder) Encode(v any) error {
	fmt.Println(v)
	return nil
}
