package copy

import "strings"

var forbiddenSubstrings = []string{
	"scam",
	"fraud",
	"same-company",
}

// ContainsForbiddenLanguage reports whether s contains any forbidden substring,
// compared case-insensitively.
func ContainsForbiddenLanguage(s string) bool {
	lower := strings.ToLower(s)
	for _, forbidden := range forbiddenSubstrings {
		if strings.Contains(lower, forbidden) {
			return true
		}
	}
	return false
}
