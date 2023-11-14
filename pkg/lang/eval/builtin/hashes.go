// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"io"

	"go.xrstf.de/otto/pkg/lang/ast"
	"go.xrstf.de/otto/pkg/lang/eval"
	"go.xrstf.de/otto/pkg/lang/eval/types"
)

func hashFunc(ctx types.Context, args []ast.Expression, h hash.Hash) (any, error) {
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

	if _, err := io.WriteString(h, string(str)); err != nil {
		return nil, fmt.Errorf("error when hashing: %w", err)
	}

	encoded := hex.EncodeToString(h.Sum(nil))

	return ast.String(encoded), nil
}

// (sha1 VAL:string)
func sha1Function(ctx types.Context, args []ast.Expression) (any, error) {
	return hashFunc(ctx, args, sha1.New())
}

// (sha256 VAL:string)
func sha256Function(ctx types.Context, args []ast.Expression) (any, error) {
	return hashFunc(ctx, args, sha256.New())
}

// (sha512 VAL:string)
func sha512Function(ctx types.Context, args []ast.Expression) (any, error) {
	return hashFunc(ctx, args, sha512.New())
}
