// SPDX-License-Identifier: 0BSD

package ics

import (
	"bytes"
)

func (c *Calendar) Bytes() []byte {
	var b bytes.Buffer
	c.writeTo(&b)
	return b.Bytes()
}

func (c *Calendar) writeTo(b *bytes.Buffer) {
	writeContainerTo(b, "VCALENDAR", func() {
		writePropertiesTo(b, []propertyWriter{
			constProperty{key: "VERSION", value: "2.0"},
			textProperty{key: "PRODID", value: c.ID},
		})
		for _, tz := range c.extractTZs() {
			tz.writeTo(b)
		}
		for _, event := range c.Events {
			event.writeTo(b)
		}
	})
}

func (tz TZ) writeTo(b *bytes.Buffer) {
	writeContainerTo(b, "VTIMEZONE", func() {
		writePropertiesTo(b, []propertyWriter{
			textProperty{key: "TZID", value: tz.id},
		})
		for _, range_ := range tz.ranges {
			range_.writeTo(b)
		}
	})
}

func (r tzRange) writeTo(b *bytes.Buffer) {
	writeContainerTo(b, r.dst, func() {
		writePropertiesTo(b, []propertyWriter{
			textProperty{key: "TZNAME", value: r.name},
			constProperty{key: "TZOFFSETFROM", value: r.from},
			constProperty{key: "TZOFFSETTO", value: r.to},
			dateProperty{key: "DTSTART", value: r.start},
			constProperty{key: "RRULE", value: r.rrule},
		})
	})
}

func (e Event) writeTo(b *bytes.Buffer) {
	writeContainerTo(b, "VEVENT", func() {
		writePropertiesTo(b, []propertyWriter{
			textProperty{key: "SUMMARY", value: e.Name},
			textProperty{key: "UID", value: e.UID},
			dateProperty{key: "DTSTAMP", value: e.LastModified, utc: true},
			dateProperty{key: "DTSTART", value: e.StartDate, utc: false, params: []param{
				{name: "TZID", value: e.TZ.id},
			}},
			durationProperty{key: "DURATION", value: e.Duration},
		})
	})
}

func writePropertyTo(b *bytes.Buffer, prop propertyWriter) {
	col := 0
	col = writeAndFold(b, col, prop.writableKey())
	for _, param_ := range prop.writableParams() {
		col = writeParamTo(b, col, param_)
	}
	col = writeAndFold(b, col, ":")
	col = writeAndFold(b, col, prop.writableValue())
	// The final newline isn’t subject to wrapping.
	// If it were to wrap, we’d end up with a blank line.
	_ = writeAndFold(b, 0, "\r\n")
}

func writeParamTo(b *bytes.Buffer, col int, p param) int {
	col = writeAndFold(b, col, ";")
	col = writeAndFold(b, col, p.name)
	col = writeAndFold(b, col, "=")
	col = writeAndFold(b, col, p.value)
	return col
}

const foldLimit = 75

// writeAndFold prints the specified string after wrapping to 75 bytes.
func writeAndFold(b *bytes.Buffer, col int, value string) int {
	for _, r := range value {
		for {
			n, _ := b.WriteRune(r)
			col += n
			if col <= foldLimit {
				break
			}
			b.Truncate(b.Len() - n)
			b.WriteString("\r\n ")
			col = 1
		}
	}
	return col
}

func writeContainerTo(b *bytes.Buffer, name string, f func()) {
	writePropertyTo(b, constProperty{key: "BEGIN", value: name})
	f()
	writePropertyTo(b, constProperty{key: "END", value: name})
}

func writePropertiesTo(b *bytes.Buffer, properties []propertyWriter) {
	for _, prop := range properties {
		writePropertyTo(b, prop)
	}
}
