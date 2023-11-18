// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package main

import (
	"bufio"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"go.xrstf.de/otto/pkg/lang/eval"
	"go.xrstf.de/otto/pkg/lang/eval/builtin"
	"go.xrstf.de/otto/pkg/lang/eval/types"
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

func runConsole(opts *options, args []string) error {
	files, err := loadFiles(opts, args)
	if err != nil {
		return fmt.Errorf("failed to read inputs: %w", err)
	}

	var document types.Document

	if len(files) > 0 {
		document, err = types.NewDocument(files[0])
		if err != nil {
			return fmt.Errorf("cannot use %s as document: %w", args[0], err)
		}
	} else {
		document, _ = types.NewDocument(nil)
	}

	vars := eval.NewVariables().
		Set("files", files)

	ctx := eval.NewContext(document, builtin.Functions, vars)

	fmt.Println("Welcome to Otti 🐘")
	fmt.Println("Type `help` fore more information, `exit` or Ctrl-C to exit.")
	fmt.Println("")

	reader := bufio.NewScanner(os.Stdin)
	printPrompt()

	for reader.Scan() {
		input := cleanInput(reader.Text())

		newCtx, stop, err := processReplInput(ctx, opts, &document, input)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		if stop {
			break
		}
		printPrompt()

		ctx = newCtx
	}

	fmt.Println()

	return nil
}

func processReplInput(ctx types.Context, opts *options, doc *types.Document, input string) (newCtx types.Context, stop bool, err error) {
	if command, exists := replCommands[input]; exists {
		return ctx, false, command()
	}

	if strings.EqualFold("exit", input) {
		return ctx, true, nil
	}

	// parse input
	program, err := parseScript(input)
	if err != nil {
		return ctx, false, err
		// fmt.Println(caretError(err, string(content)))
		// os.Exit(1)
	}

	newCtx, evaluated, err := eval.Run(ctx, program)
	if err != nil {
		return ctx, false, err
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.Encode(evaluated)

	return newCtx, false, nil
}
