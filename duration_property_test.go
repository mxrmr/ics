// SPDX-License-Identifier: 0BSD

package ics

import (
	"testing"
	"time"
)

func TestDurationProperty(t *testing.T) {
	testCases := []struct {
		input string
		want  string
	}{
		{"0s", "PT0S"},
		{"30s", "PT30S"},
		{"1m30s", "PT1M30S"},
		{"11m27s", "PT11M27S"},
		{"11m", "PT11M"},
		{"7h", "PT7H"},
		{"8h34m", "PT8H34M"},
		{"9h34m17s", "PT9H34M17S"},
		{"10h19s", "PT10H0M19S"},
		{"24h", "P1D"},
		{"144h", "P6D"},
		{"145h", "P6DT1H"},
		{"48h1s", "P2DT1S"},
		{"72h10m", "P3DT10M"},
		{"73h5m17s", "P3DT1H5M17S"},
	}
	for _, tc := range testCases {
		input, err := time.ParseDuration(tc.input)
		if err != nil {
			t.Errorf("ParseDuration(%q) error: %v", tc.input, err)
			continue
		}
		got := durationProperty{value: input}.writableValue()
		if got != tc.want {
			t.Errorf("T{value: %q}.writableValue() = %v, want %v", tc.input, got, tc.want)
		}
	}
}
