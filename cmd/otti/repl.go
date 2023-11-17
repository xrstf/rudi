// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package main

import (
	"bufio"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"go.xrstf.de/otto/pkg/lang/ast"
	"go.xrstf.de/otto/pkg/lang/eval"
	"go.xrstf.de/otto/pkg/lang/eval/builtin"
	"go.xrstf.de/otto/pkg/lang/eval/types"
	"go.xrstf.de/otto/pkg/lang/parser"
)

//go:embed help.txt
var helpText string

func printPrompt() {
	fmt.Print("⮞ ")
}

func displayHelp() error {
	fmt.Println(helpText)
	return nil
}

func cleanInput(text string) string {
	return strings.TrimSpace(text)
}

type replCommandFunc func() error

var replCommands = map[string]replCommandFunc{
	"help": displayHelp,
}

func replRun(opts *options, args []string) error {
	if len(args) == 0 {
		return errors.New("no input file given")
	}

	document, err := loadDocument(opts, args[0])
	if err != nil {
		return fmt.Errorf("failed to read %q: %w", args[0], err)
	}

	fmt.Println("Welcome to Otti 🐘")
	fmt.Println("Type `help` fore more information, `exit` or Ctrl-C to exit.")
	fmt.Println("")

	reader := bufio.NewScanner(os.Stdin)
	printPrompt()

	for reader.Scan() {
		input := cleanInput(reader.Text())

		stop, err := processReplInput(opts, &document, input)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		if stop {
			break
		}
		printPrompt()
	}

	fmt.Println()

	return nil
}

func processReplInput(opts *options, doc *types.Document, input string) (stop bool, err error) {
	if command, exists := replCommands[input]; exists {
		return false, command()
	}

	if strings.EqualFold("exit", input) {
		return true, nil
	}

	// parse input
	got, err := parser.Parse("(repl)", []byte(input))
	if err != nil {
		return false, err
		// fmt.Println(caretError(err, string(content)))
		// os.Exit(1)
	}

	program, ok := got.(ast.Program)
	if !ok {
		return false, fmt.Errorf("parsed input is not a ast.Program, but %T", got)
	}

	vars := eval.NewVariables()
	ctx := eval.NewContext(*doc, builtin.Functions, vars)

	evaluated, err := eval.Run(ctx, program)
	if err != nil {
		return false, err
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.Encode(evaluated)

	return false, nil
}
