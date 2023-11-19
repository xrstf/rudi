// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package util

import (
	"fmt"

	"go.xrstf.de/otto"
)

func SetupOttoContext(files []any) (otto.Context, error) {
	var (
		document otto.Document
		err      error
	)

	if len(files) > 0 {
		document, err = otto.NewDocument(files[0])
		if err != nil {
			return otto.Context{}, fmt.Errorf("cannot use first input as document: %w", err)
		}
	} else {
		document, _ = otto.NewDocument(nil)
	}

	vars := otto.NewVariables().
		Set("files", files)

	ctx := otto.NewContext(document, otto.NewBuiltInFunctions(), vars)

	return ctx, nil
}
