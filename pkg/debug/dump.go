// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package debug

import (
	"fmt"
	"io"

	"go.xrstf.de/rudi/pkg/eval/types"
	"go.xrstf.de/rudi/pkg/lang/ast"
)

var Indent = "  "

const DoNotIndent = -1

func Dump(p any, out io.Writer) error {
	return dumpAny(p, out, 0)
}

func DumpSingleline(p any, out io.Writer) error {
	return dumpAny(p, out, DoNotIndent)
}

func dumpAny(val any, out io.Writer, depth int) error {
	switch asserted := val.(type) {
	case nil:
		return DumpNull(out)
	case ast.Null:
		return DumpNull(out)
	case bool:
		return DumpBool(asserted, out)
	case ast.Bool:
		return DumpBool(bool(asserted), out)
	case int:
		return DumpNumber(asserted, out)
	case int32:
		return DumpNumber(asserted, out)
	case int64:
		return DumpNumber(asserted, out)
	case float32:
		return DumpNumber(asserted, out)
	case float64:
		return DumpNumber(asserted, out)
	case ast.Number:
		return DumpNumber(asserted.Value, out)
	case string:
		return DumpString(asserted, out)
	case ast.String:
		return DumpString(string(asserted), out)
	case []any:
		return DumpVector(asserted, out, depth)
	case ast.VectorNode:
		return DumpVectorNode(asserted, out, depth)
	case map[string]any:
		return DumpObject(asserted, out, depth)
	case ast.ObjectNode:
		return DumpObjectNode(asserted, out, depth)
	case ast.Symbol:
		return DumpSymbol(&asserted, out, depth)
	case ast.Tuple:
		return DumpTuple(&asserted, out, depth)
	case ast.Identifier:
		return DumpIdentifier(&asserted, out)
	case ast.Statement:
		return DumpStatement(&asserted, out, depth)
	case ast.Program:
		return DumpProgram(&asserted, out, depth)
	}

	wrapped, err := types.WrapNative(val)
	if err != nil {
		return fmt.Errorf("cannot dump values of type %T", val)
	}

	// as long as dumpAny() can handle all possible types returned by WrapNative, this won't be an infinite loop
	return dumpAny(wrapped, out, depth)
}

func writeString(out io.Writer, str string) error {
	_, err := out.Write([]byte(str))
	return err
}
