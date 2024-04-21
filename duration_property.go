// SPDX-License-Identifier: 0BSD

package ics

import (
	"math"
	"strconv"
	"strings"
	"time"
)

type durationProperty struct {
	key    string
	value  time.Duration
	params []param
}

func (p durationProperty) writableKey() string {
	return p.key
}

func (p durationProperty) writableValue() string {
	return formatDuration(p.value)
}

func (p durationProperty) writableParams() []param {
	return p.params
}

func splitDuration(t time.Duration) (days, hours, minutes, seconds int) {
	total := int(math.Round(t.Seconds()))
	seconds = total % 60
	total /= 60
	minutes = total % 60
	total /= 60
	hours = total % 24
	total /= 24
	days = total
	return
}

func formatDuration(t time.Duration) string {
	days, hours, minutes, seconds := splitDuration(t)
	components := [3]struct {
		value  int
		suffix string
	}{
		{hours, "H"},
		{minutes, "M"},
		{seconds, "S"},
	}
	// Find the range of non-zero H/M/S components.
	start, end := len(components), len(components)
	for idx := 0; idx < len(components); idx++ {
		if components[idx].value > 0 {
			start = idx
			break
		}
	}
	for idx := len(components) - 1; idx >= 0; idx-- {
		if components[idx].value > 0 {
			end = idx + 1
			break
		}
	}
	var res strings.Builder
	res.Grow(16) // 16 = “P”, “0000”, “D”, “T”, 3x ( “00”, “_” )
	res.WriteString("P")
	if days > 0 {
		res.WriteString(strconv.Itoa(days))
		res.WriteString("D")
	} else if start == end {
		// If there aren't days, then we might need "0S" as a fallback.
		start--
	}
	if start < end {
		res.WriteString("T")
		for _, component := range components[start:end] {
			res.WriteString(strconv.Itoa(component.value))
			res.WriteString(component.suffix)
		}
	}
	return res.String()
}
