// SPDX-License-Identifier: 0BSD

package ics

import (
	"slices"
	"strings"
	"testing"
)

func tzSliceString(tzs []*TZ) string {
	var tzIDs []string
	for _, tz := range tzs {
		tzIDs = append(tzIDs, tz.id)
	}
	return strings.Join(tzIDs, ", ")
}

func TestExtractTZs(t *testing.T) {
	pt := NewTZAmerica("PT", -8, "PST", "PDT")
	ct := NewTZAmerica("CT", -6, "CST", "CDT")
	et := NewTZAmerica("ET", -5, "EST", "EDT")

	testCases := []struct {
		input []*TZ
		want  []*TZ
	}{
		// test empty
		{[]*TZ{}, []*TZ{}},
		// test removing duplicates
		{[]*TZ{pt}, []*TZ{pt}},
		{[]*TZ{pt, pt}, []*TZ{pt}},
		{[]*TZ{pt, pt, pt}, []*TZ{pt}},
		{[]*TZ{et}, []*TZ{et}},
		{[]*TZ{et, et}, []*TZ{et}},
		{[]*TZ{et, et, et}, []*TZ{et}},
		// test multiple results are sorted/de-duped
		{[]*TZ{et, ct}, []*TZ{ct, et}},
		{[]*TZ{ct, et}, []*TZ{ct, et}},
		{[]*TZ{et, et, ct}, []*TZ{ct, et}},
		{[]*TZ{et, ct, et}, []*TZ{ct, et}},
		{[]*TZ{ct, et, et}, []*TZ{ct, et}},
		{[]*TZ{ct, et, pt}, []*TZ{ct, et, pt}},
		{[]*TZ{et, pt, ct}, []*TZ{ct, et, pt}},
		{[]*TZ{et, ct, pt, ct}, []*TZ{ct, et, pt}},
		{[]*TZ{et, ct, pt, et, ct}, []*TZ{ct, et, pt}},
		{[]*TZ{et, ct, pt, pt, et, ct}, []*TZ{ct, et, pt}},
	}

	for _, tc := range testCases {
		cal := &Calendar{}
		for _, tz := range tc.input {
			cal.Events = append(cal.Events, Event{TZ: tz})
		}
		got := cal.extractTZs()
		if !slices.Equal(got, tc.want) {
			t.Errorf(
				"[%v].extractTZs() = %v, want %v",
				tzSliceString(tc.input), tzSliceString(got), tzSliceString(tc.want),
			)
		}
	}
}
