package linter

import (
	"fmt"
	"regexp"

	"github.com/yuin/goldmark"
	// "github.com/yuin/goldmark/extension".
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
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

// ParseText will use parse input text using Goldmark and invoke the string replacement only on plain text.
//
// If the text is parsed and part of a table, header, and other such special text formats, then this won't be reformatted to avoid breaking markdown text.
func ParseText(source []byte) error {
	var err error
	doc := goldmark.DefaultParser().Parse(text.NewReader(source))
	err = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		status := ast.WalkContinue

		if n.Kind() == ast.KindText || n.Kind() == ast.KindTextBlock {
			// Apply any logic here

			// gets the title
			// fmt.Println(string(n.Text(source)))
			// get block content
			// fmt.Println(string(n.NextSibling().Text(source)))

			// Stop walking
		}
		return status, err
	})
	if err != nil {
		return fmt.Errorf("parsetext: %w", err)
	}
	return nil
}
