package strutil

import (
	"strings"
)

// CleanEmpty used to clean empty strings
func CleanEmpty(s []string) []string {
	var out []string
	for _, item := range s {
		if strings.TrimSpace(item) != "" {
			out = append(out, item)
		}
	}
	return out
}

// Include  is used to return `true` if the target string t is in the
// slice.
// Ex: fmt.Println(Include(strs, "grape"))
func Include(vs []string, t string) bool {
	return Index(vs, t) >= 0
}

// Index is used to return the first index of the target string `t`, or
// -1 if no match is found.
// Ex: fmt.Println(Index(strs, "pear"))
func Index(vs []string, t string) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}
