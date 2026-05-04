package schedule

import (
	"time"
)

type Filter struct {
	Day      *string   // ISO date to derive day, nil = no filter
	Week     *string   // ISO date to derive week, nil = no filter
	Subgroup *Subgroup // nil = no filter
}

func matchesDay(e Entry, day *string) bool {
	if day == nil {
		return true
	}

	if e.Date == nil {
		return false
	}

	return *e.Date == *day
}

func matchesWeek(e Entry, week *string) bool {
	if week == nil {
		return true
	}

	if e.Date == nil {
		return false
	}

	td, err := time.Parse("2006-01-02", *week)
	if err != nil {
		return false
	}
	ty, tw := td.ISOWeek()

	d, err := time.Parse("2006-01-02", *e.Date)
	if err != nil {
		return false
	}

	y, w := d.ISOWeek()

	return y == ty && w == tw
}

func matchesSubgroup(e Entry, group *Subgroup) bool {
	if group == nil {
		return true
	}

	if e.Subgroup == nil {
		return true
	}

	return *group == *e.Subgroup
}
