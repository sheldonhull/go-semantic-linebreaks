package main_test

import (
	"bytes"
	"testing"

	iz "github.com/matryer/is"
	"github.com/pterm/pterm"
	proj "github.com/sheldonhull/go-semantic-linebreaks/cmd/go-semantic-linebreaks"
)

func TestMain(t *testing.T) {
	is := iz.New(t)
	pterm.DisableStyling()
	args := []string{"-source", "test.md"}
	var stdout bytes.Buffer
	err := proj.Run(args, &stdout)
	is.NoErr(err) // run should not fail
}
