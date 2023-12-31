// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package hashing

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"io"

	"go.xrstf.de/rudi/pkg/runtime/functions"
	"go.xrstf.de/rudi/pkg/runtime/types"
)

var (
	Functions = types.Functions{
		"sha1":   functions.NewBuilder(sha1Function).WithDescription("return the lowercase hex representation of the SHA-1 hash").Build(),
		"sha256": functions.NewBuilder(sha256Function).WithDescription("return the lowercase hex representation of the SHA-256 hash").Build(),
		"sha512": functions.NewBuilder(sha512Function).WithDescription("return the lowercase hex representation of the SHA-512 hash").Build(),
	}
)

func sha1Function(value string) (any, error) {
	return hashFunc(value, sha1.New())
}

func sha256Function(value string) (any, error) {
	return hashFunc(value, sha256.New())
}

func sha512Function(value string) (any, error) {
	return hashFunc(value, sha512.New())
}

func hashFunc(value string, h hash.Hash) (any, error) {
	if _, err := io.WriteString(h, value); err != nil {
		return nil, fmt.Errorf("error when hashing: %w", err)
	}

	encoded := hex.EncodeToString(h.Sum(nil))

	return encoded, nil
}
