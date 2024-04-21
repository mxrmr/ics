// SPDX-License-Identifier: 0BSD

package ics

import "time"

type Event struct {
	Name         string
	UID          string
	StartDate    time.Time
	Duration     time.Duration
	TZ           *TZ
	LastModified time.Time
}
