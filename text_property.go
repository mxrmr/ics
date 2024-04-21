// SPDX-License-Identifier: 0BSD

package ics

import (
	"strings"
	"unicode/utf8"
)

type textProperty struct {
	key    string
	value  string
	params []param
}

func (p textProperty) writableKey() string {
	return p.key
}

func (p textProperty) writableValue() string {
	return textPropertyReplacer.Replace(p.value)
}

func (p textProperty) writableParams() []param {
	return p.params
}

func escapedReplacements() []string {
	return []string{
		"\\", "\\\\",
		"\n", "\\n",
		",", "\\,",
		";", "\\;",
	}
}

func invalidReplacements() []string {
	allowedAscii := []struct {
		start, end int
	}{
		// TSAFE-CHAR
		// WSP
		{' ', ' '},
		{'\t', '\t'},

		// Ranges
		{0x21, 0x21}, // 0x22 = "
		{0x23, 0x2B}, // 0x2C = ,
		{0x2D, 0x39}, // 0x3A = :  0x3B = ;
		{0x3C, 0x5B}, // 0x5C = \
		{0x5D, 0x7E}, // 0x7F = <del>

		// ":"
		{':', ':'},

		// DQUOTE
		{'"', '"'},
	}
	// Mark the allowed characters in the bit array.
	var include [0x80]bool
	for _, chars := range allowedAscii {
		for i := chars.start; i <= chars.end; i++ {
			include[i] = true
		}
	}
	// Return replacements for all invalid characters.
	// Characters that can be escaped are handled first by `escapedReplacements`.
	var res []string
	for char, inc := range include {
		if !inc {
			res = append(res, string(rune(char)))
			res = append(res, string(rune(utf8.RuneError)))
		}
	}
	return res
}

var textPropertyReplacer = strings.NewReplacer(
	append(escapedReplacements(), invalidReplacements()...)...,
)
