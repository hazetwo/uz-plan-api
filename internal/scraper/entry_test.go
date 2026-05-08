package scraper

import (
	"reflect"
	"testing"
)

func Test_strPtr(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want *string
	}{
		{"valid", "test", new("test")},
		{"empty", "", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := strPtr(tt.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("strPtr() = %v, want %v", got, tt.want)
			}
		})
	}
}
