package model

import (
	"reflect"
	"testing"
)

func TestParseSubgroup(t *testing.T) {
	tests := []struct {
		name     string
		subgroup string
		want     *Subgroup
	}{
		{"A", "A", new(A)},
		{"B", "B", new(B)},
		{"empty", "", nil},
		{"invalid", "test", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseSubgroup(tt.subgroup); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseSubgroup() = %v, want %v", got, tt.want)
			}
		})
	}
}
