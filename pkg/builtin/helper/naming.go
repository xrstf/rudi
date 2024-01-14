// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package helper

import (
	"fmt"

	"go.xrstf.de/rudi/pkg/lang/ast"
)

func DecodeNamingVector(expr ast.Expression) (indexName string, valueName string, err error) {
	namingVector, ok := expr.(ast.VectorNode)
	if !ok {
		return "", "", fmt.Errorf("expected a vector, but got %T", expr)
	}

	size := len(namingVector.Expressions)
	if size < 1 || size > 2 {
		return "", "", fmt.Errorf("expected 1 or 2 identifiers in the naming vector, got %d", size)
	}

	if size == 1 {
		varNameIdent, ok := namingVector.Expressions[0].(ast.Identifier)
		if !ok {
			return "", "", fmt.Errorf("value variable name must be an identifier, got %T", namingVector.Expressions[0])
		}

		valueName = varNameIdent.Name
	} else {
		indexIdent, ok := namingVector.Expressions[0].(ast.Identifier)
		if !ok {
			return "", "", fmt.Errorf("index variable name must be an identifier, got %T", namingVector.Expressions[0])
		}

		varNameIdent, ok := namingVector.Expressions[1].(ast.Identifier)
		if !ok {
			return "", "", fmt.Errorf("value variable name must be an identifier, got %T", namingVector.Expressions[0])
		}

		indexName = indexIdent.Name
		valueName = varNameIdent.Name

		if indexName == valueName {
			return "", "", fmt.Errorf("cannot use %s for both value and index variable", indexName)
		}
	}

	return indexName, valueName, nil
}
