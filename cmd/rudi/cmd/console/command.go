// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package console

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go.xrstf.de/rudi"
	"go.xrstf.de/rudi/cmd/rudi/docs"
	cmdtypes "go.xrstf.de/rudi/cmd/rudi/types"
	"go.xrstf.de/rudi/cmd/rudi/util"
	"go.xrstf.de/rudi/pkg/eval/types"

	colorjson "github.com/TylerBrock/colorjson"
	"github.com/chzyer/readline"
)

func helpCommand(ctx types.Context, opts *cmdtypes.Options) error {
	content, err := docs.RenderFile("cmd-console.md", nil)
	if err != nil {
		return err
	}

	fmt.Print(content)

	return nil
}

func helpTopicCommand(topic string) error {
	rendered, err := util.RenderHelpTopic(topic, 2)
	if err != nil {
		return err
	}

	fmt.Print(rendered)

	return nil
}

type replCommandFunc func(ctx types.Context, opts *cmdtypes.Options) error

var replCommands = map[string]replCommandFunc{
	"help": helpCommand,
}

func Run(ctx context.Context, opts *cmdtypes.Options, args []string, rudiVersion string) error {
	rl, err := readline.New("⮞ ")
	if err != nil {
		return fmt.Errorf("failed to setup readline prompt: %w", err)
	}

	fileContents, err := util.LoadFiles(opts, args)
	if err != nil {
		return fmt.Errorf("failed to read inputs: %w", err)
	}

	rudiCtx, err := util.SetupRudiContext(opts, args, fileContents)
	if err != nil {
		return fmt.Errorf("failed to setup context: %w", err)
	}

	fmt.Printf("Welcome to 🚂Rudi %s\n", rudiVersion)
	fmt.Println("Type `help` fore more information, `exit` or Ctrl-C to exit.")
	fmt.Println("")

	for {
		line, err := rl.Readline()
		if err != nil { // io.EOF
			break
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		newCtx, stop, err := processInput(ctx, rudiCtx, opts, line)
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

		rudiCtx = newCtx
	}

	fmt.Println()

	return nil
}

func processInput(ctx context.Context, rudiCtx types.Context, opts *cmdtypes.Options, input string) (newCtx types.Context, stop bool, err error) {
	if command, exists := replCommands[input]; exists {
		return rudiCtx, false, command(rudiCtx, opts)
	}

	if prefix := "help "; strings.HasPrefix(input, prefix) {
		topicName := strings.TrimPrefix(input, prefix)
		return rudiCtx, false, helpTopicCommand(topicName)
	}

	if strings.EqualFold("exit", input) {
		return rudiCtx, true, nil
	}

	// parse input
	program, err := rudi.Parse("(repl)", input)
	if err != nil {
		return rudiCtx, false, err
	}

	// TODO: Setup a new context that can be cancelled with Ctrl-C to interrupt long running statements.

	newCtx, evaluated, err := program.RunContext(rudiCtx.WithGoContext(ctx))
	if err != nil {
		return rudiCtx, false, err
	}

	f := colorjson.NewFormatter()
	f.Indent = 0
	f.EscapeHTML = false

	encoded, err := f.Marshal(evaluated)
	if err != nil {
		return rudiCtx, false, fmt.Errorf("failed to encode %v: %w", evaluated, err)
	}

	fmt.Println(string(encoded))

	return newCtx, false, nil
}
