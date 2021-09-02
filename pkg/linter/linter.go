package linter

import (
	"regexp"
)

// CountViolations counts the number of lines that would need to be fixed by adding semantic line break. It returns an integer value of the violation count found.
func CountViolations(content []byte) int {
	// re := regexp.MustCompile(`(?is)(\w[.?])(\s+)(\w)?`)
	// re := regexp.MustCompile(`(?is)(?:[a-zA-Z"';\]][.?!])(\s+)[a-zA-Z"';\]]`)
	re := regexp.MustCompile(`(?is)([a-zA-Z"';\]][.?!])\s+`)
	matches := re.FindAllString(string(content), -1)

	return len(matches)
}

// FormatSemanticLineBreak takes a byte array and searches for any violations of semantic line breaks and then fixes with line breaks.
func FormatSemanticLineBreak(content []byte) (formatted string) {
	// re := regexp.MustCompile(`(?is)(?:[a-zA-Z][.?])(\s)(?:[a-zA-Z])`)
	re := regexp.MustCompile(`(?is)([a-zA-Z"';\]][.?!])\s+`)
	formatted = re.ReplaceAllString(string(content), "$1\n")

	return formatted
}
