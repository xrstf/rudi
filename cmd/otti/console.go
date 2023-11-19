// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package main

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"go.xrstf.de/otto"
	"go.xrstf.de/otto/docs"
	"go.xrstf.de/otto/pkg/eval/types"
)

//go:embed help.txt
var helpText string

func printPrompt() {
	fmt.Print("⮞ ")
}

func displayHelp(ctx types.Context, helpTopics []docs.Topic, opts *options) error {
	var builder strings.Builder
	builder.WriteString(strings.TrimSpace(helpText))
	builder.WriteString("\n\n")

	width := 0
	for _, topic := range helpTopics {
		if l := len(topic.CliNames[0]); l > width {
			width = l
		}
	}

	format := fmt.Sprintf("* %%-%ds – %%s\n", width)
	for _, topic := range helpTopics {
		builder.WriteString(fmt.Sprintf(format, topic.CliNames[0], topic.Description))
	}

	printMarkdown(builder.String())

	return nil
}

func cleanInput(text string) string {
	return strings.TrimSpace(text)
}

type replCommandFunc func(ctx types.Context, helpTopics []docs.Topic, opts *options) error

var replCommands = map[string]replCommandFunc{
	"help": displayHelp,
}

func runConsole(opts *options, args []string) error {
	rl, err := readline.New("⮞ ")
	if err != nil {
		return fmt.Errorf("failed to setup readline prompt: %w", err)
	}

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

	helpTopics := docs.Topics()

	for {
		line, err := rl.Readline()
		if err != nil { // io.EOF
			break
		}
		line = cleanInput(line)
		if line == "" {
			continue
		}

		newCtx, stop, err := processReplInput(ctx, helpTopics, opts, line)
		if err != nil {
			parseErr := &otto.ParseError{}
			if errors.As(err, parseErr) {
				fmt.Println(parseErr.Snippet())
				fmt.Println(parseErr)
			} else {
				fmt.Printf("Error: %v\n", err)
			}
		}
		if stop {
			break
		}

		ctx = newCtx
	}

	fmt.Println()

	return nil
}

func processReplInput(ctx types.Context, helpTopics []docs.Topic, opts *options, input string) (newCtx types.Context, stop bool, err error) {
	if command, exists := replCommands[input]; exists {
		return ctx, false, command(ctx, helpTopics, opts)
	}

	if prefix := "help "; strings.HasPrefix(input, prefix) {
		topicName := strings.TrimPrefix(input, prefix)
		return ctx, false, renderHelpTopic(helpTopics, topicName)
	}

	if strings.EqualFold("exit", input) {
		return ctx, true, nil
	}

	// parse input
	program, err := otto.ParseScript("(repl)", input)
	if err != nil {
		return ctx, false, err
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
