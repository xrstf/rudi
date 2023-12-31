// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package console

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"go.xrstf.de/rudi"
	"go.xrstf.de/rudi/cmd/rudi/docs"
	"go.xrstf.de/rudi/cmd/rudi/options"
	"go.xrstf.de/rudi/cmd/rudi/util"
	"go.xrstf.de/rudi/pkg/runtime/types"

	colorjson "github.com/TylerBrock/colorjson"
	"github.com/chzyer/readline"
)

func helpCommand(ctx types.Context, opts *options.Options) error {
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

type replCommandFunc func(ctx types.Context, opts *options.Options) error

var replCommands = map[string]replCommandFunc{
	"help": helpCommand,
}

func Run(handler *util.SignalHandler, opts *options.Options, library rudi.Program, args []string, rudiVersion string) error {
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
	fmt.Println("Type `help` for more information, `exit` or Ctrl-D to exit, Ctrl-C to interrupt statements.")
	fmt.Println("")

	// Evaluate the library (its return value is irrelevant, as the main program has to have
	// at least 1 statement, which will overwrite the total return value anyway).
	if library != nil {
		_, err = runProgram(handler, rudiCtx, library)
		if err != nil {
			return fmt.Errorf("failed to evaluate library: %w", err)
		}
	}

	for {
		line, err := rl.Readline()

		// treat interrupts as "clear input"
		if errors.Is(err, readline.ErrInterrupt) {
			continue
		}

		// io.EOF (Ctrl-D)
		if err != nil {
			break
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		stop, err := processInput(handler, rudiCtx, opts, line)
		if err != nil {
			parseErr := &rudi.ParseError{}
			if errors.As(err, parseErr) {
				fmt.Fprintln(os.Stderr, parseErr.Snippet())
				fmt.Fprintln(os.Stderr, parseErr)
			} else {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			}
		}
		if stop {
			break
		}
	}

	fmt.Println()

	return nil
}

func processInput(handler *util.SignalHandler, rudiCtx types.Context, opts *options.Options, input string) (stop bool, err error) {
	if command, exists := replCommands[input]; exists {
		return false, command(rudiCtx, opts)
	}

	if prefix := "help "; strings.HasPrefix(input, prefix) {
		topicName := strings.TrimPrefix(input, prefix)
		return false, helpTopicCommand(topicName)
	}

	if strings.EqualFold("exit", input) {
		return true, nil
	}

	// parse input
	program, err := rudi.Parse("(repl)", input)
	if err != nil {
		return false, err
	}

	// run the program
	evaluated, err := runProgram(handler, rudiCtx, program)
	if err != nil {
		return false, err
	}

	f := colorjson.NewFormatter()
	f.Indent = 0
	f.EscapeHTML = false

	encoded, err := f.Marshal(evaluated)
	if err != nil {
		return false, fmt.Errorf("failed to encode %v: %w", evaluated, err)
	}

	fmt.Println(string(encoded))

	return false, nil
}

func runProgram(handler *util.SignalHandler, rudiCtx types.Context, prog rudi.Program) (any, error) {
	// allow to interrupt the statement
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	handler.SetCancelFn(cancel)

	return prog.RunContext(rudiCtx.WithGoContext(ctx))
}
