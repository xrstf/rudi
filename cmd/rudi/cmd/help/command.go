// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package help

import (
	"fmt"

	"go.xrstf.de/rudi/cmd/rudi/docs"
	"go.xrstf.de/rudi/cmd/rudi/options"
	"go.xrstf.de/rudi/cmd/rudi/util"

	"github.com/spf13/pflag"
)

func Run(opts *options.Options, args []string) error {
	// do not show function docs for "--help help if"
	if !opts.ShowHelp && len(args) == 2 && args[0] == "help" {
		rendered, err := util.RenderHelpTopic(args[1], 0)
		if err == nil {
			fmt.Print(rendered)
			return nil
		}

		fmt.Printf("Error: %v\n", err)
		fmt.Println()
	}

	content, err := docs.RenderFile("cmd-help.md", nil)
	if err != nil {
		return err
	}

	fmt.Print(content)
	pflag.Usage()

	return nil
}
