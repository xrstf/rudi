// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package pathexpr

// These types are deprecated, use jsonpath package instead.

type ObjectReader interface {
	GetObjectKey(name string) (any, error)
}

type ObjectWriter interface {
	ObjectReader
	SetObjectKey(name string, value any) (any, error)
}
