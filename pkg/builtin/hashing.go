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
