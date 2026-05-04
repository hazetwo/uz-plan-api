package schedule

type Repository interface {
	GetFields() (map[string]string, bool)
	StoreFields(fields map[string]string)
	GetGroups(fieldID string) (map[string]string, bool)
	StoreGroups(fieldId string, groups map[string]string)
	GetSchedule(groupID string) ([]Entry, error)
	StoreSchedule(groupID string, entries []Entry)
}
