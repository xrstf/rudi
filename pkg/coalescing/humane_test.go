// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package coalescing

import (
	"testing"
)

func TestHumaneCoalescer(t *testing.T) {
	testCoalescer(t, NewHumane(), getHumaneTestcases())
}

func getHumaneTestcases() []testcase {
	return []testcase{
		//         (source, canBeNull, toBool, toInt, toFloat, toNumber, toString, toVector, toObject)
		// nil source value
		newTestcase(nil, true, false, int64(0), 0.0, newNum(int64(0)), "", []any{}, map[string]any{}),
		// boolean source values
		newTestcase(true, invalid, true, int64(1), 1.0, newNum(int64(1)), "true", invalid, invalid),
		newTestcase(false, true, false, int64(0), 0.0, newNum(int64(0)), "false", invalid, invalid),
		// numeric source values
		newTestcase(0, true, false, int64(0), 0.0, newNum(int64(0)), "0", invalid, invalid),
		newTestcase(0.0, true, false, int64(0), 0.0, newNum(int64(0)), "0", invalid, invalid),
		newTestcase(0.1, invalid, true, invalid, 0.1, newNum(0.1), "0.1", invalid, invalid),
		newTestcase(1, invalid, true, int64(1), 1.0, newNum(int64(1)), "1", invalid, invalid),
		newTestcase(1.0, invalid, true, int64(1), 1.0, newNum(int64(1)), "1", invalid, invalid),
		newTestcase(-3.14, invalid, true, invalid, -3.14, newNum(-3.14), "-3.14", invalid, invalid),
		// string source values
		newTestcase("", true, false, int64(0), 0.0, newNum(int64(0)), "", invalid, invalid),
		newTestcase(" ", invalid, true, int64(0), 0.0, newNum(int64(0)), " ", invalid, invalid),
		newTestcase("\n", invalid, true, int64(0), 0.0, newNum(int64(0)), "\n", invalid, invalid),
		newTestcase("0", invalid, false, int64(0), 0.0, newNum(int64(0)), "0", invalid, invalid),
		newTestcase("000", invalid, true, int64(0), 0.0, newNum(int64(0)), "000", invalid, invalid),
		newTestcase(" 0 ", invalid, true, int64(0), 0.0, newNum(int64(0)), " 0 ", invalid, invalid),
		newTestcase(" 000 ", invalid, true, int64(0), 0.0, newNum(int64(0)), " 000 ", invalid, invalid),
		newTestcase("1", invalid, true, int64(1), 1.0, newNum(int64(1)), "1", invalid, invalid),
		newTestcase("001", invalid, true, int64(1), 1.0, newNum(int64(1)), "001", invalid, invalid),
		newTestcase("1.0", invalid, true, int64(1), 1.0, newNum(int64(1)), "1.0", invalid, invalid),
		newTestcase(" 1.1 ", invalid, true, invalid, 1.1, newNum(1.1), " 1.1 ", invalid, invalid),
		// vector source values
		newTestcase([]any{}, true, false, invalid, invalid, invalid, invalid, []any{}, map[string]any{}),
		newTestcase([]any{""}, invalid, true, invalid, invalid, invalid, invalid, []any{""}, invalid),
		newTestcase([]any{"foo"}, invalid, true, invalid, invalid, invalid, invalid, []any{"foo"}, invalid),
		// object source values
		newTestcase(map[string]any{}, true, false, invalid, invalid, invalid, invalid, []any{}, map[string]any{}),
		newTestcase(map[string]any{"": ""}, invalid, true, invalid, invalid, invalid, invalid, invalid, map[string]any{"": ""}),
	}
}
