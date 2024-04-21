// SPDX-License-Identifier: 0BSD

package ics

import (
	"bytes"
	"strconv"
	"testing"
	"time"
)

// TestCalendar tests a full calendar encoding with multiple Events.
func TestCalendar(t *testing.T) {
	cal := &Calendar{}
	var lines []string
	for _, snippet := range tzSnippets {
		lines = append(lines, snippet.lines...)
	}
	for _, snippet := range evSnippets {
		cal.Events = append(cal.Events, snippet.event)
		lines = append(lines, snippet.lines...)
	}
	testWriter(t, cal.writeTo, flattenLines(calStart, lines, calEnd))
}

// TestCalendarEmpty tests an empty calendar.
func TestCalendarEmpty(t *testing.T) {
	cal := &Calendar{}
	l := flattenLines(calStart, calEnd)
	testWriter(t, cal.writeTo, l)
}

// TestCalendarProdId tests the PRODID property.
func TestCalendarProdId(t *testing.T) {
	testCases := []string{
		"ProdID With Spaces",
		"A Different ProdID",
		"-//Blah/Bleh v1.0//EN",
	}
	for _, tc := range testCases {
		cal := &Calendar{ID: tc}
		lines := flattenLines(
			calStart[:len(calStart)-1],
			[]string{"PRODID:" + tc},
			calEnd,
		)
		testWriter(t, cal.writeTo, lines)
	}
}

// TestTZ tests time zone encodings.
func TestTZ(t *testing.T) {
	for _, snippet := range tzSnippets {
		testWriter(t, snippet.tz.writeTo, snippet.lines)
	}
}

// TestEvent tests event encodings.
func TestEvent(t *testing.T) {
	for _, snippet := range evSnippets {
		testWriter(t, snippet.event.writeTo, snippet.lines)
	}
}

// SAMPLE DATA

var calStart = []string{
	"BEGIN:VCALENDAR",
	"VERSION:2.0",
	"PRODID:",
}

var calEnd = []string{
	"END:VCALENDAR",
}

var tzSnippets = []struct {
	tz    *TZ
	lines []string
}{
	{NewTZAmerica("America/Los_Angeles", -8, "PST", "PDT"), []string{
		"BEGIN:VTIMEZONE",
		"TZID:America/Los_Angeles",
		"BEGIN:STANDARD",
		"TZNAME:PST",
		"TZOFFSETFROM:-0700",
		"TZOFFSETTO:-0800",
		"DTSTART:20071104T020000",
		"RRULE:FREQ=YEARLY;BYMONTH=11;BYDAY=1SU",
		"END:STANDARD",
		"BEGIN:DAYLIGHT",
		"TZNAME:PDT",
		"TZOFFSETFROM:-0800",
		"TZOFFSETTO:-0700",
		"DTSTART:20070311T020000",
		"RRULE:FREQ=YEARLY;BYMONTH=3;BYDAY=2SU",
		"END:DAYLIGHT",
		"END:VTIMEZONE",
	}},
	{NewTZAmerica("America/New_York", -5, "EST", "EDT"), []string{
		"BEGIN:VTIMEZONE",
		"TZID:America/New_York",
		"BEGIN:STANDARD",
		"TZNAME:EST",
		"TZOFFSETFROM:-0400",
		"TZOFFSETTO:-0500",
		"DTSTART:20071104T020000",
		"RRULE:FREQ=YEARLY;BYMONTH=11;BYDAY=1SU",
		"END:STANDARD",
		"BEGIN:DAYLIGHT",
		"TZNAME:EDT",
		"TZOFFSETFROM:-0500",
		"TZOFFSETTO:-0400",
		"DTSTART:20070311T020000",
		"RRULE:FREQ=YEARLY;BYMONTH=3;BYDAY=2SU",
		"END:DAYLIGHT",
		"END:VTIMEZONE",
	}},
}

var evSnippets = []struct {
	event Event
	lines []string
}{
	{makeEvent("Another", "PST", 201902210457, 201902011234, 4*time.Minute), []string{
		"BEGIN:VEVENT",
		"SUMMARY:Another",
		"UID:190221@test-domain",
		"DTSTAMP:20190201T203400Z",
		"DTSTART;TZID=America/Los_Angeles:20190221T045700",
		"DURATION:PT4M",
		"END:VEVENT",
	}},
	{makeEvent("Event Name", "EDT", 202011072329, 201906022210, 90*time.Minute), []string{
		"BEGIN:VEVENT",
		"SUMMARY:Event Name",
		"UID:201107@test-domain",
		"DTSTAMP:20190603T021000Z",
		"DTSTART;TZID=America/New_York:20201107T232900",
		"DURATION:PT1H30M",
		"END:VEVENT",
	}},
}

func testWriter(t *testing.T, f func(*bytes.Buffer), lines []string) {
	var buffer bytes.Buffer
	f(&buffer)
	got := buffer.Bytes()

	buffer = bytes.Buffer{}
	for _, line := range lines {
		buffer.WriteString(line)
		buffer.WriteString("\r\n")
	}
	want := buffer.Bytes()

	if !bytes.Equal(got, want) {
		t.Errorf("\ngot:\n%v\nwant:\n%v", string(got), string(want))
	}
}

func flattenLines(liness ...[]string) []string {
	var res []string
	for _, lines := range liness {
		res = append(res, lines...)
	}
	return res
}

func parseTZOffset(offset string) *time.Location {
	parsedOffset, _ := strconv.Atoi(offset)
	return time.FixedZone(offset, parsedOffset/100*3600)
}

// findTZNamed searches tzSnippets for a range with a particular name.
func findTZNamed(tzName string) (tz *TZ, loc *time.Location) {
	for _, snippet := range tzSnippets {
		for _, r := range snippet.tz.ranges {
			if r.name == tzName {
				tz, loc = snippet.tz, parseTZOffset(r.to)
			}
		}
	}
	return
}

// parseTime converts an integer of the form yyyymmddHHMM to a time.Time.
func parseTime(value int, loc *time.Location) time.Time {
	minute := value % 100
	value /= 100
	hour := value % 100
	value /= 100
	day := value % 100
	value /= 100
	month := value % 100
	value /= 100
	year := value
	return time.Date(year, time.Month(month), day, hour, minute, 0, 0, loc)
}

// makeEvent makes an event with a TZ from tzSnippets
func makeEvent(name, tzName string, start, modified int, duration time.Duration) Event {
	tz, loc := findTZNamed(tzName)
	startDate := parseTime(start, loc)
	lastModified := parseTime(modified, loc).UTC()
	return Event{
		Name:         name,
		UID:          startDate.Format("060102") + "@test-domain",
		StartDate:    startDate,
		Duration:     duration,
		TZ:           tz,
		LastModified: lastModified,
	}
}
