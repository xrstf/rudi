// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"encoding/base64"
	"fmt"

	"go.xrstf.de/otto/pkg/lang/ast"
	"go.xrstf.de/otto/pkg/lang/eval"
	"go.xrstf.de/otto/pkg/lang/eval/types"
)

// (to-base64 VAL:string)
func toBase64Function(ctx types.Context, args []ast.Expression) (any, error) {
	if size := len(args); size != 1 {
		return nil, fmt.Errorf("expected 1 argument, got %d", size)
	}

	_, value, err := eval.EvalExpression(ctx, args[0])
	if err != nil {
		return nil, err
	}

	str, ok := value.(ast.String)
	if !ok {
		return nil, fmt.Errorf("argument is not string, but %T", value)
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(str))

	return ast.String(encoded), nil
}

// (from-base64 VAL:string)
func fromBase64Function(ctx types.Context, args []ast.Expression) (any, error) {
	if size := len(args); size != 1 {
		return nil, fmt.Errorf("expected 1 argument, got %d", size)
	}

	_, value, err := eval.EvalExpression(ctx, args[0])
	if err != nil {
		return nil, err
	}

	str, ok := value.(ast.String)
	if !ok {
		return nil, fmt.Errorf("argument is not string, but %T", value)
	}

	decoded, err := base64.StdEncoding.DecodeString(string(str))
	if err != nil {
		return nil, fmt.Errorf("argument is not valid base64: %w", err)
	}

	return ast.String(string(decoded)), nil
}