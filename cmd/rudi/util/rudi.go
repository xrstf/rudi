// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package util

import (
	"fmt"

	"go.xrstf.de/rudi"
)

func SetupRudiContext(files []any) (rudi.Context, error) {
	var (
		document rudi.Document
		err      error
	)

	if len(files) > 0 {
		document, err = rudi.NewDocument(files[0])
		if err != nil {
			return rudi.Context{}, fmt.Errorf("cannot use first input as document: %w", err)
		}
	} else {
		document, _ = rudi.NewDocument(nil)
	}

	vars := rudi.NewVariables().
		Set("files", files)

	ctx := rudi.NewContext(document, vars, rudi.NewBuiltInFunctions(), nil)

	return ctx, nil
}
