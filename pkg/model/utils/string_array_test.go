package utils

import (
	"testing"
)

func TestStringArray_Contains(t *testing.T) {
	type args struct {
		search string
	}
	tests := []struct {
		name string
		arr  StringArray
		args args
		want bool
	}{
		{name: "Contains", arr: StringArray{"U", "PG", "12A", "12", "15", "18", "R18"}, args: struct{ search string }{search: "PG"}, want: true},
		{name: "NotContains", arr: StringArray{"U", "PG", "12A", "12", "15", "18", "R18"}, args: struct{ search string }{search: "66"}, want: false},
		{name: "EmptyNotContains", arr: StringArray{"U", "PG", "12A", "12", "15", "18", "R18"}, args: struct{ search string }{search: ""}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.arr.Contains(tt.args.search); got != tt.want {
				t.Errorf("StringArray.Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringArray_String(t *testing.T) {
	tests := []struct {
		name string
		arr  StringArray
		want string
	}{
		{name: "1", arr: StringArray{"1", "2", "3", "4", "5"}, want: "1, 2, 3, 4, 5"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.arr.String(); got != tt.want {
				t.Errorf("StringArray.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
