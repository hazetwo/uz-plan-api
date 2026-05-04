package schedule

type Repository interface {
	GetFields() (map[string]string, bool)
	StoreFields(fields map[string]string) error

	GetGroups(fieldID string) (map[string]string, bool)
	StoreGroups(fieldID string, groups map[string]string) error

	GetSchedule(groupID string) ([]Entry, bool)
	StoreSchedule(groupID string, entries []Entry) error
}
