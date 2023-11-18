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

	"go.xrstf.de/otto"
	"go.xrstf.de/otto/pkg/eval/types"
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

	ctx, err := setupOttoContext(files)
	if err != nil {
		return fmt.Errorf("failed to setup context: %w", err)
	}

	fmt.Println("Welcome to Otti 🐘")
	fmt.Println("Type `help` fore more information, `exit` or Ctrl-C to exit.")
	fmt.Println("")

	reader := bufio.NewScanner(os.Stdin)
	printPrompt()

	for reader.Scan() {
		input := cleanInput(reader.Text())

		newCtx, stop, err := processReplInput(ctx, opts, input)
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

func processReplInput(ctx types.Context, opts *options, input string) (newCtx types.Context, stop bool, err error) {
	if command, exists := replCommands[input]; exists {
		return ctx, false, command()
	}

	if strings.EqualFold("exit", input) {
		return ctx, true, nil
	}

	// parse input
	program, err := otto.ParseScript("(repl)", input)
	if err != nil {
		return ctx, false, err
		// fmt.Println(caretError(err, string(content)))
		// os.Exit(1)
	}

	newCtx, evaluated, err := otto.RunProgram(ctx, program)
	if err != nil {
		return ctx, false, err
	}

	encoder := json.NewEncoder(os.Stdout)
	if err := encoder.Encode(evaluated); err != nil {
		return ctx, false, fmt.Errorf("failed to encode %v: %w", evaluated, err)
	}

	return newCtx, false, nil
}
