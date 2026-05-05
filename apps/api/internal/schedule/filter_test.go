package schedule

import "testing"

func TestMatchesDay(t *testing.T) {
	tests := []struct {
		name   string
		modify func(entry *Entry)
		date   string
		want   bool
	}{
		{"same day", func(e *Entry) { e.Date = new("2026-01-01") }, "2026-01-01", true},
		{"same day in the past", func(e *Entry) { e.Date = new("1999-05-05") }, "1999-05-05", true},
		{"different year", func(e *Entry) { e.Date = new("2025-12-11") }, "2024-12-11", false},
		{"different day", func(e *Entry) { e.Date = new("2026-11-10") }, "2026-11-11", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Entry{}
			tt.modify(&e)
			if got := matchesDay(e, new(tt.date)); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}

		})
	}
}

func TestMatchesDay_Nil(t *testing.T) {
	t.Run("nil argument date", func(t *testing.T) {
		if got := matchesDay(Entry{}, nil); got != true {
			t.Errorf("got %v, want true", got)
		}
	})
	t.Run("nil entry date", func(t *testing.T) {
		if got := matchesDay(Entry{Date: nil}, new("2026-03-03")); got != false {
			t.Errorf("got %v, want true", got)
		}
	})

}

func TestMatchesWeek(t *testing.T) {
	tests := []struct {
		name   string
		modify func(entry *Entry)
		date   string
		want   bool
	}{
		{"same week", func(e *Entry) { e.Date = new("2026-05-04") }, "2026-05-08", true},
		{"same week in the past", func(e *Entry) { e.Date = new("1999-02-08") }, "1999-02-11", true},
		{"different year", func(e *Entry) { e.Date = new("2025-12-11") }, "2024-12-11", false},
		{"different day", func(e *Entry) { e.Date = new("2026-05-10") }, "2026-05-11", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Entry{}
			tt.modify(&e)
			if got := matchesWeek(e, new(tt.date)); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}

		})
	}
}

func TestMatchesWeek_Nil(t *testing.T) {
	t.Run("nil argument week", func(t *testing.T) {
		if got := matchesWeek(Entry{}, nil); got != true {
			t.Errorf("got %v, want true", got)
		}
	})
	t.Run("nil entry date", func(t *testing.T) {
		if got := matchesWeek(Entry{Date: nil}, new("2026-03-03")); got != false {
			t.Errorf("got %v, want true", got)
		}
	})

}

func TestMatchesSubgroup(t *testing.T) {
	tests := []struct {
		name     string
		modify   func(entry *Entry)
		subgroup Subgroup
		want     bool
	}{
		{"A", func(e *Entry) { e.Subgroup = ParseSubgroup("A") }, "A", true},
		{"B", func(e *Entry) { e.Subgroup = ParseSubgroup("B") }, "B", true},
		{"empty falls for all", func(e *Entry) { e.Subgroup = ParseSubgroup("") }, "", true},
		{"different subgroup", func(e *Entry) { e.Subgroup = ParseSubgroup("A") }, "C", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Entry{}
			tt.modify(&e)
			if got := matchesSubgroup(e, new(tt.subgroup)); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}

		})
	}
}

func TestMatchesSubgroup_Nil(t *testing.T) {
	t.Run("nil argument group", func(t *testing.T) {
		if got := matchesSubgroup(Entry{}, nil); got != true {
			t.Errorf("got %v, want true", got)
		}
	})
	t.Run("nil entry date", func(t *testing.T) {
		if got := matchesSubgroup(Entry{Subgroup: nil}, ParseSubgroup("A")); got != true {
			t.Errorf("got %v, want true", got)
		}
	})

}
