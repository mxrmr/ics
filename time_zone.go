// SPDX-License-Identifier: 0BSD

package ics

import (
	"fmt"
	"time"
)

// A TZ defines a particular time zone.
type TZ struct {
	id     string
	ranges []tzRange
}

// A tzRange defines a particular standard/daylight interval of a time zone.
type tzRange struct {
	dst      string
	name     string
	start    time.Time
	rrule    string
	from, to string
}

// NewTZAmerica creates a time zone with America’s Daylight Savings Time rules.
func NewTZAmerica(id string, offset int, standard, daylight string) *TZ {
	standardOffset := fmt.Sprintf("%+05d", offset*100)
	daylightOffset := fmt.Sprintf("%+05d", (offset+1)*100)
	return &TZ{
		id: id,
		ranges: []tzRange{
			{
				name:  standard,
				start: tzDate(2007, 11, 4, 2, 0, daylight, offset+1),
				dst:   "STANDARD",
				rrule: "FREQ=YEARLY;BYMONTH=11;BYDAY=1SU",
				from:  daylightOffset,
				to:    standardOffset,
			},
			{
				name:  daylight,
				start: tzDate(2007, 3, 11, 2, 0, daylight, offset),
				dst:   "DAYLIGHT",
				rrule: "FREQ=YEARLY;BYMONTH=3;BYDAY=2SU",
				from:  standardOffset,
				to:    daylightOffset,
			},
		},
	}
}

func tzDate(year, month, day, hour, minute int, tzName string, tzOffset int) time.Time {
	return time.Date(
		year, time.Month(month), day, hour, minute, 0, 0,
		time.FixedZone(tzName, tzOffset*3600),
	)
}
