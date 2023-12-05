// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package builtin

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
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
