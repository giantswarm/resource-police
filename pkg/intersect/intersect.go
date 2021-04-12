// Package intersect is a version of https://github.com/juliangruber/go-intersect/blob/master/intersect.go
// adapted especially for string slices.
package intersect

// StringSlice returns a slice of strings
// occurring both in a and b.
// Complexity: O(n^2)
func StringSlice(a []string, b []string) []string {
	set := []string{}

	for i := 0; i < len(a); i++ {
		for j := 0; j < len(b); j++ {
			if a[i] == b[j] {
				set = append(set, a[i])
			}
		}
	}

	return set
}
