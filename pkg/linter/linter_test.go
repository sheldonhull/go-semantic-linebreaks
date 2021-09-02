package linter_test

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/matryer/is"
	proj "github.com/sheldonhull/go-semantic-sentences/pkg/linter"
)

// // Logger contains the package level logger provided from internal logger package that wraps up zerolog.
// // var Log *logger.Logger //nolint: gochecknoglobals.
// //nolint
// func TestMain(m *testing.M) {
// 	// c := zl.Config{
// 	// 	Enable:                false,
// 	// 	ConsoleLoggingEnabled: false,
// 	// 	EncodeLogsAsJson:      false,
// 	// 	FileLoggingEnabled:    false,
// 	// 	Directory:             "",
// 	// 	Filename:              "",
// 	// 	MaxSize:               0,
// 	// 	MaxBackups:            0,
// 	// 	MaxAge:                0,
// 	// 	Level:                 "info",
// 	// }
// 	// _ = zl.InitLogger(c)

// 	m.Run()
// 	os.Exit(m.Run())
// }

func TestCountViolations(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name     string
		filepath string
		want     int
	}

	testCases := []testCase{
		{
			name:     "1 violation counted",
			filepath: "test-files/1-violations-multiple-lines.md",
			want:     1,
		},
		{
			name:     "2 violations counted",
			filepath: "test-files/2-violations-multiple-lines.md",
			want:     2,
		},
		{
			name:     "11 violations counted",
			filepath: "test-files/11-violations-one-line.md",
			want:     11,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			is := is.New(t)
			f, err := filepath.Abs(tc.filepath)
			if err != nil {
				t.Fatalf("cannot find test file: %q", tc.filepath)
			}
			content, err := ioutil.ReadFile(f)
			if err != nil {
				t.Fatalf("ioutil.ReadFile(f): %q", err)
			}
			got := proj.CountViolations(content)
			want := tc.want
			is.Equal(want, got) // violation count matches expected count
		})
	}
}

// func TestFixViolations(t *testing.T) {
// 	t.Parallel()

// 	type testCase struct {
// 		name          string
// 		filepath      string
// 		filepathFixed string
// 	}

// 	testCases := []testCase{
// 		{
// 			name:          "1 violation fixed",
// 			filepath:      "test-files/1-violations-multiple-lines.md",
// 			filepathFixed: "test-files/1-violations-multiple-lines-fixed.md",
// 		},
// 		{
// 			name:          "2 violations",
// 			filepath:      "test-files/2-violations-multiple-lines.md",
// 			filepathFixed: "test-files/2-violations-multiple-lines-fixed.md",
// 		},
// 		{
// 			name:          "11 violations",
// 			filepath:      "test-files/11-violations-one-line.md",
// 			filepathFixed: "test-files/11-violations-one-line-fixed.md",
// 		},
// 	}

// 	for _, tc := range testCases {
// 		tc := tc
// 		t.Run(tc.name, func(t *testing.T) {
// 			t.Parallel()
// 			is := is.New(t)
// 			f, err := filepath.Abs(tc.filepath)
// 			if err != nil {
// 				t.Fatal("cannot find test file: ", tc.filepath)
// 			}
// 			content, err := ioutil.ReadFile(f)
// 			if err != nil {
// 				t.Fatalf("ioutil.ReadFile(f): %v", err)
// 			}
// 			fixed, err := filepath.Abs(tc.filepathFixed)
// 			if err != nil {
// 				t.Fatalf("cannot find test file: %q", tc.filepathFixed)
// 			}
// 			fixedContent, err := ioutil.ReadFile(fixed)
// 			if err != nil {
// 				t.Fatalf("ioutil.ReadFile(fixed): %v", err)
// 			}

// 			got := proj.FormatSemanticLineBreak(content)
// 			gotTrimmed := strings.Trim(got, " ")
// 			want := strings.Trim(string(fixedContent), " ")
// 			is.Equal(string(want), gotTrimmed) // FormatSemanticLineBreak matches fixed file
// 		})
// 	}
// }

// func TestRun(t *testing.T) {
// 	type testCase struct {
// 		name          string
// 		filepath      string
// 		filepathFixed string
// 	}

// 	testCases := []testCase{
// 		{
// 			name:          "1 violation fixed",
// 			filepath:      "test-files/1-violations-multiple-lines.md",
// 			filepathFixed: "test-files/1-violations-multiple-lines-fixed.md",
// 		},
// 		{
// 			name:          "2 violations",
// 			filepath:      "test-files/2-violations-multiple-lines.md",
// 			filepathFixed: "test-files/2-violations-multiple-lines-fixed.md",
// 		},
// 		{
// 			name:          "11 violations",
// 			filepath:      "test-files/11-violations-one-line.md",
// 			filepathFixed: "test-files/11-violations-one-line-fixed.md",
// 		},
// 	}

// 	for _, tc := range testCases {
// 		tc := tc
// 		t.Run(tc.name, func(t *testing.T) {
// 			t.Parallel()
// 			is := is.New(t)
// 			f, err := filepath.Abs(tc.filepath)
// 			if err != nil {
// 				t.Fatal("cannot find test file: ", tc.filepath)
// 			}
// 			content, err := ioutil.ReadFile(f)
// 			if err != nil {
// 				t.Fatalf("ioutil.ReadFile(f): %v", err)
// 			}
// 			fixed, err := filepath.Abs(tc.filepathFixed)
// 			if err != nil {
// 				t.Fatalf("cannot find test file: %q", tc.filepathFixed)
// 			}
// 			fixedContent, err := ioutil.ReadFile(fixed)
// 			if err != nil {
// 				t.Fatalf("ioutil.ReadFile(fixed): %v", err)
// 			}
// 			var stdout bytes.Buffer
// 			args := []string{"-source", tc.filepath}
// 			err = proj.Run(args, &stdout)
// 			is.NoErr(err) // run should not fail

// 			got := proj.FormatSemanticLineBreak(content)
// 			gotTrimmed := strings.Trim(got, " ")
// 			want := strings.Trim(string(fixedContent), " ")
// 			is.Equal(string(want), gotTrimmed) // FormatSemanticLineBreak matches fixed file
// 		})
// 	}
// }
