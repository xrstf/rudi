// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package jsonpath

import "go.xrstf.de/rudi/pkg/lang/ast"

type Path []Step

func (p Path) IsValid() bool {
	for _, s := range p {
		switch s.(type) {
		case VectorStep,
			ObjectStep,
			DynamicVectorStep,
			DynamicObjectStep:
			continue
		default:
			return false
		}
	}

	return true
}

func (p Path) IsDynamic() bool {
	for _, s := range p {
		if isDynamicStep(s) {
			return true
		}
	}

	return false
}

func isDynamicStep(s Step) bool {
	switch s.(type) {
	case DynamicVectorStep, DynamicObjectStep:
		return true
	default:
		return false
	}
}

type Step any

type VectorStep interface {
	Index() int
}

type IndexStep int

func (i IndexStep) Index() int {
	return int(i)
}

type ObjectStep interface {
	Key() string
}

type KeyStep string

func (k KeyStep) Key() string {
	return string(k)
}

type DynamicVectorStep interface {
	Keep(index int, value any) (bool, error)
}

type DynamicObjectStep interface {
	Keep(key string, value any) (bool, error)
}

func FromEvaluatedPath(evaledPath ast.EvaluatedPathExpression) Path {
	p := Path{}
	for _, step := range evaledPath.Steps {
		if i := step.IntegerValue; i != nil {
			p = append(p, IndexStep(*i))
		} else if s := step.StringValue; s != nil {
			p = append(p, KeyStep(*s))
		}
	}

	return p
}

func toIntegerStep(s Step) (int, bool) {
	switch asserted := s.(type) {
	case int:
		return asserted, true
	case int32:
		return int(asserted), true
	case int64:
		return int(asserted), true
	default:
		return 0, false
	}
}

func toStringStep(s Step) (string, bool) {
	switch asserted := s.(type) {
	case string:
		return asserted, true
	default:
		return "", false
	}
}
