// SPDX-License-Identifier: 0BSD

package ics

type propertyWriter interface {
	writableKey() string
	writableValue() string
	writableParams() []param
}
