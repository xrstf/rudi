// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package console

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"go.xrstf.de/rudi"
	cmdtypes "go.xrstf.de/rudi/cmd/rudi/types"
	"go.xrstf.de/rudi/cmd/rudi/util"
	"go.xrstf.de/rudi/docs"
	"go.xrstf.de/rudi/pkg/eval/types"

	"github.com/chzyer/readline"
)

//go:embed help.md
var helpText string

func printPrompt() {
	fmt.Print("⮞ ")
}

func helpCommand(ctx types.Context, helpTopics []docs.Topic, opts *cmdtypes.Options) error {
	fmt.Println(util.RenderMarkdown(strings.TrimSpace(helpText), 2))
	fmt.Println(util.RenderHelpTopics(helpTopics, 2))

	return nil
}

func helpTopicCommand(helpTopics []docs.Topic, topic string) error {
	rendered, err := util.RenderHelpTopic(helpTopics, topic, 2)
	if err != nil {
		return err
	}

	fmt.Println(rendered)

	return nil
}

type replCommandFunc func(ctx types.Context, helpTopics []docs.Topic, opts *cmdtypes.Options) error

var replCommands = map[string]replCommandFunc{
	"help": helpCommand,
}

func Run(opts *cmdtypes.Options, args []string) error {
	rl, err := readline.New("⮞ ")
	if err != nil {
		return fmt.Errorf("failed to setup readline prompt: %w", err)
	}

	files, err := util.LoadFiles(opts, args)
	if err != nil {
		return fmt.Errorf("failed to read inputs: %w", err)
	}

	ctx, err := util.SetupRudiContext(files)
	if err != nil {
		return fmt.Errorf("failed to setup context: %w", err)
	}

	fmt.Println("Welcome to Rudi 🐘")
	fmt.Println("Type `help` fore more information, `exit` or Ctrl-C to exit.")
	fmt.Println("")

	helpTopics := docs.Topics()

	for {
		line, err := rl.Readline()
		if err != nil { // io.EOF
			break
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		newCtx, stop, err := processInput(ctx, helpTopics, opts, line)
		if err != nil {
			parseErr := &rudi.ParseError{}
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

func processInput(ctx types.Context, helpTopics []docs.Topic, opts *cmdtypes.Options, input string) (newCtx types.Context, stop bool, err error) {
	if command, exists := replCommands[input]; exists {
		return ctx, false, command(ctx, helpTopics, opts)
	}

	if prefix := "help "; strings.HasPrefix(input, prefix) {
		topicName := strings.TrimPrefix(input, prefix)
		return ctx, false, helpTopicCommand(helpTopics, topicName)
	}

	if strings.EqualFold("exit", input) {
		return ctx, true, nil
	}

	// parse input
	program, err := rudi.ParseScript("(repl)", input)
	if err != nil {
		return ctx, false, err
	}

	newCtx, evaluated, err := rudi.RunProgram(ctx, program)
	if err != nil {
		return ctx, false, err
	}

	encoder := json.NewEncoder(os.Stdout)
	if err := encoder.Encode(evaluated); err != nil {
		return ctx, false, fmt.Errorf("failed to encode %v: %w", evaluated, err)
	}

	return newCtx, false, nil
}
