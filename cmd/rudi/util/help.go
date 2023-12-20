// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package util

import (
	"fmt"
	"strings"

	"go.xrstf.de/rudi/cmd/rudi/batteries"
	"go.xrstf.de/rudi/cmd/rudi/docs"
	rudidocs "go.xrstf.de/rudi/pkg/docs"
)

func RenderHelpTopic(selectedTopic string, indent int) (string, error) {
	// Did the user request a regular documentation page?
	topics := docs.Topics()

	for _, topic := range topics {
		for _, cliName := range topic.CliNames {
			if strings.EqualFold(cliName, selectedTopic) {
				return topic.Render(nil)
			}
		}
	}

	// Check if a function of that name is available with documentation.
	modules := []rudidocs.Module{}
	modules = append(modules, batteries.SafeBuiltInModules...)
	modules = append(modules, batteries.UnsafeBuiltInModules...)
	modules = append(modules, batteries.ExtendedModules...)

	for _, mod := range modules {
		for funcName := range mod.Functions {
			if strings.EqualFold(funcName, selectedTopic) {
				return docs.RenderFunction(funcName, nil)
			}
		}
	}

	return "", fmt.Errorf("no help available for %q", selectedTopic)
}
