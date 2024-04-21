// SPDX-License-Identifier: 0BSD

package ics

import (
	"slices"
	"sort"
)

type Calendar struct {
	ID     string
	Events []Event
}

// extractTZs builds a sorted slice of unique *TZs
func (c *Calendar) extractTZs() []*TZ {
	var res []*TZ
	for _, event := range c.Events {
		insertionIdx := sort.Search(len(res), func(idx int) bool {
			return res[idx].id >= event.TZ.id
		})
		if insertionIdx < len(res) && res[insertionIdx].id == event.TZ.id {
			continue
		}
		res = slices.Insert(res, insertionIdx, event.TZ)
	}
	return res
}
