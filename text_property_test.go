// SPDX-License-Identifier: 0BSD

package ics

import "testing"

func TestTextProperty(t *testing.T) {
	testCases := []struct {
		input, want string
	}{
		{"", ""},
		{"0", "0"},
		{"a", "a"},
		{"é", "é"},
		{"\x00", "\ufffd"},
		{"\x01", "\ufffd"},
		{"\x02", "\ufffd"},
		{"\x7f", "\ufffd"},
		{"\t", "\t"},
		{" ", " "},
		{"\n", "\\n"},
		{";", "\\;"},
		{",", "\\,"},
		{":", ":"},
		{"\\n", "\\\\n"},
		{"\\\n", "\\\\\\n"},
		{"\"", "\""},
	}
	l := len(testCases)
	for start := 0; start < l; start++ {
		var input, want string
		for index := 0; index < l; index++ {
			tc := testCases[(start+index)%l]
			input += tc.input
			want += tc.want
		}
		testCases = append(testCases, struct{ input, want string }{input, want})
	}
	for _, tc := range testCases {
		got := textProperty{value: tc.input}.writableValue()
		if got != tc.want {
			t.Errorf("T{v: %q}.writableValue() = %q, want %q", tc.input, got, tc.want)
		}
	}
}
