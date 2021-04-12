// Package intersect is a version of https://github.com/juliangruber/go-intersect/blob/master/intersect.go
// adapted especially for string slices.
package intersect

import (
	"reflect"
	"testing"
)

func TestStringSlice(t *testing.T) {
	type args struct {
		a []string
		b []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "First",
			args: args{[]string{"bar", "foo", "one", "two"}, []string{"one", "two", "three", "four", "five"}},
			want: []string{"one", "two"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringSlice(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringSliceSorted() = %v, want %v", got, tt.want)
			}
		})
	}
}
