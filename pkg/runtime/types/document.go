// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package types

type Document struct {
	data any
}

func NewDocument(data any) (Document, error) {
	return Document{
		data: data,
	}, nil
}

func (d *Document) Data() any {
	return d.data
}

func (d *Document) Set(wrappedData any) {
	d.data = wrappedData
}
