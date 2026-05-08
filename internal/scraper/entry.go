package scraper

import "uz-plan-api/internal/model"

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

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func FromScraper(raw RawEntry) model.Entry {

	group := model.ParseSubgroup(raw.Group)

	return model.Entry{
		Subgroup:  group,
		Start:     strPtr(raw.Start),
		End:       strPtr(raw.End),
		Date:      strPtr(raw.Date),
		Subject:   strPtr(raw.Subject),
		ClassType: strPtr(raw.ClassType),
		Teacher:   strPtr(raw.Teacher),
		Classroom: strPtr(raw.Classroom),
	}

}
