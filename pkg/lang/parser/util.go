// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package parser

func toAnySlice(v any) []any {
	if v == nil {
		return nil
	}
	return v.([]any)
}
