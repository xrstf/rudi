// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"go.xrstf.de/otto/pkg/lang/ast"
	"go.xrstf.de/otto/pkg/lang/parser"

	"gopkg.in/yaml.v3"
)

func loadFiles(opts *options, filenames []string) ([]any, error) {
	results := make([]any, len(filenames))

	for i, filename := range filenames {
		data, err := loadFile(opts, filename)
		if err != nil {
			return nil, fmt.Errorf("failed to read %s: %w", filename, err)
		}

		results[i] = data
	}

	return results, nil
}

func loadFile(opts *options, filename string) (any, error) {
	if filename == "" {
		return nil, errors.New("no filename provided")
	}

	var input io.Reader

	if filename == "-" {
		input = os.Stdin
	} else {
		f, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		input = f
	}

	var doc any

	decoder := yaml.NewDecoder(input)
	if err := decoder.Decode(&doc); err != nil {
		return nil, fmt.Errorf("failed to parse document as YAML/JSON: %w", err)
	}

	return doc, nil
}

func parseScript(script string) (ast.Program, error) {
	got, err := parser.Parse("(repl)", []byte(script))
	if err != nil {
		return ast.Program{}, err
		// fmt.Println(caretError(err, script))
		// os.Exit(1)
	}

	program, ok := got.(ast.Program)
	if !ok {
		// this should never happen
		return ast.Program{}, fmt.Errorf("parsed input is not a ast.Program, but %T", got)
	}

	return program, nil
}
