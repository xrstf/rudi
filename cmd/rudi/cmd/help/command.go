// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package help

import (
	_ "embed"
	"fmt"
	"strings"

	"go.xrstf.de/rudi/cmd/rudi/types"
	"go.xrstf.de/rudi/cmd/rudi/util"
	"go.xrstf.de/rudi/docs"

	"github.com/spf13/pflag"
)

//go:embed help.md
var helpText string

func Run(opts *types.Options, args []string) error {
	helpTopics := docs.Topics()

	// do not show function docs for "--help help if"
	if !opts.ShowHelp && len(args) == 2 && args[0] == "help" {
		rendered, err := util.RenderHelpTopic(helpTopics, args[1], 0)
		if err == nil {
			fmt.Println(rendered)
			return nil
		}

		fmt.Printf("Error: %v\n", err)
		fmt.Println()
		fmt.Println("The following topics are available:")
		fmt.Println()
		fmt.Println(util.RenderHelpTopics(helpTopics, 0))

		return nil
	}

	fmt.Println(util.RenderMarkdown(strings.TrimSpace(helpText), 0))
	fmt.Println(util.RenderHelpTopics(helpTopics, 0))
	fmt.Println()

	pflag.Usage()

	return nil
}
