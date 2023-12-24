// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package encoding

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"go.xrstf.de/rudi/pkg/runtime/functions"
	"go.xrstf.de/rudi/pkg/runtime/types"
)

var (
	Functions = types.Functions{
		"to-base64":   functions.NewBuilder(toBase64Function).WithDescription("apply base64 encoding to the given string").Build(),
		"from-base64": functions.NewBuilder(fromBase64Function).WithDescription("decode a base64 encoded string").Build(),
		"to-json":     functions.NewBuilder(toJSONFunction).WithDescription("encode the given value using JSON").Build(),
		"from-json":   functions.NewBuilder(fromJSONFunction).WithDescription("decode a JSON string").Build(),
	}
)

func toBase64Function(value string) (any, error) {
	encoded := base64.StdEncoding.EncodeToString([]byte(value))

	return encoded, nil
}

func fromBase64Function(encoded string) (any, error) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("not valid base64: %w", err)
	}

	return string(decoded), nil
}

func toJSONFunction(value any) (any, error) {
	encoded, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	return string(encoded), nil
}

func fromJSONFunction(encoded string) (any, error) {
	var result any
	if err := json.Unmarshal([]byte(encoded), &result); err != nil {
		return nil, err
	}

	return result, nil
}
