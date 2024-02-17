// SPDX-FileCopyrightText: 2024 Christoph Mewes
// SPDX-License-Identifier: MIT

package jsonpath

import "strconv"

type Path []Step

func (p Path) IsValid() bool {
	for _, s := range p {
		switch s.(type) {
		case SingleStep, FilterStep:
			continue
		default:
			return false
		}
	}

	return true
}

func (p Path) HasFilterSteps() bool {
	for _, s := range p {
		if isFilterStep(s) {
			return true
		}
	}

	return false
}

func isFilterStep(s Step) bool {
	_, ok := s.(FilterStep)
	return ok
}

type Step any

type SingleStep interface {
	ToIndex() (int, bool)
	ToKey() (string, bool)
}

type IndexStep int

func (i IndexStep) ToIndex() (int, bool) {
	return int(i), true
}

func (i IndexStep) ToKey() (string, bool) {
	return "", false
}

type KeyStep string

func (k KeyStep) ToIndex() (int, bool) {
	return 0, false
}

func (k KeyStep) ToKey() (string, bool) {
	return string(k), true
}

type flexStep string

func (f flexStep) ToIndex() (int, bool) {
	index, err := strconv.ParseInt(string(f), 10, 64)
	if err != nil {
		return 0, false
	}

	return int(index), false
}

func (f flexStep) ToKey() (string, bool) {
	return string(f), true
}

type FilterStep interface {
	Keep(key any, value any) (bool, error)
}

func indexOrKey(s SingleStep) (*int, *string) {
	index, ok := s.ToIndex()
	if ok {
		return &index, nil
	}

	key, ok := s.ToKey()
	if ok {
		return nil, &key
	}

	return nil, nil
}
