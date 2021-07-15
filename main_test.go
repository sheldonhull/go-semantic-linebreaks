package main_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/matryer/is"
	proj "github.com/sheldonhull/go-semantic-sentences"
	zl "github.com/sheldonhull/go-semantic-sentences/internal/logger"
)

// Logger contains the package level logger provided from internal logger package that wraps up zerolog.
// var Log *logger.Logger //nolint: gochecknoglobals.
//nolint
func TestMain(m *testing.M) {
	c := zl.Config{
		Enable:                false,
		ConsoleLoggingEnabled: false,
		EncodeLogsAsJson:      false,
		FileLoggingEnabled:    false,
		Directory:             "",
		Filename:              "",
		MaxSize:               0,
		MaxBackups:            0,
		MaxAge:                0,
		Level:                 "info",
	}
	_ = zl.InitLogger(c)

	m.Run()
	os.Exit(m.Run())
}

func TestCountViolations(t *testing.T) {
	t.Parallel()

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
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			f, err := filepath.Abs("test-files/1-violations-multiple-lines.md")
			if err != nil {
				t.Fatal("cannot find test file: [test-files/1-violations-multiple-lines.md]")
			}
			content, err := ioutil.ReadFile(f)
			if err != nil {
				t.Fatalf("ioutil.ReadFile(f): %v", err)
			}

			want := 1
			got := proj.CountViolations(content)
			is.Equal(want, got) // violation count
		})
	}
}
