package model

type Predicate func(Entry) bool

func FilterEntries(entries []Entry, predicates ...Predicate) []Entry {
	var result []Entry
	for _, e := range entries {
		match := true
		for _, p := range predicates {
			if !p(e) {
				match = false
				break
			}
		}
		if match {
			result = append(result, e)
		}
	}
	return result
}

func DayPredicate(day *string) Predicate {
	return func(e Entry) bool {
		return matchesDay(e, day)
	}
}

func WeekPredicate(week *string) Predicate {
	return func(e Entry) bool {
		return matchesWeek(e, week)
	}
}

func SubgroupPredicate(subgroup *Subgroup) Predicate {
	return func(e Entry) bool {
		return matchesSubgroup(e, subgroup)
	}
}
