package schedule

type Subgroup string

const (
	A Subgroup = "A"
	B Subgroup = "B"
)

func ParseSubgroup(s string) *Subgroup {
	switch Subgroup(s) {
	case A, B:
		return new(Subgroup(s))
	default:
		return nil
	}
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

type Entry struct {
	Subgroup  *Subgroup `json:"subgroup"`
	Start     *string   `json:"start"`
	End       *string   `json:"end"`
	Date      *string   `json:"date"`
	Subject   *string   `json:"subject"`
	ClassType *string   `json:"classType"`
	Teacher   *string   `json:"teacher"`
	Classroom *string   `json:"classroom"`
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

	group := ParseSubgroup(raw.Group)

	return Entry{
		Subgroup:  group,
		Start:     strPtr(raw.Start),
		End:       strPtr(raw.End),
		Date:      strPtr(raw.Date),
		Subject:   strPtr(raw.Subject),
		ClassType: strPtr(raw.ClassType),
		Teacher:   strPtr(raw.Teacher),
		Classroom: strPtr(raw.Classroom),
	}, nil
}
