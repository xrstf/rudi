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

	"go.xrstf.de/rudi/pkg/eval/types"
)

// (sha1 VAL:string)
func sha1Function(ctx types.Context, args []any) (any, error) {
	return hashFunc(ctx, args[0], sha1.New())
}

// (sha256 VAL:string)
func sha256Function(ctx types.Context, args []any) (any, error) {
	return hashFunc(ctx, args[0], sha256.New())
}

// (sha512 VAL:string)
func sha512Function(ctx types.Context, args []any) (any, error) {
	return hashFunc(ctx, args[0], sha512.New())
}

func hashFunc(ctx types.Context, arg any, h hash.Hash) (any, error) {
	str, err := ctx.Coalesce().ToString(arg)
	if err != nil {
		return nil, err
	}

	if _, err := io.WriteString(h, string(str)); err != nil {
		return nil, fmt.Errorf("error when hashing: %w", err)
	}

	encoded := hex.EncodeToString(h.Sum(nil))

	return encoded, nil
}
