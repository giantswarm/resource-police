// Package intersect is a version of https://github.com/juliangruber/go-intersect/blob/master/intersect.go
// adapted especially for string slices.
package intersect

import (
	"sort"
)

// StringSliceSorted returns a slice of strings
// occurring both in a and b.
// a needs to be sorted.
// Complexity: O(n * log(n))
func StringSliceSorted(a, b []string) []string {
	set := make([]string, 0)

	for i := 0; i < len(a); i++ {
		el := a[i]
		idx := sort.Search(len(b), func(i int) bool {
			return b[i] == el
		})
		if idx < len(b) && b[idx] == el {
			set = append(set, el)
		}
	}

	return set
}
