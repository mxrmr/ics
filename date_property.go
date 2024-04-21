// SPDX-License-Identifier: 0BSD

package ics

import "time"

type dateProperty struct {
	key    string
	value  time.Time
	utc    bool
	params []param
}

func (p dateProperty) writableKey() string {
	return p.key
}

func (p dateProperty) writableValue() string {
	fmt := "20060102T150405"
	if p.utc {
		fmt += "Z"
	}
	return p.value.Format(fmt)
}

func (p dateProperty) writableParams() []param {
	return p.params
}
