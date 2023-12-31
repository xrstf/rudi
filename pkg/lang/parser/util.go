// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package parser

func toAnySlice(v any) []any {
	if v == nil {
		return nil
	}
	return v.([]any)
}
