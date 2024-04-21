// SPDX-License-Identifier: 0BSD

package ics

type constProperty struct {
	key    string
	value  string
	params []param
}

func (p constProperty) writableKey() string {
	return p.key
}

func (p constProperty) writableValue() string {
	return p.value
}

func (p constProperty) writableParams() []param {
	return p.params
}
