// SPDX-License-Identifier: 0BSD

package ics

import (
	"testing"
	"time"
)

func TestDateProperty(t *testing.T) {
	date := func(year, month, day, hour, minute, second int) time.Time {
		return time.Date(year, time.Month(month), day, hour, minute, second, 0, time.UTC)
	}
	testCases := []struct {
		input time.Time
		utc   bool
		want  string
	}{
		{date(2022, 12, 31, 0, 34, 43), false, "20221231T003443"},
		{date(2021, 1, 1, 23, 12, 23), false, "20210101T231223"},
		{date(2020, 3, 12, 5, 0, 59), false, "20200312T050059"},
		{date(2019, 9, 20, 12, 59, 7), false, "20190920T125907"},
		{date(2018, 11, 9, 17, 24, 0), false, "20181109T172400"},
		{date(2019, 9, 20, 12, 59, 7), true, "20190920T125907Z"},
		{date(2018, 11, 9, 17, 24, 0), true, "20181109T172400Z"},
	}
	for _, tc := range testCases {
		got := dateProperty{value: tc.input, utc: tc.utc}.writableValue()
		if got != tc.want {
			t.Errorf("T{value: %q, utc: %v}.writableValue() = %v, want %v", tc.input, tc.utc, got, tc.want)
		}
	}
}
