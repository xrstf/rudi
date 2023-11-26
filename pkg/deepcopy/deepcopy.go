package deepcopy

import (
	"fmt"

	"go.xrstf.de/rudi/pkg/lang/ast"
)

func Clone[T any](val T) (T, error) {
	cloned, err := clone(val)
	if err != nil {
		empty := new(T)
		return *empty, err
	}

	if cloned == nil {
		empty := new(T)
		return *empty, nil
	}

	return cloned.(T), nil
}

func MustClone[T any](val T) T {
	cloned, err := Clone(val)
	if err != nil {
		panic(err)
	}

	return cloned
}

func clonePtr[T any](ptr *T) *T {
	if ptr == nil {
		return nil
	}

	cloned, _ := Clone(*ptr)
	return &cloned
}

//nolint:gocyclo
func clone(val any) (any, error) {
	switch asserted := val.(type) {
	// Go native types

	case nil:
		return asserted, nil
	case bool:
		return asserted, nil
	case int:
		return asserted, nil
	case int32:
		return asserted, nil
	case int64:
		return asserted, nil
	case float32:
		return asserted, nil
	case float64:
		return asserted, nil
	case string:
		return asserted, nil
	case map[string]any:
		return cloneObject(asserted)
	case []any:
		return cloneVector(asserted)

	// pointer to Go types

	case *bool:
		return clonePtr(asserted), nil
	case *int:
		return clonePtr(asserted), nil
	case *int32:
		return clonePtr(asserted), nil
	case *int64:
		return clonePtr(asserted), nil
	case *float32:
		return clonePtr(asserted), nil
	case *float64:
		return clonePtr(asserted), nil
	case *string:
		return clonePtr(asserted), nil
	case *map[string]any:
		return clonePtr(asserted), nil
	case *[]any:
		return clonePtr(asserted), nil

	// AST literals

	case ast.Null:
		return ast.Null{}, nil
	case ast.Bool:
		return asserted, nil
	case ast.Number:
		return ast.Number{Value: asserted.Value}, nil
	case ast.String:
		return asserted, nil
	case ast.Object:
		cloned, err := Clone(asserted.Data)
		if err != nil {
			return nil, err
		}
		return ast.Object{Data: cloned}, nil
	case ast.Vector:
		cloned, err := Clone(asserted.Data)
		if err != nil {
			return nil, err
		}
		return ast.Vector{Data: cloned}, nil

	// pointer to AST literals

	case *ast.Null:
		return clonePtr(asserted), nil
	case *ast.Bool:
		return clonePtr(asserted), nil
	case *ast.Number:
		return clonePtr(asserted), nil
	case *ast.String:
		return clonePtr(asserted), nil
	case *ast.Object:
		return clonePtr(asserted), nil
	case *ast.Vector:
		return clonePtr(asserted), nil

	default:
		return nil, fmt.Errorf("cannot deep-copy %T", val)
	}
}

func cloneVector(obj []any) ([]any, error) {
	result := make([]any, len(obj))
	for i, item := range obj {
		cloned, err := clone(item)
		if err != nil {
			return nil, err
		}

		result[i] = cloned
	}

	return result, nil
}

func cloneObject(obj map[string]any) (map[string]any, error) {
	result := map[string]any{}
	for key, value := range obj {
		cloned, err := clone(value)
		if err != nil {
			return nil, err
		}

		result[key] = cloned
	}

	return result, nil
}
