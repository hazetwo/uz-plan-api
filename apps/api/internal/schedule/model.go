package schedule

import (
	"fmt"
	"time"
)

type SubGroup string

const (
	A SubGroup = "A"
	B SubGroup = "B"
)

func ParseSubGroup(s string) (*SubGroup, error) {
	if s == "" {
		return nil, nil
	}
	switch SubGroup(s) {
	case A, B:
		return new(SubGroup(s)), nil
	default:
		return nil, fmt.Errorf("unknown subgroup: %s", s)
	}
}

type Date struct {
	start *time.Time
	end   *time.Time
	date  time.Time
}

func ParseDate(rawStart string, rawEnd string, rawDate string) (Date, error) {
	date, err := time.Parse("2006-01-02", rawDate)
	if err != nil {
		return Date{}, fmt.Errorf("invalid date: %w", err)
	}

	d := Date{date: date}

	if rawStart != "" {
		start, err := time.Parse("15:04", rawStart)
		if err != nil {
			return Date{}, fmt.Errorf("invalid start time: %w", err)
		}
		d.start = &start
	}

	if rawEnd != "" {
		end, err := time.Parse("15:04", rawEnd)
		if err != nil {
			return Date{}, fmt.Errorf("invalid end time: %w", err)
		}
		d.end = &end
	}

	return d, nil
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func (d Date) String() string {
	s := d.date.Format("2006-01-02")
	if d.start != nil && d.end != nil {
		s += fmt.Sprintf(" %s-%s", d.start.Format("15:04"), d.end.Format("15:04"))
	}
	return s
}

type Entry struct {
	subGroup  *SubGroup
	date      *Date
	subject   *string
	classType *string
	teacher   *string
	classroom *string
}

func (e Entry) SubGroup() *SubGroup {
	return e.subGroup
}

func (e Entry) Date() *Date {
	return e.date
}

func (e Entry) Subject() *string {
	return e.subject
}

func (e Entry) ClassType() *string {
	return e.classType
}

func (e Entry) Teacher() *string {
	return e.teacher
}

func (e Entry) Classroom() *string {
	return e.classroom
}

type RawEntry struct {
	Group     string
	Start     string
	End       string
	Date      string
	Subject   string
	ClassType string
	Teacher   string
	Classroom string
}

func FromScraper(raw RawEntry) (Entry, error) {
	hasStart := raw.Start != ""
	hasEnd := raw.End != ""

	if hasStart != hasEnd {
		return Entry{}, fmt.Errorf("incomplete time fields: start=%q end=%q", raw.Start, raw.End)
	}

	if (hasStart || hasEnd) && raw.Date == "" {
		return Entry{}, fmt.Errorf("time fields require a date")
	}

	var date *Date
	if raw.Date != "" {
		d, err := ParseDate(raw.Start, raw.End, raw.Date)
		if err != nil {
			return Entry{}, err
		}
		date = &d
	}

	group, err := ParseSubGroup(raw.Group)
	if err != nil {
		return Entry{}, err
	}

	return Entry{
		subGroup:  group,
		date:      date,
		subject:   strPtr(raw.Subject),
		classType: strPtr(raw.ClassType),
		teacher:   strPtr(raw.Teacher),
		classroom: strPtr(raw.Classroom),
	}, nil
}
