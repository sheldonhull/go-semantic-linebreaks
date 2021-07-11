package main

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestCountViolations(t *testing.T) {
	type testCase struct {
		name     string
		filepath string
		want     int
	}
	testCases := []testCase{
		{
			name:     "1 violation",
			filepath: "test-files/1-violations-multiple-lines.md",
			want:     1,
		},
		{
			name:     "2 violations",
			filepath: "test-files/2-violations-multiple-lines.md",
			want:     2,
		},
		{
			name:     "18 violations",
			filepath: "test-files/18-violations-one-line.md",
			want:     18,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f, err := filepath.Abs("test-files/1-violations-multiple-lines.md")
			if err != nil {
				t.Fatal("cannot find test file: [test-files/1-violations-multiple-lines.md]")
			}
			content, err := ioutil.ReadFile(f)
			if err != nil {
				t.Fatalf("ioutil.ReadFile(f): %v", err)
			}
			got := CountViolations(content)
			want := 1
			if got != want {
				t.Errorf("CountViolations() = %v, want %v", got, want)
			}
		})
	}
}
