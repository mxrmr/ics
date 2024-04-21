// SPDX-License-Identifier: 0BSD

package ics

import "testing"

func TestTZId(t *testing.T) {
	testCases := []string{
		"America/New_York",
		"America/Los_Angeles",
	}
	for _, tc := range testCases {
		tz := NewTZAmerica(tc, 0, "", "")
		if tz.id != tc {
			t.Errorf("tz.id = %v, want %v", tz.id, tc)
		}
	}
}

func TestTZOffset(t *testing.T) {
	testCases := []struct {
		offset   int
		from, to string
	}{
		{-5, "-0500", "-0400"},
		{-6, "-0600", "-0500"},
		{-7, "-0700", "-0600"},
		{-8, "-0800", "-0700"},
	}
	for _, tc := range testCases {
		tz := NewTZAmerica("", tc.offset, "", "")
		for _, range_ := range tz.ranges {
			// Determine the expected transition direction.
			from, to := tc.from, tc.to
			if range_.dst == "STANDARD" {
				from, to = to, from
			}
			if range_.from != from {
				t.Errorf("{%d, %v}.from = %v, want %v", tc.offset, range_.dst, range_.from, from)
			}
			if range_.to != to {
				t.Errorf("{%d, %v}.to = %v, want %v", tc.offset, range_.dst, range_.to, to)
			}
		}
	}
}
