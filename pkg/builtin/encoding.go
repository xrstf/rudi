// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"encoding/base64"
	"fmt"

	"go.xrstf.de/rudi/pkg/eval/types"
)

// (to-base64 VAL:string)
func toBase64Function(ctx types.Context, args []any) (any, error) {
	str, err := ctx.Coalesce().ToString(args[0])
	if err != nil {
		return nil, err
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(str))

	return encoded, nil
}

// (from-base64 VAL:string)
func fromBase64Function(ctx types.Context, args []any) (any, error) {
	str, err := ctx.Coalesce().ToString(args[0])
	if err != nil {
		return nil, err
	}

	decoded, err := base64.StdEncoding.DecodeString(string(str))
	if err != nil {
		return nil, fmt.Errorf("not valid base64: %w", err)
	}

	return string(decoded), nil
}
