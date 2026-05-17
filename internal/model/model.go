package model

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

type Field struct {
	Faculty string `json:"faculty"`
	Name    string `json:"name"`
}

type Fields map[string]Field
type Groups map[string]string
